// Copyright 2025 Fred Bairn
// Licensed under the Apache License, Version 2.0.

package fid

import "errors"

var (
	ErrorBadCheckDigit = errors.New("bad check digit")
	ErrorInvalidID     = errors.New("invalid ID format")
	ErrorIdBadLength   = errors.New("bad length")
	ErrorIdOverflow    = errors.New("ID overflow")
)
