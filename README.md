# 🐾 skypaw
**Lightweight weather CLI-Tool written on Go.**

![Skypaw Demo](./demo/demo.gif)

## 🚀 Installation
### Windows (WinGet)
Open the cmd and run:
```
winget install skypaw
```
### Arch Linux (AUR)
Run in shell:
```
yay -S skypaw-bin
```

### Others platforms
Check the [release page](https://github.com/zenpaw-labs/skypaw/releases).

## ⌨️ Usage
Simply run the following command:
```
skypaw
```

## 🛠 Development
### Build from sources
To build the production-ready binary (only your platform):
```
go build -ldflags="-X 'github.com/zenpaw-labs/skypaw/cmd.semVersion=dev' -s -w" ./cmd/skypaw
```
To run the build binaries ready for release:
```
goreleaser release --snapshot --clean
```

## 📄 License
This project is licensed under the [MIT License](LICENSE)

## ☕ Support the project
#### If you find this project useful, you can buy me a coffee!
[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/I3I51YMHG4)<br>