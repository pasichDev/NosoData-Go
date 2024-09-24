package legacy

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Friends-Of-Noso/NosoData-Go/utils"
)

type LegacyPSO struct {
	Block       int32
	MNLockCount int32
	PSOCount    int32
	MNLocks     []LegacyMNLockItem
	PSOS        []LegacyPSOItem
}

func (p *LegacyPSO) ReadFromFile(f string) error {
	// Check if the file exists before trying to open it
	if !utils.FileExists(f) {
		return fmt.Errorf("file %s not found", f)
	}

	file, err := os.Open(f)
	if err != nil {
		return fmt.Errorf("cannot open file: %s", f)
	}
	defer file.Close()

	// Field Block
	err = binary.Read(file, binary.LittleEndian, &p.Block)
	if err != nil {
		return err
	}

	// Field MNLockCount
	err = binary.Read(file, binary.LittleEndian, &p.MNLockCount)
	if err != nil {
		return err
	}

	// Field PSOCount
	err = binary.Read(file, binary.LittleEndian, &p.PSOCount)
	if err != nil {
		return err
	}

	// Field MNLocks
	if p.MNLockCount > 0 {
		for i := 0; i < int(p.MNLockCount); i++ {
			mli := LegacyMNLockItem{}
			err := mli.ReadFromStream(file)
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			p.MNLocks = append(p.MNLocks, mli)
		}
	}

	// Field PSOS
	// if p.PSOCount > 0 {
	// 	for i := 0; i < int(p.PSOCount); i++ {
	// 		psoi := LegacyPSOItem{}
	// 		err := psoi.ReadFromStream(file)
	// 		if err == io.EOF {
	// 			break
	// 		}
	// 		if err != nil {
	// 			return err
	// 		}
	// 		p.PSOS = append(p.PSOS, psoi)
	// 	}
	// }

	return nil
}

type LegacyMNLockItem struct {
	Address PascalShortString // Capacity 32 aligned makes it 35
	Expire  int32
}

func (m *LegacyMNLockItem) ReadFromStream(f *os.File) error {
	// Check if the stream is nil
	if f == nil {
		return errors.New("nil reader provided")
	}

	// Field Address
	m.Address = *NewPascalShortString(35)
	err := m.Address.ReadFromStream(f)
	if err != nil {
		return err
	}

	// Field Expire
	err = binary.Read(f, binary.LittleEndian, &m.Expire)
	if err != nil {
		return err
	}

	return nil
}

type LegacyPSOItem struct {
	Mode    int32
	Hash    string
	Owner   string
	Expire  int32
	Members string
	Params  string
}
