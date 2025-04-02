package urlshortener

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

type Storage interface {
	Store(string) string
	Get(string) (string, bool)
}

type inMemoryStorage struct {
	mp        map[string]string
	inverseMP map[string]string
	mutex     *sync.Mutex
}

func (s inMemoryStorage) Store(url string) string {
	s.mutex.Lock()

	key, exists := s.inverseMP[url]
	if !exists {
		key = generateKey()
		s.mp[key] = url
		s.inverseMP[url] = key
	}

	s.mutex.Unlock()

	return key
}

func (s inMemoryStorage) Get(key string) (string, bool) {
	url, found := s.mp[key]
	return url, found
}

func GetRAMStorage() Storage {
	return inMemoryStorage{
		mp:        make(map[string]string),
		inverseMP: make(map[string]string),
		mutex:     &sync.Mutex{},
	}
}

func generateKey() string {
	return fmt.Sprint(rand.Uint64() % 1_000_000)
}
