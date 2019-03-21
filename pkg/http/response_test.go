package http

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResponse(t *testing.T) {
	rr := httptest.NewRecorder()

	body := struct{}{}
	code := http.StatusOK
	JSONResponse(rr, body, code)

	// Expect Content-Type to be "application/json"
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type json, got %s\n", contentType)
	}

	// Expect code to be written into the header
	if rr.Code != code {
		t.Errorf("Expected code %d, got %d\n", code, rr.Code)
	}

	// Expect body to not be empty
	if len(rr.Body.Bytes()) == 0 {
		t.Errorf("Body was empty, expected: %v\n", body)
	}
}
