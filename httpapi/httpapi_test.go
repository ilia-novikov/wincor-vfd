package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilia-novikov/wincor-vfd/display"
)

type fakeDevice struct {
	calls  []string
	opErr  error
	closed bool

	text string
	x, y int
}

func (f *fakeDevice) PrintText(text string) error {
	f.calls = append(f.calls, "print")
	f.text = text
	return f.opErr
}

func (f *fakeDevice) ClearDisplay() error {
	f.calls = append(f.calls, "clear")
	return f.opErr
}

func (f *fakeDevice) Initialize() error {
	f.calls = append(f.calls, "init")
	return f.opErr
}

func (f *fakeDevice) MoveCursor(x, y int) error {
	f.calls = append(f.calls, "cursor")
	f.x, f.y = x, y
	return f.opErr
}

func (f *fakeDevice) SendBackspace() error {
	f.calls = append(f.calls, "backspace")
	return f.opErr
}

func (f *fakeDevice) SendLineFeed() error {
	f.calls = append(f.calls, "linefeed")
	return f.opErr
}

func (f *fakeDevice) SendCarriageReturn() error {
	f.calls = append(f.calls, "carriage-return")
	return f.opErr
}

func (f *fakeDevice) DeleteToEndOfLine() error {
	f.calls = append(f.calls, "delete-to-eol")
	return f.opErr
}

func (f *fakeDevice) Close() error {
	f.closed = true
	return nil
}

func openerFor(dev *fakeDevice, openErr error) OpenFunc {
	return func(charset display.Charset) (Device, error) {
		if openErr != nil {
			return nil, openErr
		}
		dev.calls = append(dev.calls, "open:"+charsetName(charset))
		return dev, nil
	}
}

func charsetName(c display.Charset) string {
	switch c {
	case display.CharsetUS:
		return "us"
	case display.CharsetRU:
		return "ru"
	default:
		return "unknown"
	}
}

func okProbe() error { return nil }

