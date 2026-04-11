package packet

import (
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"io"
)

// Packet represents a packet that can be sent or received over the Portal socket connection.
type Packet interface {
	ID() uint16
	Marshal(w *protocol.Writer)
	Unmarshal(r *protocol.Reader)
}

// Header is the header of a packet. It contains the packet ID encoded as 2 bytes (little-endian).
type Header struct {
	PacketID uint16
}

// Write writes the header to the given writer.
func (header *Header) Write(w io.ByteWriter) error {
	if err := w.WriteByte(byte(header.PacketID)); err != nil {
		return err
	}
	if err := w.WriteByte(byte(header.PacketID >> 8)); err != nil {
		return err
	}
	return nil
}

// Read reads the header from the given reader.
func (header *Header) Read(r io.ByteReader) error {
	b1, err := r.ReadByte()
	if err != nil {
		return err
	}
	b2, err := r.ReadByte()
	if err != nil {
		return err
	}
	header.PacketID = uint16(b1) | uint16(b2)<<8
	return nil
}
