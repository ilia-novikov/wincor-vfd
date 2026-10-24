package display

import "fmt"

// DeviceNotFoundError means no matching display was found on the USB bus.
type DeviceNotFoundError struct {
}

func (e *DeviceNotFoundError) Error() string {
	return "display interface not found"
}

// CharsetError means an unsupported Charset was used.
type CharsetError struct {
	Charset Charset
}

func (e *CharsetError) Error() string {
	return fmt.Sprintf("only US and RU charsets are supported, got %q", e.Charset)
}

// PayloadTooLargeError means a command payload exceeded the display's
// maximum report size.
type PayloadTooLargeError struct {
	Size int
	Max  int
}

func (e *PayloadTooLargeError) Error() string {
	return fmt.Sprintf("payload of %d bytes exceeds maximum of %d bytes", e.Size, e.Max)
}
