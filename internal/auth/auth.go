package auth

import (
	"errors"
	"net/http"
	"strings"
	"testing"
)

var ErrNoAuthHeaderIncluded = errors.New("no authorization header included")

// GetAPIKey -
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	return splitAuth[1], nil
}

// Testing GetAPIKey function

func TestGetAPIKey(t *testing.T) {
	tests := []struct {
		name       string
		headers    http.Header
		wantKey    string
		wantErr    error  // use this for sentinel errors (e.g., ErrNoAuthHeaderIncluded)
		wantErrStr string // use this for non-sentinel errors created with errors.New(...)
	}{
		{
			name:    "valid header",
			headers: http.Header{"Authorization": {"ApiKey my-secret-key"}},
			wantKey: "my-secret-key",
		},
		{
			name:    "missing Authorization header",
			headers: http.Header{},
			wantErr: ErrNoAuthHeaderIncluded, // compare with sentinel from main package
		},
		{
			name:       "wrong prefix",
			headers:    http.Header{"Authorization": {"Bearer token"}},
			wantErrStr: "malformed authorization header",
		},
		{
			name:       "missing key after ApiKey",
			headers:    http.Header{"Authorization": {"ApiKey"}},
			wantErrStr: "malformed authorization header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := GetAPIKey(tt.headers)

			// key check
			if key != tt.wantKey {
				t.Fatalf("key: got %q, want %q", key, tt.wantKey)
			}

			// error checks
			switch {
			case tt.wantErr != nil: // sentinel error
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error: got %v, want %v (sentinel)", err, tt.wantErr)
				}
			case tt.wantErrStr != "": // plain error text
				if err == nil || err.Error() != tt.wantErrStr {
					t.Fatalf("error: got %v, want %q", err, tt.wantErrStr)
				}
			default: // expect no error
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}
