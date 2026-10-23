# wincor-vfd

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

## License

MIT -> [LICENSE](LICENSE).
