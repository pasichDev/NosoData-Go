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

type LegacyWallet struct {
	AccountsCount int64                 `json:"accounts-count"`
	Accounts      []LegacyWalletAccount `json:"accounts"`
}

func (w *LegacyWallet) ReadFromFile(f string) error {
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
		a := LegacyWalletAccount{}
		err := a.ReadFromStream(file)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		w.AccountsCount += 1
		w.Accounts = append(w.Accounts, a)
	}

	return nil
}

type LegacyWalletAccount struct {
	Hash          PascalShortString `json:"hash"`        // Capacity 40
	Custom        PascalShortString `json:"custom"`      // Capacity 40
	PrivateKey    PascalShortString `json:"private-key"` // Capacity 255
	PublicKey     PascalShortString `json:"public-key"`  // Capacity 255
	Balance       int64             `json:"balance"`
	Pending       int64             `json:"pending"`
	Score         int64             `json:"score"`
	LastOperation int64             `json:"last-operation"`
}

func (a *LegacyWalletAccount) ReadFromStream(f *os.File) error {
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

	// Field PublicKey
	a.PublicKey = *NewPascalShortString(255)
	err = a.PublicKey.ReadFromStream(f)
	if err != nil {
		return err
	}

	// Field PrivateKey
	a.PrivateKey = *NewPascalShortString(255)
	err = a.PrivateKey.ReadFromStream(f)
	if err != nil {
		return err
	}

	// Field Balance
	err = binary.Read(f, binary.LittleEndian, &a.Balance)
	if err != nil {
		return err
	}

	// Field Pending
	err = binary.Read(f, binary.LittleEndian, &a.Pending)
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

func (w *LegacyWallet) AsJSON() string {
	jsonData, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		fmt.Printf("error %v", err)
		return ""
	}
	return string(jsonData)
}
