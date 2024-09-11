package internal

import (
	"bytes"
	"errors"
	"log"
	"strings"
	"testing"
)

// captureOutput redirects log's output to buffer and returns it.
func captureOutput(f func()) string {
    var buf bytes.Buffer
    log.SetOutput(&buf)
    f()
    log.SetOutput(nil) // Reset output to default
    return buf.String()
}

// TestLogError tests the LogError function.
func TestLogError(t *testing.T) {
    tests := []struct {
        name    string
        err     error
        wantLog bool
    }{
        {"WithError", errors.New("test error"), true},
        {"WithNil", nil, false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            output := captureOutput(func() {
                LogError(tt.err)
            })

            if tt.wantLog && !strings.Contains(output, tt.err.Error()) {
                t.Errorf("LogError() = %v, want %v", output, tt.err.Error())
            }

            if !tt.wantLog && output != "" {
                t.Errorf("LogError() logged something, but expected nothing")
            }
        })
    }
}