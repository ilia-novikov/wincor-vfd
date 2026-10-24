package display

import (
	"fmt"

	hid "github.com/sstallion/go-hid"
)

// TODO: Add support for other display params?
const (
	VendorID  = 0x0AA7
	ProductID = 0x0200
	UsagePage = 0xFF45
	Usage     = 0x2400

	ReportSize   = 32
	CharsPerLine = 20
)

var (
	STX = []uint8{0x02, 0x00}
)

type Display struct {
	device  *hid.Device
	charset Charset
}

// Open finds and opens the display, using the given charset.
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

// Close releases the HID device handle.
func (d *Display) Close() error {
	if d.device == nil {
		return nil
	}

	return d.device.Close()
}

// Initialize clears the display, selects the configured charset, and
// resets the cursor to (0, 0).
func (d *Display) Initialize() error {
	if err := d.ClearDisplay(); err != nil {
		return fmt.Errorf("initializing display: clearing: %w", err)
	}
	if err := d.SetCharset(); err != nil {
		return fmt.Errorf("initializing display: setting charset: %w", err)
	}
	if err := d.MoveCursor(0, 0); err != nil {
		return fmt.Errorf("initializing display: moving cursor: %w", err)
	}

	return nil
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
		return "", &DeviceNotFoundError{}
	}

	return path, nil
}

func createPacket(data []uint8) ([]uint8, error) {
	msg := make([]uint8, 0, len(STX)+1+len(data))
	msg = append(msg, STX...)
	msg = append(msg, uint8(len(data)))
	msg = append(msg, data...)

	// report[0] is the HID report ID
	if len(msg) > ReportSize-1 {
		return nil, &PayloadTooLargeError{Size: len(data), Max: ReportSize - 1 - len(STX) - 1}
	}

	report := make([]uint8, ReportSize)
	report[0] = 0x00
	copy(report[1:], msg)
	return report, nil
}

func (d *Display) WriteData(data []uint8) error {
	report, err := createPacket(data)
	if err != nil {
		return err
	}

	if _, err := d.device.Write(report); err != nil {
		return fmt.Errorf("writing to display: %w", err)
	}
	return nil
}
