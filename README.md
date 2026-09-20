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
#### If you find this project useful, you can support us with a coin!
[![Donate via Monobank](https://img.shields.io/badge/monobank-23232b?style=for-the-badge&logo=data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABwAAAAcBAMAAACAI8KnAAABfmlDQ1BJQ0MgUHJvZmlsZQAAeJylkL9LQlEcxY+aKGk41ODQ8AZpCIWypbEfgxAiYgZZLXp9avCePt5TIhobWh1cKlqy6D+oLfoHgiCopiBqbiiIIOR1rk8QQqfu497vh3O/5757D+BuakK3RmYAvVo3M4klZT23ofheEEAYfszCmxeWsZhOJzF0fD3AJet9TJ41vG/gCBRVSwAuP3leGGadvEBO7dQNyU3yhKjki+QzctTkBcl3Ui84/Ca57PC3ZDObWQbcQbJSdjgqueCwfIsiKqZO1sgRXWuI3n3kS4JqdW2VdbI7LWSQwBIUFNDANjTUEWOtMrPBvnjXl0KNHsHVwC5MOsqo0Bul2uCpKmuJuspPYweHzP5vplZpLu78IbgCeF9t+3Ma8B0DnQPb/jm17U4b8DwBN62+v9ZinO/Um30tcgKE9oHL675WOAeumHH42cib+a7k4XSXSsDHBTCWA8aZ9ejmf/edvHv7aD8C2T0geQscHgFT7A9t/QIdXnS014JJvwAAAB5QTFRFJCYuMzVC4+PkWFphlpaafH2CPkFIJSYuAAAAAAAAFb5GYwAAAAh0Uk5T///////////zgVLUAAAAp0lEQVR42qWQrwrCUBjFf/tEvnonGC0mEfEJTDar7yHe5jPY7gP4IAaDsCYMtqRiM1lka3LH/BNmEE2yXzhwyjmHE4R8Ivxh6xAsvpK9944lzrEEQVGLOgsWhFVUbotBsekTQWN0zQ+3S97N/FpPPYHpIzDnvGmG7dKJrwoAjvLenHQMjFtiRdO7aWVhbNIynCDMnvvdPF4kM431d1WNc1RVUZRKavACBhwu8Eos8DoAAAAASUVORK5CYII=)](https://send.monobank.ua/23MWDFXsRa)