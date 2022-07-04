# Gopher Find 🐹🔎

Gopher Find is a blazingly fast alternative to [Sherlock](https://github.com/sherlock-project/sherlock) written in [Go](https://go.dev/).

![gopher-find-usage.gif](./.github/assets/gopher-find-usage.gif)

This project is still under development and new features will be added soon.

All the websites Gopher Find look for were extracted from the Sherlock project.

# Installation and usage

Assuming you already have Go installed into your machine.

```sh
git clone https://github.com/drull1000/gopher-find
```

```sh
cd gopher-find
```

```
go mod tidy
```

```sh
go run cmd/main.go <username>
```

or if you want to compile the program:

```sh
go build cmd/main.go
```

```sh
./main <username>
```
>Note that the json file needs to be in the correct path.

The script will start hunting the urls for you and after going through all of them, it will generate a `<username>.txt` containing all the valid links.

## Known issues

- Some few websites return false positives. This is caused because some of them either require javascript, captchas or use cloudflare to block requests.

## Contributing

Seeing something wrong? Want to help?

Feel free to fork the project, make changes, open pull requests and issues. Any contribution is welcome.
