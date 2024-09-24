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

type LegacyPSO struct {
	Block       int32              `json:"block"`
	MNLockCount int32              `json:"mn-locks-count"`
	PSOCount    int32              `json:"psos-count"`
	MNLocks     []LegacyMNLockItem `json:"mn-locks"`
	PSOS        []LegacyPSOItem    `json:"psos"`
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

func (p *LegacyPSO) AsJSON() string {
	jsonData, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		fmt.Printf("error %v", err)
		return ""
	}
	return string(jsonData)
}

type LegacyMNLockItem struct {
	Address PascalShortString `json:"address"` // Capacity 32 aligned makes it 35
	Expire  int32             `json:"expire"`
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
	Mode    int32  `json:"mode"`
	Hash    string `json:"hash"`
	Owner   string `json:"owner"`
	Expire  int32  `json:"expire"`
	Members string `json:"members"`
	Params  string `json:"params"`
}