func doRequest(t *testing.T, h http.Handler, method, target string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var reader *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshaling body: %v", err)
		}
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, target, reader)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestHealthzAlwaysOK(t *testing.T) {
	h := newServer("", openerFor(&fakeDevice{}, nil), func() error { t.Fatal("probe should not be called"); return nil })

	rec := doRequest(t, h, http.MethodGet, "/healthz", nil)
	if rec.Code != http.StatusOK {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestReadyz(t *testing.T) {
	tests := []struct {
		name       string
		probe      ProbeFunc
		wantStatus int
	}{
		{"present", okProbe, http.StatusOK},
		{"absent", func() error { return &display.DeviceNotFoundError{} }, http.StatusServiceUnavailable},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newServer("", openerFor(&fakeDevice{}, nil), tt.probe)
			rec := doRequest(t, h, http.MethodGet, "/readyz", nil)
			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d", rec.Code, tt.wantStatus)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	dev := &fakeDevice{}
	h := newServer("", openerFor(dev, nil), okProbe)

	rec := doRequest(t, h, http.MethodPost, "/v1/print", map[string]string{"text": "hello"})

	if rec.Code != http.StatusNoContent {
		t.Fatalf("got status %d, want %d, body=%s", rec.Code, http.StatusNoContent, rec.Body)
	}
	if dev.text != "hello" {
		t.Errorf("got printed text %q, want %q", dev.text, "hello")
	}
	if !dev.closed {
		t.Error("device was not closed")
	}
}

func TestPrintEmptyTextIsBadRequest(t *testing.T) {
	h := newServer("", openerFor(&fakeDevice{}, nil), okProbe)

	rec := doRequest(t, h, http.MethodPost, "/v1/print", map[string]string{"text": ""})
	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestPrintInvalidJSON(t *testing.T) {
	h := newServer("", openerFor(&fakeDevice{}, nil), okProbe)

	req := httptest.NewRequest(http.MethodPost, "/v1/print", bytes.NewBufferString("not json"))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestCursor(t *testing.T) {
	dev := &fakeDevice{}
	h := newServer("", openerFor(dev, nil), okProbe)

	rec := doRequest(t, h, http.MethodPut, "/v1/cursor", map[string]int{"x": 3, "y": 5})

	if rec.Code != http.StatusNoContent {
		t.Fatalf("got status %d, want %d, body=%s", rec.Code, http.StatusNoContent, rec.Body)
	}
	if dev.x != 3 || dev.y != 5 {
		t.Errorf("got cursor (%d, %d), want (3, 5)", dev.x, dev.y)
	}
}

func TestSimpleOps(t *testing.T) {
	tests := []struct {
		method, path, wantCall string
	}{
		{http.MethodPost, "/v1/clear", "clear"},
		{http.MethodPost, "/v1/init", "init"},
		{http.MethodPost, "/v1/backspace", "backspace"},
		{http.MethodPost, "/v1/linefeed", "linefeed"},
		{http.MethodPost, "/v1/carriage-return", "carriage-return"},
		{http.MethodPost, "/v1/delete-to-eol", "delete-to-eol"},
	}

	for _, tt := range tests {
		t.Run(tt.wantCall, func(t *testing.T) {
			dev := &fakeDevice{}
			h := newServer("", openerFor(dev, nil), okProbe)

			rec := doRequest(t, h, tt.method, tt.path, nil)

			if rec.Code != http.StatusNoContent {
				t.Fatalf("got status %d, want %d, body=%s", rec.Code, http.StatusNoContent, rec.Body)
			}
			if len(dev.calls) == 0 || dev.calls[len(dev.calls)-1] != tt.wantCall {
				t.Errorf("got calls %v, want last call %q", dev.calls, tt.wantCall)
			}
		})
	}
}

func TestCharsetQueryParam(t *testing.T) {
	t.Run("valid ru charset is passed through", func(t *testing.T) {
		dev := &fakeDevice{}
		h := newServer("", openerFor(dev, nil), okProbe)

		rec := doRequest(t, h, http.MethodPost, "/v1/clear?charset=ru", nil)

		if rec.Code != http.StatusNoContent {
			t.Fatalf("got status %d, want %d, body=%s", rec.Code, http.StatusNoContent, rec.Body)
		}
		if len(dev.calls) == 0 || dev.calls[0] != "open:ru" {
			t.Errorf("got calls %v, want first call %q", dev.calls, "open:ru")
		}
	})

	t.Run("invalid charset is rejected", func(t *testing.T) {
		h := newServer("", openerFor(&fakeDevice{}, nil), okProbe)

		rec := doRequest(t, h, http.MethodPost, "/v1/clear?charset=fr", nil)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusBadRequest)
		}
	})
}

func TestDeviceErrorMapping(t *testing.T) {
	tests := []struct {
		name       string
		openErr    error
		opErr      error
		wantStatus int
	}{
		{"device not found on open", &display.DeviceNotFoundError{}, nil, http.StatusServiceUnavailable},
		{"payload too large on write", nil, &display.PayloadTooLargeError{Size: 100, Max: 10}, http.StatusBadRequest},
		{"generic write failure", nil, errAny, http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dev := &fakeDevice{opErr: tt.opErr}
			h := newServer("", openerFor(dev, tt.openErr), okProbe)

			rec := doRequest(t, h, http.MethodPost, "/v1/clear", nil)
			if rec.Code != tt.wantStatus {
				t.Errorf("got status %d, want %d, body=%s", rec.Code, tt.wantStatus, rec.Body)
			}
		})
	}
}

func TestDeviceClosedEvenOnOpError(t *testing.T) {
	dev := &fakeDevice{opErr: errAny}
	h := newServer("", openerFor(dev, nil), okProbe)

	doRequest(t, h, http.MethodPost, "/v1/clear", nil)

	if !dev.closed {
		t.Error("device was not closed after an operation error")
	}
}

func TestAuth(t *testing.T) {
	dev := &fakeDevice{}
	h := newServer("secret", openerFor(dev, nil), okProbe)

	t.Run("missing token is rejected", func(t *testing.T) {
		rec := doRequest(t, h, http.MethodPost, "/v1/clear", nil)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("wrong token is rejected", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/clear", nil)
		req.Header.Set("Authorization", "Bearer wrong")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusUnauthorized {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusUnauthorized)
		}
	})

	t.Run("correct token is accepted", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/v1/clear", nil)
		req.Header.Set("Authorization", "Bearer secret")
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)

		if rec.Code != http.StatusNoContent {
			t.Errorf("got status %d, want %d, body=%s", rec.Code, http.StatusNoContent, rec.Body)
		}
	})

	t.Run("health endpoints bypass auth", func(t *testing.T) {
		rec := doRequest(t, h, http.MethodGet, "/healthz", nil)
		if rec.Code != http.StatusOK {
			t.Errorf("got status %d, want %d", rec.Code, http.StatusOK)
		}
	})
}

type sentinelError struct{}

func (sentinelError) Error() string { return "boom" }

var errAny = sentinelError{}
