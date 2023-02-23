package cokacoin

import (
	"errors"
	"fmt"
)

var (
	ErrNotInitialized = errors.New("coka coin store is empty.")
)

type Store interface {
	Load(hash []byte) (*Block, error)
	Append(b *Block) error
	LastHash() ([]byte, error)
}

func Iterate(store Store, fn func(b *Block) error) error {
	last, err := store.LastHash()
	if err != nil {
		return fmt.Errorf("cannot get last hash: %v\n", err)
	}

	for {
		b, err := store.Load(last)
		if err != nil {
			return fmt.Errorf("cannot iterate more: %v\n", err)
		}

		if err := fn(b); err != nil {
			return fmt.Errorf("error occures while fn(b): %v\n", err)
		}

		if len(b.PrevHash) == 0 {
			return nil
		}

		last = b.PrevHash
	}

}
