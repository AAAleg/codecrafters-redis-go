package main

type Storage interface {
	Set(key, value string)
	Get(key string) (string, bool)
}

type InMemoryStorage struct {
	data map[string]string
}

func NewStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

func (s *InMemoryStorage) Set(key, value string) {
	s.data[key] = value
}

func (s *InMemoryStorage) Get(key string) (value string, ok bool) {
	value, ok = s.data[key]
	return
}
