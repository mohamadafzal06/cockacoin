package cokacoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type folderConfig struct {
	LastHash []byte
}

type folderStore struct {
	root       string
	config     *folderConfig
	configPath string
}

func (fs *folderStore) Load(hash []byte) (*Block, error) {
	path := filepath.Join(fs.root, fmt.Sprintf("%x.json", hash))
	var b Block
	if err := readJSON(path, &b); err != nil {
		return nil, fmt.Errorf("cannot load from file: %v\n", err)
	}

	return &b, nil
}

func (fs *folderStore) Append(b *Block) error {
	if !bytes.Equal(fs.config.LastHash, b.PrevHash) {
		return fmt.Errorf("store is out of sync")
	}

	path := filepath.Join(fs.root, fmt.Sprintf("%x.json", b.Hash))
	if err := writJSON(path, b); err != nil {
		return fmt.Errorf("cannot write to file: %v\n", err)
	}

	fs.config.LastHash = b.Hash
	if err := writJSON(fs.configPath, fs.config); err != nil {
		return fmt.Errorf("cannot update the block: %v\n", err)
	}

	return nil
}

func (fs *folderStore) LastHash() ([]byte, error) {
	if len(fs.config.LastHash) == 0 {
		return nil, ErrNotInitialized
	}

	return fs.config.LastHash, nil
}
func readJSON(path string, v interface{}) error {
	fl, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open file: %v\n", err)

	}
	defer fl.Close()

	dec := json.NewDecoder(fl)

	if err := dec.Decode(v); err != nil {
		return fmt.Errorf("cannot decode: %v\n", err)
	}

	return nil
}

func writJSON(path string, v interface{}) error {
	fl, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("cannot open file: %v\n", err)

	}
	defer fl.Close()

	enc := json.NewEncoder(fl)
	// just for more readablity
	enc.SetIndent("", "  ")

	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("cannot encode: %v\n", err)
	}

	return nil
}

func NewFolderStore(root string) (Store, error) {
	fs := &folderStore{
		root:       root,
		config:     &folderConfig{},
		configPath: filepath.Join(root, "config.json"),
	}

	if err := readJSON(fs.configPath, fs.config); err != nil {
		log.Print("Empty Store")
		fs.config.LastHash = nil
	}

	return fs, nil

}
