<h1 align="center">Pathfinding</h1>

This project consists on a visual implementation of Dijkstra and A*, written in Golang with [Ebitengine](https://ebitengine.org/).

Work in progress.

Inspiration: [https://youtu.be/cSxnOm5aceA](https://youtu.be/cSxnOm5aceA)

### Compile
#### Windows:
```bash
go mod tidy
```
```bash
go-winres simply --icon assets/icons/greenFlag.png --manifest gui
```
```bash
go build -o pathfinding.exe -ldflags -H=windowsgui
```