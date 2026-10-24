```
┌─────────────────────────────────────────────────────────────┐
│  ▄▄   ▄▄ ▄▄ ▄▄  ▄▄  ▄▄▄▄  ▄▄▄  ▄▄▄▄      ▄▄ ▄▄ ▄▄▄▄▄ ▄▄▄▄   │
│  ██ ▄ ██ ██ ███▄██ ██▀▀▀ ██▀██ ██▄█▄ ▄▄▄ ██▄██ ██▄▄  ██▀██  │
│   ▀█▀█▀  ██ ██ ▀██ ▀████ ▀███▀ ██ ██      ▀█▀  ██    ████▀  │
└─────────────────────────────────────────────────────────────┘
```

---

A Go driver for Wincor Nixdorf VFD (vacuum fluorescent display) customer-facing
displays used in ATMs/POS terminals, communicating over USB HID. Was tested on
BA63, but other models should be potentially supported as well.

## Install

```sh
go get github.com/ilia-novikov/wincor-vfd
```

## Usage

```go
d, err := display.Open(display.CharsetUS)
if err != nil {
	log.Fatalf("opening display: %v", err)
}
defer d.Close()

d.Initialize()
d.PrintText("Hello, wincor-vfd!")
```

US (`CharsetUS`) and Russian/CP866 (`CharsetRU`) charsets are supported. See
`display/operations.go` for the full operation set (cursor movement, clearing,
backspace, line feed, etc). Some more charsets could be potentially added,
please let me know if you want to add a missing one.

## CLI

`cmd/vfdctl` exposes the same operations from the command line:

```sh
go run ./cmd/vfdctl init
go run ./cmd/vfdctl print "Hello, wincor-vfd!"
go run ./cmd/vfdctl --charset ru print "Привет"
go run ./cmd/vfdctl move 3 1
```

Use `--help` to get a more detailed description of supported commands.

## REST API

`cmd/vfdhttp` exposes the same operations over HTTP, for cases where the
display is driven remotely rather than from a local process. Each request
opens and closes the device itself, so the server never holds it open and
can coexist with `vfdctl` or other clients between requests.

```sh
go run ./cmd/vfdhttp --addr localhost:8080

curl -X POST localhost:8080/v1/print -d '{"text": "Hello, wincor-vfd!"}'
curl -X PUT localhost:8080/v1/cursor -d '{"x": 3, "y": 1}'
```

Setting the `VFD_API_TOKEN` environment variable requires an `Authorization: Bearer <token>`
header on all `/v1/*` routes. See `httpapi/openapi.yaml` for the full reference.

## License

See [LICENSE](LICENSE).
