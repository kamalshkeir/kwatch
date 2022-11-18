package kwatch

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/kamalshkeir/klog"
)

func changesDetected(since time.Time, root string, dirs ...string) bool {
	var changed bool
	var err error
	if len(dirs) == 0 {
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if info.ModTime().After(since) {
				changed = true
			}
			return nil
		})
	} else {
		for _, d := range dirs {
			err = filepath.Walk(root+"/"+d, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				if info.ModTime().After(since) {
					changed = true
				}
				return nil
			})
		}
	}
	if err != nil {
		return false
	}
	return changed
}

type BuildFunc func() error

func ExecCommand(command string, args ...string) BuildFunc {
	return func() error {
		cmd := exec.Command(command, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("error building: \"%s %s\": %w", command, strings.Join(args, " "), err)
		}
		return nil
	}
}

type LaunchFunc func() (stop func(), err error)

func LaunchCommand(command string, args ...string) LaunchFunc {
	return func() (func(), error) {
		cmd := exec.Command(command, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Start()
		if err != nil {
			return nil, fmt.Errorf("error LaunchCommand: \"%s %s\": %w", command, strings.Join(args, " "), err)
		}
		return func() {
			cmd.Process.Kill()
		}, nil
	}
}

func Start(before []BuildFunc, run LaunchFunc, after []BuildFunc) (func(), error) {
	for _, fn := range before {
		err := fn()
		if err != nil {
			return nil, err
		}
	}
	stop, err := run()
	if err != nil {
		return nil, err
	}
	for _, fn := range after {
		err := fn()
		if err != nil {
			stop()
			return nil, err
		}
	}
	return stop, nil
}

func Watch(every time.Duration, root string, dirs ...string) {
	if every == 0 {
		every = 500 * time.Millisecond
	}
	root = filepath.ToSlash(root)

	ssp := strings.Split(root, "/")
	projName := ssp[len(ssp)-1]
	if projName == "" {
		projName = ssp[len(ssp)-2]
		if projName == "" {
			klog.Printf("rd%s is empty",projName)
			return
		}
	}
	if runtime.GOOS == "windows" {
		projName += ".exe"
	}
	var stop func()
	var err error
	var lastScan time.Time

	if len(dirs) > 0 {
		for _, d := range dirs {
			klog.Printfs("grwatching %s", root+"/"+d)
		}
	} else {
		klog.Printfs("grwatching %s", root)
	}

	for {
		if !changesDetected(lastScan, root, dirs...) {
			time.Sleep(every)
			continue
		}

		if stop != nil {
			klog.Printfs("ylRestarting...")
			stop()
		}
		stop, err = Start(
			[]BuildFunc{ExecCommand("go", "install", root)},
			LaunchCommand(projName),
			nil,
		)
		if err != nil {
			klog.Printfs("error katcher:", err)
		}
		lastScan = time.Now()
		time.Sleep(every)
	}
}