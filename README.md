# Gopher Find 🐹🔎

Gopher Find is a blazingly fast alternative to [Sherlock](https://github.com/sherlock-project/sherlock) written in [Go](https://go.dev/).

![gopher-find-usage.gif](./.github/assets/gopher-find-usage.gif)

This project is still under development and new features will be added soon.

All the websites Gopher Find look for were extracted from the Sherlock project.

![links.gif](./.github/assets/links.gif)

# Installation and usage

Assuming you already have Go installed into your machine.

`git clone https://github.com/drull1000/gopher-find`

`cd Gopher-Find`

`go mod tidy`

`go run cmd/main.go <username>`

or if you want to compile the program

`go build cmd/main.go <username>`
`./main `

The script will start hunting the urls for you and after going through all of them, it will generate a `<username>.txt` containing all the valid links.

## Known issues

- Some few websites return false positives. This is caused because some of them either require javascript, captchas or use cloudflare to block requests.

## Contributing

Feel free to fork the project, make changes and open pull requests. Any contribution is welcome.
