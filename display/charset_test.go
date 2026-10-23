package display

import (
	"bytes"
	"errors"
	"testing"
)

func TestEncodeText(t *testing.T) {
	tests := []struct {
		name    string
		charset Charset
		text    string
		want    []uint8
	}{
		{
			name:    "US ASCII",
			charset: CharsetUS,
			text:    "Hello!",
			want:    []uint8("Hello!"),
		},
		{
			name:    "RU CP866",
			charset: CharsetRU,
			text:    "Привет",
			want:    []uint8{0x8f, 0xe0, 0xa8, 0xa2, 0xa5, 0xe2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encodeText(tt.charset, tt.text)
			if err != nil {
				t.Fatalf("encodeText() error = %v", err)
			}
			if !bytes.Equal(got, tt.want) {
				t.Errorf("encodeText() = %x, want %x", got, tt.want)
			}
		})
	}
}

func TestEncodeTextUnsupportedCharset(t *testing.T) {
	_, err := encodeText(Charset(0xFF), "text")

	var charsetErr *CharsetError
	if !errors.As(err, &charsetErr) {
		t.Fatalf("encodeText() error = %v, want *CharsetError", err)
	}
}
