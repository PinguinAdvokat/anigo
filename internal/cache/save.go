package cache

import (
	"encoding/gob"
	"os"
)

func (c *Cache) Save() error {
	f, err := os.Create(c.FilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := gob.NewEncoder(f)
	c.Range(func(key, value interface{}) bool {
		enc.Encode([]interface{}{key, value})
		return true
	})
	return nil
}
