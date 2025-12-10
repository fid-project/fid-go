// Copyright 2025 Fred Bairn
// Licensed under the Apache License, Version 2.0.

package fid

import "strings"

const alphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ"
const checkAlphabet = "0123456789ABCDEFGHJKMNPQRSTVWXYZ*~$=U"

func encode(n uint64) string {
	if n == 0 {
		return "0"
	}
	var buf [13]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = alphabet[n%32]
		n /= 32
	}
	return string(buf[i:])
}

func decode(s string) (uint64, error) {
	var n uint64
	s = normalize(s)

	const max = ^uint64(0)

	for _, r := range s {
		i := indexId(r)
		if i < 0 {
			return 0, ErrorInvalidID
		}
		if n > (max-uint64(i))/32 {
			return 0, ErrorIdOverflow
		}
		n = n*32 + uint64(i)
	}
	return n, nil
}

func computeCheckUint64(v uint64) byte {
	return checkAlphabet[v%37]
}

func normalize(s string) string {
	s = strings.ToUpper(s)
	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch r {
		case '-', '_', ' ', '.':
			continue
		case 'O':
			r = '0'
		case 'I', 'L':
			r = '1'
		}
		b.WriteRune(r)
	}
	return b.String()
}

func indexId(r rune) int {
	switch r {
	case 'O':
		r = '0'
	case 'I', 'L':
		r = '1'
	}
	return strings.IndexRune(alphabet, r)
}
