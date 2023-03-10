package rtcp

import (
	"encoding/binary"
	"errors"
)

const (
	nameLength = 4
)

type App struct {
	SSRC uint32
	D    []byte
}

func (a *App) DestinationSSRC() []uint32 {
	return []uint32{a.SSRC}
}

func (a *App) Marshal() ([]byte, error) {
	rawPacket := make([]byte, a.len())
	packetBody := rawPacket[headerLength:]

	binary.BigEndian.PutUint32(packetBody, a.SSRC)

	// Copy the data to packet body.
	copy(packetBody[ssrcLength+nameLength:], a.D)

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
		return errors.New("not a valid app packet")
	}

	packetBody := rawPacket[headerLength:]

	a.SSRC = binary.BigEndian.Uint32(packetBody)
	a.D = packetBody[ssrcLength+nameLength : len(packetBody)-int(h.Count)]

	return nil
}

func (a *App) Header() Header {
	return Header{
		// A hacky way of sending padding information to the other end.
		// So that, the other end can trim the padding to get the original byte string.
		Count: uint8(getPadding(len(a.D))),

		Type:   TypeApplicationDefined,
		Length: uint16(a.len()/4 - 1),
	}
}

func (a *App) len() int {
	l := headerLength + ssrcLength + nameLength + len(a.D)
	return l + getPadding(l)
}

var _ Packet = (*App)(nil)
