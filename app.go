package rtcp

import (
	"encoding/binary"
	"errors"
)

const (
	nameLength = 4

	// FIXME: Just hard code SSRC for testing.
	debugSSRC uint32 = 3572843026
)

type App []byte

func (a *App) DestinationSSRC() []uint32 {
	return []uint32{debugSSRC}
}

func (a *App) Marshal() ([]byte, error) {
	rawPacket := make([]byte, a.len())
	packetBody := rawPacket[headerLength:]

	binary.BigEndian.PutUint32(packetBody, debugSSRC)

	// Copy the data to packet body.
	copy(packetBody[ssrcLength+nameLength:], *a)

	rawHeader, err := a.Header().Marshal()
	if err != nil {
		return nil, err
	}
	copy(rawPacket, rawHeader)

	return rawPacket, nil
}

func (a *App) Unmarshal(rawPacket []byte) error {
	var h Header
	if err := h.Unmarshal(rawPacket); err != nil {
		return err
	}

	if h.Type != TypeApplicationDefined {
		return errors.New("not an app packet")
	}

	data := rawPacket[headerLength+ssrcLength+nameLength:]
	data = data[:len(data)-int(h.Count)] // trim padding at the end

	*a = data

	return nil
}

func (a *App) Header() Header {
	return Header{
		// A hacky way of sending padding information to the other end.
		// So that, the other end can trim the padding to get the original byte string.
		Count: uint8(getPadding(len(*a))),

		Type:   TypeApplicationDefined,
		Length: uint16(a.len()/4 - 1),
	}
}

func (a *App) len() int {
	l := headerLength + ssrcLength + nameLength + len(*a)
	return l + getPadding(l)
}

var _ Packet = (*App)(nil)
