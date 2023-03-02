package main

import "time"

type Storage interface {
	Set(key, value string)
	Get(key string) (string, bool)
}

type InMemoryStorage struct {
	data     map[string]string
	expirity map[string]int64
}

func NewStorage() *InMemoryStorage {
	return &InMemoryStorage{
		data:     make(map[string]string),
		expirity: make(map[string]int64),
	}
}

func (s *InMemoryStorage) Set(key, value string, expiration int64) {
	s.data[key] = value
	if expiration != 0 {
		s.expirity[key] = time.Now().UnixMilli() + expiration
	}
}

func (s *InMemoryStorage) Get(key string) (value string, ok bool) {
	value, ok = s.data[key]
	if expiration, eok := s.expirity[key]; eok && expiration < time.Now().UnixMilli() {
		delete(s.expirity, key)
		return "", false
	}
	return
}
