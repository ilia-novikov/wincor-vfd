package display

import "fmt"

type DeviceNotFoundError struct {
}

func (e *DeviceNotFoundError) Error() string {
	return "display interface not found"
}

type CharsetError struct {
	Charset Charset
}

func (e *CharsetError) Error() string {
	return fmt.Sprintf("only US and RU charsets are supported, got %q", e.Charset)
}

type PayloadTooLargeError struct {
	Size int
	Max  int
}

func (e *PayloadTooLargeError) Error() string {
	return fmt.Sprintf("payload of %d bytes exceeds maximum of %d bytes", e.Size, e.Max)
}
