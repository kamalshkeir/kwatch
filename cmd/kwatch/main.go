package main

import (
	"flag"
	"strings"
	"time"

	"github.com/kamalshkeir/klog"
	"github.com/kamalshkeir/kwatch"
)

func main() {
	root := flag.String("root", "", "root is fullPath to the project")
	watch := flag.String("watch", "", "directory to watch inside root,if empty, will take all files and dirs inside root")
	every := flag.Int("every", 313, "time in milliseconds")
	flag.Parse()
	if root == nil || *root == "" {
		klog.Printfs("rderror: root tag not specified")
		return
	}
	sp := strings.Split(*watch, ",")
	kwatch.Watch(time.Duration(*every)*time.Millisecond, *root, sp...)
}