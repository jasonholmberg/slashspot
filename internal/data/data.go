package data

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
)

var store map[string]Spot

// Open - open the spot store
func Open() {
	store = make(map[string]Spot)
	Load()
}

// IsOpen - data store is open
func IsOpen() bool {
	return store != nil
}

// Save - saves the spot stroe
func Save() error {
	// Save to disk
	os.MkdirAll(os.Getenv("SPOT_DATA_DIR"), os.ModePerm)
	f, err := os.Create(FilePath())
	defer f.Close()
	if err != nil {
		log.Printf("Error saving spot store %v", err)
	}
	if store == nil {
		store = make(map[string]Spot)
		return nil
	}
	r, err := marshal(store)
	if err != nil {
		errMsg := "Error marshalling spot-store"
		return errors.New(errMsg)
	}
	_, err = io.Copy(f, r)
	return err
}

// Load - load the data file
func Load() (map[string]Spot, error) {
	// Load from disk
	f, err := os.Open(FilePath())
	defer f.Close()
	if err != nil {
		log.Print("No data file to load, creating one")
		Save()
	}
	unmarshal(f, &store)
	return store, nil
}

// Delete - deletes a spot from the store by id. Save should be called after this operation.
func Delete(id string) {
	delete(store, id)
}

// Insert - adds a spot to the store. Save should be called after this operation.
func Insert(s Spot) {
	store[s.Key()] = s
}

// FilePath - path to data file
func FilePath() string {
	return filepath.Join(os.Getenv("SPOT_DATA_DIR"), os.Getenv("SPOT_DATA_FILE"))
}

func marshal(spots map[string]Spot) (io.Reader, error) {
	b, err := json.MarshalIndent(spots, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

func unmarshal(r io.Reader, data *map[string]Spot) error {
	return json.NewDecoder(r).Decode(data)
}
