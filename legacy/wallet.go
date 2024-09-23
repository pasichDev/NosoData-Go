package legacy

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Friends-Of-Noso/NosoData-Go/utils"
)

type LegacyWallet struct {
	AccountsCount int64
	Accounts      []LegacyWalletAccount
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
	Hash          PascalShortString // Capacity 40
	Custom        PascalShortString // Capacity 40
	PrivateKey    PascalShortString // Capacity 255
	PublicKey     PascalShortString // Capacity 255
	Balance       int64
	Pending       int64
	Score         int64
	LastOperation int64
}

func (a *LegacyWalletAccount) ReadFromStream(r io.Reader) error {
	// Check if the stream is nil
	if r == nil {
		return errors.New("nil reader provided")
	}

	// Field Hash
	a.Hash = *NewPascalShortString(40)
	err := a.Hash.ReadFromStream(r)
	if err != nil {
		return err
	}

	// Field Custom
	a.Custom = *NewPascalShortString(40)
	err = a.Custom.ReadFromStream(r)
	if err != nil {
		return err
	}

	// Field PublicKey
	a.PublicKey = *NewPascalShortString(255)
	err = a.PublicKey.ReadFromStream(r)
	if err != nil {
		return err
	}

	// Field PrivateKey
	a.PrivateKey = *NewPascalShortString(255)
	err = a.PrivateKey.ReadFromStream(r)
	if err != nil {
		return err
	}

	// Field Balance
	err = binary.Read(r, binary.LittleEndian, &a.Balance)
	if err != nil {
		return err
	}

	// Field Pending
	err = binary.Read(r, binary.LittleEndian, &a.Pending)
	if err != nil {
		return err
	}

	// Field Score
	err = binary.Read(r, binary.LittleEndian, &a.Score)
	if err != nil {
		return err
	}

	// Field LastOperation
	err = binary.Read(r, binary.LittleEndian, &a.LastOperation)
	if err != nil {
		return err
	}

	return nil
}
