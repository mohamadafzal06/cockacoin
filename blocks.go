package cokacoin

import (
	"bytes"
	"errors"
	"fmt"
	"time"
)

type Block struct {
	Timestamp time.Time
	Data      []byte
	Nonce     int32
	PrevHash  []byte
	Hash      []byte
}
type Blockchain struct {
	Difficulty int
	Mask       []byte
	store      Store
}

func NewBlockchain(difficulty int, store Store) (*Blockchain, error) {
	mask := GenerateMask(difficulty)
	bc := Blockchain{
		Difficulty: difficulty,
		store:      store,
		Mask:       mask,
	}
	_, err := bc.store.LastHash()

	if err == nil {
		return &bc, nil
	}

	if !errors.Is(err, ErrNotInitialized) {
		return nil, err // don't need for wrapping, it. we did this previously

	}

	gb := NewBlock("Genesis Block", bc.Mask, []byte{})
	if err := bc.store.Append(gb); err != nil {
		return nil, fmt.Errorf("error while add Genesis Block to the Chain: %v\n", err)
	}

	return &bc, err
}

func (bc *Blockchain) Add(date string) (*Block, error) {
	lHash, err := bc.store.LastHash()
	if err != nil {
		return nil, fmt.Errorf("cannot load hash of last block: %v\n", err)
	}
	b := NewBlock("data", bc.Mask, lHash)
	if err := bc.store.Append(b); err != nil {
		return nil, fmt.Errorf("cannot add this new Block to the Chain: %v\n", err)
	}

	return b, nil
}

func (bc *Blockchain) Print() error {
	fmt.Printf("Difficulty: %d\nStore: %T\n", bc.Difficulty, bc.store)

	return Iterate(bc.store, func(b *Block) error {
		fmt.Print(b)
		return nil
	})

	//	return nil
}

func NewBlock(data string, mask, prevHash []byte) *Block {
	b := Block{
		Timestamp: time.Now(),
		Data:      []byte(data),
		PrevHash:  prevHash,
	}
	b.Hash, b.Nonce = DifficultHash(mask, b.Timestamp.UnixNano(), b.Data, b.PrevHash)
	return &b
}

func (b *Block) String() string {
	return fmt.Sprintf("Time: %s\nData: %s\nHash: %x\nPrevHash: %x\nNonce: %d\n------------\n",
		b.Timestamp, b.Data, b.Hash, b.PrevHash, b.Nonce)

}

func (b *Block) Validate(mask []byte) error {
	h := EasyHash(b.Timestamp.UnixNano(), b.Data, b.PrevHash, b.Nonce)
	if !bytes.Equal(h, b.Hash) {
		return fmt.Errorf("the hash is invalid it should be %x is %x\n", h, b.Hash) // %x is for print in hex
	}

	if !GoodEnough(mask, h) {
		return fmt.Errorf("hash is not good enough with mask %x\n", mask)
	}

	return nil
}

func (bc *Blockchain) Validate() error {
	return Iterate(bc.store, func(b *Block) error {

		if err := b.Validate(bc.Mask); err != nil {
			return fmt.Errorf("block is not valid: %v\n", err)
		}
		return nil

	})
}
