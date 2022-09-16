package persist

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
)

var lock sync.Mutex

var marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		log.Printf("%v", err)
		return nil, err
	}
	return bytes.NewReader(b), nil
}

var unmarshal = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// Save saves a representation of v to the file at path.
func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(path)
	if err != nil {
		log.Printf("%v", err)
	}
	defer f.Close()
	r, err := marshal(v)
	if err != nil {
		log.Printf("%v", err)
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

// Load loads the file at path into v.
func Load(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return unmarshal(f, v)
}
