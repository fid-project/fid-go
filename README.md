# FID for Go

This module implements the FID spec: a monotonic, 64-bit identifier format designed for distributed systems that need sortable IDs, optional categorisation, and resilient textual encoding.

The spec can be found at <https://github.com/fid-project/fid-spec>

## Getting Started

```
go get github.com/fid-project/fid-go
```

Use `fid.New` to generate IDs and `fid.Parse` to round-trip from strings:

```go
package main

import (
	"fmt"

	"github.com/fid-project/fid-go"
)

func main() {
	// 42 is the "kind" slot – pick any value from 0..255 to tag a domain.
	id, err := fid.New(42)
	if err != nil {
		panic(err)
	}

	fmt.Println("binary:", uint64(id))
	fmt.Println("string:", id.String()) // 14-char canonical Crockford form

	if parsed, err := fid.Parse(id.String()); err == nil {
		fmt.Printf("kind=%d node=%d timestamp=%d\n",
			parsed.Kind(), parsed.Node(), parsed.TimestampMS())
	}
}
```

`ID` implements `fmt.Stringer`, `encoding.TextMarshaler`, and `encoding/json` interfaces so IDs survive through logs, HTTP payloads, and configuration files without extra wiring.

## Node identity

ID uniqueness within the same millisecond relies on the node slot. At startup the package runs `InitNodeID`, which picks an 8-bit value using the following precedence:

1. `FID_NODE_ID=<0..255>` – use the numeric value directly.
2. `FID_NODE_ID=<string>` – hash arbitrary strings to fit the slot.
3. `/etc/machine-id` – hash the host identifier (Linux default).
4. Machine hostname hash.
5. As a last resort, hash PIDs and UID to differentiate local processes.

If you need deterministic values per deployment, set `FID_NODE_ID` via the environment, a launch flag, or call `fid.InitNodeID()` manually after preparing the environment.

## Text encoding rules

- Input strings accept lowercase, uppercase, and separators (`-`, `_`, `.`, space`).
- Ambiguous characters map to their safer equivalents (`O → 0`, `I/L → 1`).
- The decoder enforces the 13-character body; mis-sized IDs fail with `ErrorIdBadLength`.
- Check digits are validated with `ErrorBadCheckDigit`.

## Testing and benchmarks

```
go test ./...
```

The test suite exercises round-trips, monotonic generation across goroutines, kind/timestamp extraction, and includes micro-benchmarks for `fid.New` and `fid.Parse`.

## Spec alignment

- 64-bit unsigned IDs with monotonic ordering (per FID spec).
- Crockford Base32 text form plus 1-character checksum.
- User-defined kind byte and node byte for sharding or routing hints.
- Counter rollover logic sleeps for the next millisecond when ~32 IDs/ms/node is exceeded.
- Errors surface as typed Go `error` values (`ErrorBadCheckDigit`, `ErrorInvalidID`, etc.) for easy handling.

See `doc.go` for short package commentary or wire the API into your services and libraries directly.
