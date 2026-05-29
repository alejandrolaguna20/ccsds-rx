package ccsds

import "errors"

var (
	ErrPacketTooShort = errors.New("packet buffer too short: minimum 6 bytes required")
	ErrInvalidVersion = errors.New("invalid packet version: expected 000")
	ErrInvalidLength  = errors.New("packet data length mismatch: exceeds buffer capacity")
)
