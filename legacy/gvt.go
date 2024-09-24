package legacy

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Friends-Of-Noso/NosoData-Go/utils"
)

type LegacyGVT struct {
	EntryCount int64            `json:"entries-count"`
	Entries    []LegacyGVTEntry `json:"entries"`
}

func (g *LegacyGVT) ReadFromFile(f string) error {
	// Check if the file exists before trying to open it
	if !utils.FileExists(f) {
		return fmt.Errorf("file %s not found", f)
	}

	file, err := os.Open(f)
	if err != nil {
		return fmt.Errorf("cannot open file: %s", f)
	}
	defer file.Close()

	for {
		e := LegacyGVTEntry{}
		err := e.ReadFromStream(file)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		g.EntryCount += 1
		g.Entries = append(g.Entries, e)

	}

	return nil
}

func (g *LegacyGVT) AsJSON() string {
	jsonData, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		fmt.Printf("error %v", err)
		return ""
	}
	return string(jsonData)
}

type LegacyGVTEntry struct {
	Number  PascalShortString `json:"number"` // Capacity 2
	Owner   PascalShortString `json:"owner"`  // Capacity 32
	Hash    PascalShortString `json:"hash"`   // Capacity 64
	Control int32             `json:"control"`
}

func (e *LegacyGVTEntry) ReadFromStream(f *os.File) error {
	// Check if the stream is nil
	if f == nil {
		return errors.New("nil reader provided")
	}

	// Field Number
	e.Number = *NewPascalShortString(2)
	err := e.Number.ReadFromStream(f)
	if err != nil {
		return err
	}

	// Field Owner
	e.Owner = *NewPascalShortString(32)
	err = e.Owner.ReadFromStream(f)
	if err != nil {
		return err
	}

	// Field Hash
	e.Hash = *NewPascalShortString(64)
	err = e.Hash.ReadFromStream(f)
	if err != nil {
		return err
	}

	// Field Control
	err = binary.Read(f, binary.LittleEndian, &e.Control)
	if err != nil {
		return err
	}

	return nil
}
