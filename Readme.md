### Watcher or Auto-reloader 

### install
```shell
go install github.com/kamalshkeir/kwatch/cmd/kwatch
```
or get the binary from Releases

### Then you can run:
```shell
kwatch --root ${PWD} (will watch all files at root)
kwatch --root ${PWD} --watch assets/templates,assets/static (will watch only '${PWD}/assets/templates' and '${PWD}/assets/static' folders)
```