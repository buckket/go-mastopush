#mastopush

This is an example project to illustrate the usage of the **go-mastopush** library.
We create an HTTP server that handles incoming Web Push requests and forwards the decrypted
notification payload to a [Gotify](https://gotify.net/) server.

## Installation

```sh
git clone git@github.com:buckket/go-mastopush.git
cd go-mastopush

go build example/mastopush.go
```

## Usage

- Copy and adjust `config.toml`
  - The Mastodon API token needs to have the `push` scope.
- Run mastopush: `./mastopush -config config.toml`
