package legacy

import (
	"encoding/binary"
	"errors"
	"io"
)

type LegacyTransaction struct {
	Block           int32
	OrderID         PascalShortString // Capacity 64
	OrderLinesCount int32
	OrderType       PascalShortString // Capacity 6
	TimeStamp       int64
	Reference       PascalShortString // Capacity 64
	TransferIndex   int32
	Sender          PascalShortString // Capacity 120
	Address         PascalShortString // Capacity 40
	Receiver        PascalShortString // Capacity 40
	AmountFee       int64
	AmountTransfer  int64
	Signature       PascalShortString // Capacity 120
	TransferID      PascalShortString // Capacity 64
}

// ReadFromStream reads a transaction from a stream
func (t *LegacyTransaction) ReadFromStream(r io.Reader) error {
	// Check if the stream is nil
	if r == nil {
		return errors.New("nil reader provided")
	}

	// Field Block
	err := binary.Read(r, binary.LittleEndian, &t.Block)
	if err != nil {
		return err
	}

	// Field OrderID
	t.OrderID = *NewPascalShortString(64)
	err = t.OrderID.ReadFromStream(r)
	if err != nil {
		return err
	}

	// Field OrderLinesCount
	err = binary.Read(r, binary.LittleEndian, &t.OrderLinesCount)
	if err != nil {
		return err
	}

	// Field OrderType
	t.OrderType = *NewPascalShortString(6)
	err = t.OrderType.ReadFromStream(r)
	if err != nil {
		return err
	}

	// Field Timestamp
	err = binary.Read(r, binary.LittleEndian, &t.TimeStamp)
	if err != nil {
		return err
	}

	// Field Reference
	t.Reference = *NewPascalShortString(64)
	err = t.Reference.ReadFromStream(r)
	if err != nil {
		return err
	}

	// Field TransferIndex
	err = binary.Read(r, binary.LittleEndian, &t.TransferIndex)
	if err != nil {
		return err
	}

	// Field Sender
	t.Sender = *NewPascalShortString(120)
	err = t.Sender.ReadFromStream(r)
	if err != nil {
		return err
	}

	// Field Address
	t.Address = *NewPascalShortString(40)
	err = t.Address.ReadFromStream(r)
	if err != nil {
		return err
	}

	// Field Receiver
	t.Receiver = *NewPascalShortString(40)
	err = t.Receiver.ReadFromStream(r)
	if err != nil {
		return err
	}

	// Field AmountFee
	err = binary.Read(r, binary.LittleEndian, &t.AmountFee)
	if err != nil {
		return err
	}

	// Field AmountTransfer
	err = binary.Read(r, binary.LittleEndian, &t.AmountTransfer)
	if err != nil {
		return err
	}

	// Field Signature
	t.Signature = *NewPascalShortString(120)
	err = t.Signature.ReadFromStream(r)
	if err != nil {
		return err
	}

	// Field TransferID
	t.TransferID = *NewPascalShortString(64)
	err = t.TransferID.ReadFromStream(r)
	if err != nil {
		return err
	}

	return nil
}
