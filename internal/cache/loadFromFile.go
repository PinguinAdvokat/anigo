package cache

import (
	"encoding/gob"
	"io"
	"log"
	"os"
)

func (c *Cache) LoadFromFile() error {
	f, err := os.Open(c.FilePath)
	if err != nil {
		log.Printf("Cant Open cache file: %v", err)
		return err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)
	for {
		var k, v interface{}

		err := dec.Decode(&k)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		dec.Decode(&v)
		c.Store(k.(string), v)
	}
	return nil
}
