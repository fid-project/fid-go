// Package fid implements FID, a 64-bit, time-ordered identifier format.
//
// FIDs are 8-byte unsigned integers with the following layout:
//
//	[ kind:8 | time:43 | node:8 | counter:5 ]
//
// Use New to generate, Parse to decode from text, and ID.String for
// canonical Crockford Base32 encoding (13 chars + 1 check digit).
//
// Node IDs are derived from FID_NODE_ID or machine identity.
// See spec.md for full details.
package fid
