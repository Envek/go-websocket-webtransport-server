# Dual WebSocket and WebTransport chat server in Go

Just for fun and learning.

## Setup

 1. Get Go 1.18 (webtransport implementation doesn't compile on 1.19 or newer at the moment)

 2. Install dependencies

     ```sh
     go mod download
     ```

 2. Generate some local certificates with [mkcert](https://github.com/FiloSottile/mkcert)

    ```sh
    mkcert -install
    mkcert localhost 127.0.0.1 ::1
    ```

## Run

```sh
go run .
```

## Enjoy

### WebSocket

You can use [websocat](https://github.com/vi/websocat) for some chatting:

```sh
websocat --exit-on-eof ws://127.0.0.1:8090/chat
```

### WebTransport

Open http://localhost:8090/ in a modern browser\* and start chatting.

\* Currently Chrome 97+, Firefox 114+, and Edge 98+ supports WebTransport. See https://caniuse.com/webtransport for the up to date list of supported browsers.

### Or both!

Of course you can chat via both protocol at the same time!
