package spvalidator

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func expectNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func expectError(t *testing.T, err error) *ValidationError {
	t.Helper()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("expected ValidationError, got %T: %v", err, err)
	}
	return ve
}

func writeTempFile(t *testing.T, dir, name string, content []byte) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func validateErr(t *testing.T, err error, tag string) {
	t.Helper()
	ve := expectError(t, err)
	if ve.Tag != tag {
		t.Fatalf("expected tag %q, got %q", tag, ve.Tag)
	}
}

func assertFormatError(t *testing.T, err error, tag string) {
	t.Helper()
	validateErr(t, err, tag)
}

func errorString(err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprint(err)
}
