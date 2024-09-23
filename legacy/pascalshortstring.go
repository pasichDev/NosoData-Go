package legacy

import (
	"encoding/binary"
	"errors"
	"io"
)

type PascalShortString struct {
	Data     []byte // Contains raw bytes, including the length and garbage
	Length   uint8  // Cached length for quick access (matches the first byte of Data)
	Capacity int    // Maximum capacity (e.g., 20 for string[20])
}

// NewPascalShortString creates a new PascalShortString with a given capacity
func NewPascalShortString(capacity int) *PascalShortString {
	return &PascalShortString{
		Data:     make([]byte, capacity+1), // Add 1 to hold the length byte
		Capacity: capacity,
		Length:   0, // Initialize length to 0
	}
}

// ReadFromStream reads a Pascal Short String from the provided stream
func (p *PascalShortString) ReadFromStream(r io.Reader) error {
	// Check if the stream is nil
	if r == nil {
		return errors.New("nil reader provided")
	}

	if len(p.Data) != p.Capacity+1 {
		p.Data = make([]byte, p.Capacity+1) // Ensure the internal array matches capacity+1
	}

	// Read the length byte (first byte of Data)
	err := binary.Read(r, binary.LittleEndian, &p.Data[0])
	if err != nil {
		return err
	}
	if int(p.Data[0]) > p.Capacity {
		// return fmt.Errorf("capacity is %d, but read length is %d", p.Capacity, p.Data[0])
		p.Data[0] = byte(p.Capacity)
	}

	// Update the Length field from the first byte of Data
	p.Length = p.Data[0]

	// Read the actual data based on the capacity (including garbage)
	n, err := r.Read(p.Data[1 : p.Capacity+1]) // Read the string data plus garbage
	if n != p.Capacity {
		return errors.New("failed to read the entire string data")
	}

	return err
}

// WriteToStream writes the Pascal Short String to the provided stream
func (p *PascalShortString) WriteToStream(w io.Writer) error {
	// Check if the stream is nil
	if w == nil {
		return errors.New("nil writer provided")
	}

	// Write the entire byte slice, including the length and garbage bytes
	n, err := w.Write(p.Data[:p.Capacity+1])
	if n != p.Capacity+1 {
		return errors.New("failed to write the entire string data")
	}
	return err
}

// SetString sets the string value without overwriting existing content, and ensures garbage is maintained
func (p *PascalShortString) SetString(s string) error {
	strBytes := []byte(s)

	// Ensure the length of the string fits within capacity
	if len(strBytes) > p.Capacity {
		return errors.New("string exceeds capacity")
	}

	// Set the length byte in both the field and the first byte of Data
	p.Length = uint8(len(strBytes))
	p.Data[0] = p.Length

	// Overwrite only the bytes corresponding to the string, leaving the rest of the Data unchanged
	copy(p.Data[1:1+len(strBytes)], strBytes)

	// Ensure the rest of Data beyond p.Length is not touched, preserving garbage
	return nil
}

// GetString returns the actual string part (up to the length byte)
func (p *PascalShortString) GetString() string {
	return string(p.Data[1 : 1+int(p.Length)])
}
