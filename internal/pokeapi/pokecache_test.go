package pokeapi_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/faymn0/pokedex/internal/pokeapi"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key   string
		value []byte
	}{
		{
			key:   "https://example.com",
			value: []byte("test data"),
		},
		{
			key:   "https://example.com/api",
			value: []byte("more test data"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %d", i), func(t *testing.T) {
			cache := pokeapi.NewCache(interval)
			cache.Add(c.key, c.value)
			value, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key %s", c.key)
				return
			}
			if string(value) != string(c.value) {
				t.Errorf("expected to find the same value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const interval = 5 * time.Millisecond
	const waitTime = interval * 2

	const key = "https://example.com"

	cache := pokeapi.NewCache(interval)
	cache.Add(key, []byte("this data doesn't really matter"))

	_, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get(key)
	if ok {
		t.Errorf("expected key to be reaped!")
		return
	}
}
