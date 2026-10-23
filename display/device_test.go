package display

import (
	"bytes"
	"errors"
	"testing"
)

func TestCreatePacket(t *testing.T) {
	data := []uint8{0x1B, 0x52, 0x00}

	report, err := createPacket(data)
	if err != nil {
		t.Fatalf("createPacket() error = %v", err)
	}

	if len(report) != ReportSize {
		t.Fatalf("len(report) = %d, want %d", len(report), ReportSize)
	}

	want := []uint8{0x00, 0x02, 0x00, 0x03, 0x1B, 0x52, 0x00}
	if !bytes.Equal(report[:len(want)], want) {
		t.Errorf("report header = %x, want %x", report[:len(want)], want)
	}

	if !bytes.Equal(report[len(want):], make([]uint8, ReportSize-len(want))) {
		t.Errorf("report padding is not zeroed: %x", report[len(want):])
	}
}

func TestCreatePacketTooLarge(t *testing.T) {
	data := make([]uint8, ReportSize)

	_, err := createPacket(data)

	var sizeErr *PayloadTooLargeError
	if !errors.As(err, &sizeErr) {
		t.Fatalf("createPacket() error = %v, want *PayloadTooLargeError", err)
	}
}
