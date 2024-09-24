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

type LegacySummary struct {
	AccountsCount int64                  `json:"accounts-count"`
	Accounts      []LegacySummaryAccount `json:"accounts"`
}

func (s *LegacySummary) ReadFromFile(f string) error {
	// Check if the file exists before trying to open it
	if !utils.FileExists(f) {
		return fmt.Errorf("file %s not found", f)
	}

	file, err := os.Open(f)
	if err != nil {
		return fmt.Errorf("cannot open file: %s", err)
	}
	defer file.Close()

	for {
		a := LegacySummaryAccount{}
		err := a.ReadFromStream(file)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		s.AccountsCount += 1
		s.Accounts = append(s.Accounts, a)
	}

	return nil
}

func (s *LegacySummary) AsJSON() string {
	jsonData, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		fmt.Printf("error %v", err)
		return ""
	}
	return string(jsonData)
}

type LegacySummaryAccount struct {
	Hash          PascalShortString `json:"hash"`   // Capacity 40
	Custom        PascalShortString `json:"custom"` // Capacity 40
	Balance       int64             `json:"balance"`
	Score         int64             `json:"score"`
	LastOperation int64             `json:"last-operation"`
}

func (a *LegacySummaryAccount) ReadFromStream(f *os.File) error {
	// Check if the stream is nil
	if f == nil {
		return errors.New("nil reader provided")
	}

	// Field Hash
	a.Hash = *NewPascalShortString(40)
	err := a.Hash.ReadFromStream(f)
	if err != nil {
		return err
	}

	// Field Custom
	a.Custom = *NewPascalShortString(40)
	err = a.Custom.ReadFromStream(f)
	if err != nil {
		return err
	}

	// Field Balance
	err = binary.Read(f, binary.LittleEndian, &a.Balance)
	if err != nil {
		return err
	}

	// Field Score
	err = binary.Read(f, binary.LittleEndian, &a.Score)
	if err != nil {
		return err
	}

	// Field LastOperation
	err = binary.Read(f, binary.LittleEndian, &a.LastOperation)
	if err != nil {
		return err
	}

	return nil
}
