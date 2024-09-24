package legacy

import (
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/Friends-Of-Noso/NosoData-Go/utils"
)

const (
	cBlockWithPoS      int64 = 8425
	cBlockWithMNandPoS int64 = 48010
	cBlockWithMNOnly   int64 = 88500
)

type LegacyBlock struct {
	Number                      int64
	HASH                        string
	TimeStart                   int64
	TimeEnd                     int64
	TimeTotal                   int32
	TimeLast20                  int32
	TransactionsCount           int32
	Difficulty                  int32
	TargetHash                  PascalShortString // Capacity 32
	Solution                    PascalShortString // Capacity 200
	LastBlockHash               PascalShortString // Capacity 32
	NextBlockDifficulty         int32
	Miner                       PascalShortString // Capacity 40
	Fee                         int64
	Reward                      int64
	Transactions                []LegacyTransaction
	ProofOfStakeRewardCount     int32
	ProofOfStakeRewardAmount    int64
	ProofOfStakeRewardAddresses []PascalShortString
	MasterNodeRewardCount       int32
	MasterNodeRewardAmount      int64
	MasterNodeRewardAddresses   []PascalShortString
}

// TODO: Implement NewLegacyBlock function/constructor

// ReadFromFile reads the data from a block file
func (b *LegacyBlock) ReadFromFile(f string) error {
	// Check if the file exists before trying to open it
	if !utils.FileExists(f) {
		return fmt.Errorf("file %s not found", f)
	}

	file, err := os.Open(f)
	if err != nil {
		return fmt.Errorf("cannot open file: %s", err)
	}
	defer file.Close()

	return b.ReadFromStream(file)
}

// ReadFromFile reads the data from a stream
func (b *LegacyBlock) ReadFromStream(f *os.File) error {
	// Check if the stream is nil
	if f == nil {
		return errors.New("nil reader provided")
	}

	// Field HASH
	// Create a new MD5 hash object
	hash := md5.New()

	// Copy the file's content into the hash
	if _, err := io.Copy(hash, f); err != nil {
		log.Fatal(err)
	}

	// Get the final hash sum
	hashInBytes := hash.Sum(nil)

	// Convert the hash to a hexadecimal string and then to uppercase
	b.HASH = strings.ToUpper(fmt.Sprintf("%x", hashInBytes))

	// Seek back to the beginning of the file if you need to process it again
	_, err := f.Seek(0, 0)
	if err != nil {
		log.Fatal(err)
	}

	// Field Number
	err = binary.Read(f, binary.LittleEndian, &b.Number)
	if err != nil {
		return err
	}

	// Field TimeStart
	err = binary.Read(f, binary.LittleEndian, &b.TimeStart)
	if err != nil {
		return err
	}

	// Field TimeEnd
	err = binary.Read(f, binary.LittleEndian, &b.TimeEnd)
	if err != nil {
		return err
	}

	// Field TimeTotal
	err = binary.Read(f, binary.LittleEndian, &b.TimeTotal)
	if err != nil {
		return err
	}

	// Field TimeLast
	err = binary.Read(f, binary.LittleEndian, &b.TimeLast20)
	if err != nil {
		return err
	}

	// Field TransactionsCount
	err = binary.Read(f, binary.LittleEndian, &b.TransactionsCount)
	if err != nil {
		return err
	}

	// Field Difficulty
	err = binary.Read(f, binary.LittleEndian, &b.Difficulty)
	if err != nil {
		return err
	}

	// Field TargetHash
	b.TargetHash = *NewPascalShortString(32)
	err = b.TargetHash.ReadFromStream(f)
	if err != nil {
		return err
	}

	// Field Solution
	b.Solution = *NewPascalShortString(200)
	err = b.Solution.ReadFromStream(f)
	if err != nil {
		return err
	}

	// Field LastBlockHash
	b.LastBlockHash = *NewPascalShortString(32)
	err = b.LastBlockHash.ReadFromStream(f)
	if err != nil {
		return err
	}

	// Field NextBlockDifficulty
	err = binary.Read(f, binary.LittleEndian, &b.NextBlockDifficulty)
	if err != nil {
		return err
	}

	// Field Miner
	b.Miner = *NewPascalShortString(40)
	err = b.Miner.ReadFromStream(f)
	if err != nil {
		return err
	}

	// Field Fee
	err = binary.Read(f, binary.LittleEndian, &b.Fee)
	if err != nil {
		return err
	}

	// Field Reward
	err = binary.Read(f, binary.LittleEndian, &b.Reward)
	if err != nil {
		return err
	}

	// Read transactions
	if b.TransactionsCount > 0 {
		b.Transactions = make([]LegacyTransaction, b.TransactionsCount)
		var n int32
		for n = 0; n < b.TransactionsCount; n++ {
			b.Transactions[n].ReadFromStream(f)
		}
	}

	// Read PoS/MN rewards

	// Load PoS rewards
	if b.Number > cBlockWithPoS {

		// Field ProofOfStakeRewardAmount
		err = binary.Read(f, binary.LittleEndian, &b.ProofOfStakeRewardAmount)
		if err != nil {
			return err
		}

		// Field ProofOfStakeRewardCount
		err := binary.Read(f, binary.LittleEndian, &b.ProofOfStakeRewardCount)
		if err != nil {
			return err
		}

		// Field ProofOfStakeRewardAddresses
		if b.ProofOfStakeRewardCount > 0 {
			b.ProofOfStakeRewardAddresses = make([]PascalShortString, b.ProofOfStakeRewardCount)
			var n int32
			for n = 0; n < b.ProofOfStakeRewardCount; n++ {
				b.ProofOfStakeRewardAddresses[n] = *NewPascalShortString(32)
				err := b.ProofOfStakeRewardAddresses[n].ReadFromStream(f)
				if err != nil {
					return err
				}
			}
		}
	}

	// Load MN rewards
	if b.Number > cBlockWithMNandPoS {

		// Field MasterNodeRewardAmount
		err = binary.Read(f, binary.LittleEndian, &b.MasterNodeRewardAmount)
		if err != nil {
			return err
		}

		// Field MasterNodeRewardCount
		err := binary.Read(f, binary.LittleEndian, &b.MasterNodeRewardCount)
		if err != nil {
			return err
		}

		// Field ProofOfStakeRewardAddresses
		if b.MasterNodeRewardCount > 0 {
			b.MasterNodeRewardAddresses = make([]PascalShortString, b.MasterNodeRewardCount)
			var n int32
			for n = 0; n < b.MasterNodeRewardCount; n++ {
				b.MasterNodeRewardAddresses[n] = *NewPascalShortString(32)
				err := b.MasterNodeRewardAddresses[n].ReadFromStream(f)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
