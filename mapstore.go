package cokacoin

import (
	"bytes"
	"fmt"
)

type mapStore struct {
	data map[string]*Block
	last []byte
}

func (ms *mapStore) Load(hash []byte) (*Block, error) {
	x := fmt.Sprintf("%x", hash)
	if b, ok := ms.data[x]; !ok {
		return b, nil
	}

	return nil, fmt.Errorf("block is not in this store.")
}

func (ms *mapStore) Append(b *Block) error {
	if !bytes.Equal(ms.last, b.PrevHash) {
		return fmt.Errorf("store is out of sync.")
	}

	x := fmt.Sprintf("%x", b.Hash)
	if _, ok := ms.data[x]; ok {
		return fmt.Errorf("duplicat hash")
	}

	ms.data[x] = b
	ms.last = b.Hash

	return nil
}

func (ms *mapStore) LastHash() ([]byte, error) {
	if len(ms.last) == 0 {
		return nil, ErrNotInitialized
	}

	return ms.last, nil
}

func NewMapStore() Store {
	return &mapStore{
		data: make(map[string]*Block),
	}

}
