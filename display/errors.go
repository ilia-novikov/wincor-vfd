package display

import "fmt"

type CharsetError struct {
	Charset Charset
}

func (e *CharsetError) Error() string {
	return fmt.Sprintf("only US and RU charsets are supported, got %q", e.Charset)
}
