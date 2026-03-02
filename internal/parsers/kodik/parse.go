package kodik

import (
	"context"
	"fmt"
	"log"
	"net/url"
)

func (k *Kodik) Parse(ctx context.Context, kodikURL string) (map[string]string, error) {
	var payload url.Values
	v, ok := k.cache.Load(kodikURL)
	if !ok {
		log.Printf("value not cached\n")
		var err error
		payload, err = k.getPayload(ctx, kodikURL)
		if err != nil {
			return nil, err
		}
		log.Printf("caching value: %v\n", payload)
		k.cache.Store(kodikURL, payload)
	} else {
		payload = v.(url.Values)
	}

	fmt.Println("fetched payload:")
	fmt.Println(payload)

	sources, err := k.getSources(ctx, payload)
	if err != nil {
		return nil, err
	}
	return sources, nil
}
