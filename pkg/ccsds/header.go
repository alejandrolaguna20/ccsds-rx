package ccsds

// SpacePacket is a zero-allocation view of a raw CCSDS packet buffer.
// It is a type alias to []byte to allow receiver methods without copying.
type SpacePacket []byte

const (
	// PrimaryHeaderSize is the fixed 6-byte CCSDS primary header.
	PrimaryHeaderSize = 6
)

// NewSpacePacket creates a SpacePacket from a byte slice.
// This operation is O(1) and performs ZERO allocations or copies.
func NewSpacePacket(data []byte) (SpacePacket, error) {
	if len(data) < PrimaryHeaderSize {
		return nil, ErrPacketTooShort
	}
	p := SpacePacket(data)
	if err := p.Validate(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p SpacePacket) Validate() error {
	if len(p) < PrimaryHeaderSize {
		return ErrPacketTooShort
	}

	if p.Version() != 0 {
		return ErrInvalidVersion
	}

	if len(p) < p.TotalLength() {
		return ErrInvalidLength
	}

	return nil
}

func (p SpacePacket) Version() uint8 {
	return (p[0] >> 5) & 0x07
}

func (p SpacePacket) Type() uint8 {
	return (p[0] >> 4) & 0x01
}

func (p SpacePacket) HasSecondaryHeader() bool {
	return ((p[0] >> 3) & 0x01) == 1
}

func (p SpacePacket) APID() uint16 {
	return (uint16(p[0]&0x07) << 8) | uint16(p[1])
}

func (p SpacePacket) SequenceFlags() uint8 {
	return (p[2] >> 6) & 0x03
}

func (p SpacePacket) SequenceCount() uint16 {
	return (uint16(p[2]&0x3F) << 8) | uint16(p[3])
}

func (p SpacePacket) DataLength() uint16 {
	return (uint16(p[4]) << 8) | uint16(p[5])
}

func (p SpacePacket) TotalLength() int {
	return int(p.DataLength()) + 1 + PrimaryHeaderSize
}

func (p SpacePacket) UserData() []byte {
	return p[PrimaryHeaderSize:p.TotalLength()]
}
