package main

import (
	"testing"

	"github.com/ilia-novikov/wincor-vfd/display"
)

func TestParseCharset(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want display.Charset
	}{
		{name: "us", in: "us", want: display.CharsetUS},
		{name: "ru", in: "ru", want: display.CharsetRU},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCharset(tt.in)
			if err != nil {
				t.Fatalf("parseCharset() error = %v", err)
			}
			if got != tt.want {
				t.Errorf("parseCharset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseCharsetInvalid(t *testing.T) {
	if _, err := parseCharset("fr"); err == nil {
		t.Fatal("parseCharset() error = nil, want error")
	}
}
