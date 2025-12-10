// Copyright 2025 Fred Bairn
// Licensed under the Apache License, Version 2.0.

package fid

import (
	"encoding/binary"
	"os"
	"strconv"
)

// 0–255 node id, set at startup (env, flag, whatever your orchestrator provides)
var nodeId uint8

// InitNodeID re-evaluates the nodeId from environment or machine info.
// It's called automatically at package init but also testable manually.
func InitNodeID() {
	env := os.Getenv("FID_NODE_ID")

	switch {
	case env == "":
		nodeId = deriveNodeID()

	case isNumeric(env):
		v, _ := strconv.Atoi(env)
		if v >= 0 && v <= 255 {
			nodeId = uint8(v)
		} else {
			nodeId = fnv1a8([]byte(env))
		}

	default:
		nodeId = fnv1a8([]byte(env))
	}
}

func init() {
	InitNodeID()
}

func NodeID() uint8 {
	return nodeId
}

func deriveNodeID() uint8 {
	// Try /etc/machine-id first (Linux standard)
	if data, err := os.ReadFile("/etc/machine-id"); err == nil {
		return uint8(fnv1a8(data))
	}

	// Fallback: hostname hash
	if h, err := os.Hostname(); err == nil {
		return uint8(fnv1a8([]byte(h)))
	}

	// Absolute last resort: random fallback (still deterministic during process)
	// (ensures distinct processes on same host get different IDs)
	now := make([]byte, 8)
	binary.LittleEndian.PutUint64(now, uint64(os.Getpid())^uint64(os.Getppid())^uint64(os.Getuid()))
	return fnv1a8(now)
}

func fnv1a8(data []byte) uint8 {
	var hash uint32 = 2166136261
	for _, b := range data {
		hash ^= uint32(b)
		hash *= 16777619
	}
	return uint8(hash)
}

// isNumeric returns true if s represents a base-10 integer.
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
