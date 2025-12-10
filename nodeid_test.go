// Copyright 2025 Fred Bairn
// Licensed under the Apache License, Version 2.0.

package fid

import (
	"os"
	"testing"
)

// --- Node ID detection and override ---

func TestNodeIDOverride(t *testing.T) {
	const expected = 123
	os.Setenv("FID_NODE_ID", "123")
	defer os.Unsetenv("FID_NODE_ID")

	InitNodeID() // re-run init logic

	if got := NodeID(); got != expected {
		t.Fatalf("nodeID override failed: got %d, want %d", got, expected)
	}
}

func TestNodeIDFallback(t *testing.T) {
	os.Unsetenv("FID_NODE_ID")

	InitNodeID() // re-derive node id

	if NodeID() == 0 {
		t.Error("expected nonzero nodeID from fallback derivation")
	}
}
