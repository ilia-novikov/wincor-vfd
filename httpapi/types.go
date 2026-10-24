package httpapi

import "github.com/ilia-novikov/wincor-vfd/display"

// Device is the subset of *display.Display a request handler needs.
type Device interface {
	PrintText(text string) error
	ClearDisplay() error
	Initialize() error
	MoveCursor(x, y int) error
	SendBackspace() error
	SendLineFeed() error
	SendCarriageReturn() error
	DeleteToEndOfLine() error
	Close() error
}

type OpenFunc func(charset display.Charset) (Device, error)

type ProbeFunc func() error
