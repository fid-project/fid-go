// Copyright 2025 Fred Bairn
// Licensed under the Apache License, Version 2.0.

package fid

import (
	"sync"
	"testing"
	"time"
)

// --- Basic generation and parsing ---

func TestRoundTrip(t *testing.T) {
	id, err := New(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s := id.String()

	parsed, err := Parse(s)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if parsed != id {
		t.Fatalf("round-trip mismatch:\n got:  %v\n want: %v", parsed, id)
	}
}

func TestStringLengthAndCheckDigit(t *testing.T) {
	id, err := New(1)
	if err != nil {
		t.Fatal(err)
	}

	s := id.String()
	if len(s) != 14 {
		t.Errorf("unexpected string length: got %d, want 14", len(s))
	}

	// verify that changing the check digit causes a failure
	origCheck := s[len(s)-1]
	var newCheck byte = '0'
	if origCheck == newCheck {
		newCheck = '1'
	}
	bad := s[:len(s)-1] + string(newCheck)

	if _, err := Parse(bad); err == nil {
		t.Error("expected check digit error, got nil")
	}
}

// --- Monotonicity and uniqueness ---

func TestMonotonicity(t *testing.T) {
	prev, _ := New(1)
	for i := 0; i < 1000; i++ {
		id, err := New(1)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if id <= prev {
			t.Fatalf("non-monotonic: %v -> %v", prev, id)
		}
		prev = id
	}
}

func TestConcurrentGeneration(t *testing.T) {
	const goroutines = 32
	const perG = 100
	ch := make(chan ID, goroutines*perG)

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for g := 0; g < goroutines; g++ {
		go func() {
			defer wg.Done()
			for i := 0; i < perG; i++ {
				id, err := New(1)
				if err != nil {
					t.Errorf("New returned error in goroutine: %v", err)
					return
				}
				ch <- id
			}
		}()
	}

	// Close the channel when all producers are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	ids := make(map[ID]struct{})
	for id := range ch {
		if _, exists := ids[id]; exists {
			t.Fatalf("duplicate id detected: %v", id)
		}
		ids[id] = struct{}{}
	}

	expected := goroutines * perG
	if len(ids) != expected {
		t.Fatalf("expected %d ids, got %d", expected, len(ids))
	}
}

// --- Timestamp and type extraction ---

func TestTypeAndTimestamp(t *testing.T) {
	id, _ := New(9)
	if id.Kind() != 9 {
		t.Errorf("unexpected kind: got %d, want 9", id.Kind())
	}
	ts := id.TimestampMS()
	if ts < epochMS {
		t.Errorf("timestamp too early: %d", ts)
	}
	if ts > uint64(time.Now().UnixMilli())+1000 {
		t.Errorf("timestamp too far in future: %d", ts)
	}
}

// --- Benchmarks ---

func BenchmarkFIDNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := New(1); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParse(b *testing.B) {
	id, err := New(1)
	if err != nil {
		b.Fatal(err)
	}

	s := id.String()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Parse(s); err != nil {
			b.Fatal(err)
		}
	}
}
