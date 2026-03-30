package kodik

import (
	"fmt"
)

func (k *Kodik) Parse(kodikURL string) (map[string]string, error) {
	payload, err := k.getPayload(kodikURL)
	if err != nil {
		return nil, err
	}

	fmt.Println("fetched payload:")
	fmt.Println(payload)

	sources, err := k.getSources(payload)
	if err != nil {
		return nil, err
	}
	return sources, nil
}
