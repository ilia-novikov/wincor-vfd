package display

import (
	"errors"
	"fmt"

	hid "github.com/sstallion/go-hid"
)

const (
	VendorID  = 0x0AA7
	ProductID = 0x0200
	UsagePage = 0xFF45
	Usage     = 0x2400

	ReportSize   = 32
	CharsPerLine = 20
)

var (
	STX = []byte{0x02, 0x00}
)

type Display struct {
	device  *hid.Device
	charset Charset // TODO: nullable?
}

// Open finds and opens the display, using the given charset
func Open(charset Charset) (*Display, error) {
	path, err := findDevicePath()
	if err != nil {
		return nil, err
	}

	dev, err := hid.OpenPath(path)
	if err != nil {
		return nil, fmt.Errorf("opening display at %s: %w", path, err)
	}

	return &Display{device: dev, charset: charset}, nil
}

// Close releases the HID device handle
func (d *Display) Close() error {
	if d.device == nil {
		return nil
	}
	return d.device.Close()
}

// ClearDisplay clears the whole display
func (d *Display) ClearDisplay() error {
	return d.write(clearDisplayCmd)
}

// SetCharset selects the display's configured charset
func (d *Display) SetCharset() error {
	cmd, err := setCharsetCmd(d.charset)
	if err != nil {
		return err
	}
	return d.write(cmd)
}

func findDevicePath() (string, error) {
	var path string

	err := hid.Enumerate(VendorID, ProductID, func(info *hid.DeviceInfo) error {
		if info.UsagePage == UsagePage && info.Usage == Usage {
			path = info.Path
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("enumerating HID devices: %w", err)
	}

	if path == "" {
		return "", errors.New("display interface not found")
	}

	return path, nil
}

func createPacket(data []byte) []byte {
	msg := make([]byte, 0, len(STX)+1+len(data))
	msg = append(msg, STX...)
	msg = append(msg, byte(len(data)))
	msg = append(msg, data...)

	// TODO: Report size validation
	report := make([]byte, ReportSize)
	report[0] = 0x00
	copy(report[1:], msg)
	return report
}

func (d *Display) write(data []byte) error {
	_, err := d.device.Write(createPacket(data))
	return err
}
