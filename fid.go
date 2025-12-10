// Copyright 2025 Fred Bairn
// Licensed under the Apache License, Version 2.0.

package fid

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

const epochMS = 1747267200000 // 2025-05-15T00:00:00Z in Unix milliseconds

type ID uint64

// bit widths
const (
	kindBits  = 8
	timeBits  = 43
	nodeBits  = 8
	countBits = 5
)

const (
	countMax = (1 << countBits) // Max 32
	timeMax  = (1 << timeBits)
)

var state atomic.Uint64 // state: high bits = ms since epoch, low 5 bits = counter

func New(kindCode uint8) (ID, error) {
	for {
		nowMs := uint64(time.Now().UnixMilli())
		if nowMs < epochMS {
			return 0, fmt.Errorf("time before epoch")
		}
		ms := nowMs - epochMS
		if ms >= timeMax {
			return 0, fmt.Errorf("timestamp overflow")
		}

		cur := state.Load()
		lastMs := cur >> countBits
		lastCount := cur & (countMax - 1)

		if ms < lastMs {
			ms = lastMs
		}

		var nextMs, nextCount uint64
		if ms == lastMs {
			if lastCount+1 >= countMax {
				time.Sleep(time.Millisecond)
				continue
			}
			nextMs = lastMs
			nextCount = lastCount + 1
		} else {
			nextMs = ms
			nextCount = 0
		}

		next := (nextMs << countBits) | nextCount
		if state.CompareAndSwap(cur, next) {
			v := (uint64(kindCode) << (timeBits + nodeBits + countBits)) |
				(nextMs << (nodeBits + countBits)) |
				(uint64(nodeId) << countBits) |
				nextCount
			return ID(v), nil
		}
	}
}

func Parse(s string) (ID, error) {
	s = strings.TrimSpace(s)
	if len(s) < 14 {
		return 0, ErrorIdBadLength
	}

	// split into body (everything except last char) and provided check digit
	body := normalize(s[:len(s)-1])
	if len(body) != 13 {
		return 0, ErrorIdBadLength
	}

	suppliedCheck := rune(strings.ToUpper(s)[len(s)-1])

	// decode the body into the underlying uint64
	value, err := decode(body)
	if err != nil {
		return 0, err
	}

	// recompute the check digit from that number
	computedCheck := rune(computeCheckUint64(value))
	if suppliedCheck != computedCheck {
		return 0, ErrorBadCheckDigit
	}

	// wrap it back up in your ID kind
	return ID(value), nil
}

func (id ID) String() string {
	body := encode(uint64(id))
	if len(body) < 13 {
		body = strings.Repeat("0", 13-len(body)) + body
	}
	return body + string(computeCheckUint64(uint64(id)))
}

func (id ID) Kind() uint8 {
	return uint8(
		(uint64(id) >> (timeBits + nodeBits + countBits)) &
			((1 << kindBits) - 1),
	)
}

func (id ID) Node() uint8 {
	return uint8((uint64(id) >> countBits) & ((1 << nodeBits) - 1))
}

func (id ID) TimestampMS() uint64 {
	ms := (uint64(id) >> (nodeBits + countBits)) & ((1 << timeBits) - 1)
	return ms + epochMS
}

func (id ID) MarshalText() ([]byte, error) { return []byte(id.String()), nil }
func (id *ID) UnmarshalText(b []byte) error {
	x, e := Parse(string(b))
	if e != nil {
		return e
	}
	*id = x
	return nil
}

func (id ID) MarshalJSON() ([]byte, error) { return json.Marshal(id.String()) }
func (id *ID) UnmarshalJSON(b []byte) error {
	var s string
	if e := json.Unmarshal(b, &s); e != nil {
		return e
	}
	x, e := Parse(s)
	if e != nil {
		return e
	}
	*id = x
	return nil
}
