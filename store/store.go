package store

import (
	"container/list"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]any
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]any),
	}
}

func (s *Store) Set(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.data[key]
	if !ok {
		return nil, false
	}

	byteVal, ok := val.([]byte)
	return byteVal, ok
}

func (s *Store) LPush(key string, values ...[]byte) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.data[key]
	if !ok {
		l := list.New()
		for _, v := range values {
			l.PushFront(v)
		}
		s.data[key] = l
		return l.Len()
	}

	l, ok := existing.(*list.List)
	if !ok {
		return 0
	}

	for _, v := range values {
		l.PushFront(v)
	}
	return l.Len()
}

func (s *Store) LPop(key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.data[key]
	if !ok {
		return nil, false
	}

	l, ok := existing.(*list.List)
	if !ok {
		return nil, false
	}

	if l.Len() == 0 {
		return nil, false
	}

	element := l.Front()
	l.Remove(element)
	return element.Value.([]byte), true
}
