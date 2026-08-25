package utils

import (
	"testing"
	"time"
)

func TestNewDBContext(t *testing.T) {
	ctx, cancel := NewDBContext()
	defer cancel()

	if ctx == nil {
		t.Fatal("expected context, got nil")
	}

	select {
	case <-ctx.Done():
		t.Fatal("context expired too early")
	default:
	}

	time.Sleep(2100 * time.Millisecond)

	select {
	case <-ctx.Done():
		// esperando
	default:
		t.Error("expected context to be expired")
	}
}