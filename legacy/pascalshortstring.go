package legacy

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

type PascalShortString struct {
	data     []byte // Contains raw bytes, including the length and garbage
	length   uint8  // Cached length for quick access (matches the first byte of Data)
	capacity int    // Maximum capacity (e.g., 20 for string[20])
}

// NewPascalShortString creates a new PascalShortString with a given capacity
func NewPascalShortString(capacity int) *PascalShortString {
	return &PascalShortString{
		data:     make([]byte, capacity+1), // Add 1 to hold the length byte
		capacity: capacity,
		length:   0, // Initialize length to 0
	}
}

// ReadFromStream reads a Pascal Short String from the provided stream
func (p *PascalShortString) ReadFromStream(r io.Reader) error {
	// Check if the stream is nil
	if r == nil {
		return errors.New("nil reader provided")
	}

	if len(p.data) != p.capacity+1 {
		p.data = make([]byte, p.capacity+1) // Ensure the internal array matches capacity+1
	}

	// Read the length byte (first byte of Data)
	err := binary.Read(r, binary.LittleEndian, &p.data[0])
	if err != nil {
		return err
	}
	if int(p.data[0]) > p.capacity {
		// return fmt.Errorf("capacity is %d, but read length is %d", p.Capacity, p.Data[0])
		p.data[0] = byte(p.capacity)
	}

	// Update the Length field from the first byte of Data
	p.length = p.data[0]

	// Read the actual data based on the capacity (including garbage)
	n, err := r.Read(p.data[1 : p.capacity+1]) // Read the string data plus garbage
	if n != p.capacity {
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
	n, err := w.Write(p.data[:p.capacity+1])
	if n != p.capacity+1 {
		return errors.New("failed to write the entire string data")
	}
	return err
}

// SetString sets the string value without overwriting existing content, and ensures garbage is maintained
func (p *PascalShortString) SetString(s string) error {
	strBytes := []byte(s)

	// Ensure the length of the string fits within capacity
	if len(strBytes) > p.capacity {
		return errors.New("string exceeds capacity")
	}

	// Set the length byte in both the field and the first byte of Data
	p.length = uint8(len(strBytes))
	p.data[0] = p.length

	// Overwrite only the bytes corresponding to the string, leaving the rest of the Data unchanged
	copy(p.data[1:1+len(strBytes)], strBytes)

	// Ensure the rest of Data beyond p.Length is not touched, preserving garbage
	return nil
}

// GetString returns the actual string part (up to the length byte)
func (p *PascalShortString) GetString() string {
	return string(p.data[1 : 1+int(p.length)])
}

// MarshalJSON allows this struct to be used with `json.Marshal()`
func (p *PascalShortString) MarshalJSON() ([]byte, error) {
	value := p.GetString()
	return json.Marshal(value)
}
