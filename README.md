<h1 align="center">Pathfinding</h1>

<p align="center">
  <a href="./LICENSE"><img src="https://img.shields.io/badge/⚖️ license-MIT-blue" alt="MIT License"></a>
  <img src="https://img.shields.io/github/stars/keelus/pathfinding?color=red&logo=github" alt="stars">
</p>

## ℹ️ Description
A visual implementation of Dijkstra and A*, side by side, written in golang with [Ebitengine](https://ebitengine.org/).

## 📸 Screenshots
<img src="https://github.com/keelus/pathfinding/assets/86611436/fd1212cc-13b7-4bfb-977b-4e442a745291"/>


## ⬇️ Install & run it
The project is compatible with Windows, Linux and macOS.

To use it, simply download the [latest release](https://github.com/keelus/pathfinding/releases/latest) binary file and execute it.

### 🐧 Linux & macOS
To make the downloaded binary executable, run:
```bash
chmod +x pathfinding_<rest of the file>
```
Then, you can open it running:
```bash
./pathfinding_<rest of the file>
```

## Compile
### 🪟 Windows
You can compile the app in Windows directly without a C compiler. Just run:
```bash
go mod tidy
```
```bash
go build -o pathfinding.exe
```
#### Add an icon (optional. Requires [go-winres](github.com/tc-hib/go-winres))
```bash
go-winres simply --icon assets/icons/greenFlag.png --manifest gui
```
```bash
go build -o pathfinding.exe -ldflags -H=windowsgui
```
### 🐧 Linux or macOS
Compiling an Ebitengine app in linux and macOS requires having a c compiler installed. Check [ebitengine dependencies](https://ebitengine.org/en/documents/install.html#Installing_dependencies).
Once done, simply run:
```bash
go mod tidy
```
```bash
go build -o pathfinding
```

## ⚖️ License
This project is open source under the terms of the [MIT License](./LICENSE)

<br />

Made by <a href="https://github.com/keelus">keelus</a> ✌️
