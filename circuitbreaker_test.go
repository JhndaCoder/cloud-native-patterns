package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"
)

func counter() Circuit {
	m := &sync.Mutex{}
	count := 0

	return func(ctx context.Context) (string, error) {
		m.Lock()
		count++
		m.Unlock()

		return fmt.Sprintf("%d", count), nil
	}
}

func failAfter(threshold int) Circuit {
	count := 0
	return func(ctx context.Context) (string, error) {
		count++
		if count > threshold {
			return "", errors.New("INTENTIONAL FAIL!")
		}
		return "SUCCESS", nil
	}
}

func TestCircuitBreakerFailAfter5(t *testing.T) {
	circuit := failAfter(5)
	ctx := context.Background()

	for count := 1; count <= 5; count++ {
		_, err := circuit(ctx)
		t.Logf("attempt %d: %v", count, err)
		switch {
		case count <= 5 && err != nil:
			t.Errorf("expected no error, got %v", err)

		case count > 5 && err == nil:
			t.Errorf("expected error, got none")
		}
	}
}

func TestCircuitBreaker(t *testing.T) {
	circuit := failAfter(5)
	breaker := Breaker(circuit, 1)

	ctx := context.Background()

	circuitOpen := false
	doesCircuitOpen := false
	doesCircuitReclose := false
	count := 0

	for range time.NewTicker(time.Second).C {
		_, err := breaker(ctx)
		if err != nil {
			if strings.HasPrefix(err.Error(), "service unreachable") {
				if !circuitOpen {
					circuitOpen = true
					doesCircuitOpen = true

					t.Log("Circuit has opened")
				} else {
					circuitOpen = false
					doesCircuitReclose = true
					t.Log("Circuit has automatically closed")
				}
			}
		} else {
			t.Log("circuit closed and operational")
		}
		count++
		if count >= 10 {
			break
		}
	}
	if !doesCircuitOpen {
		t.Error("circuit doesn't appear to open")
	}
	if !doesCircuitReclose {
		t.Error("circuit doesn't appear to close after time")
	}
}
