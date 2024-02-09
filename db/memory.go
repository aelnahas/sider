package db

import (
	"context"
	"sync"
)

type record struct {
	key string
	val any
}

type memory struct {
	sync.Mutex
	data map[string]*record
}

func newMemory() *memory {
	return &memory{
		data: make(map[string]*record),
	}
}

func (m *memory) set(ctx context.Context, record *record) error {
	m.Lock()
	defer m.Unlock()
	m.data[record.key] = record
	return nil
}

func (m *memory) get(ctx context.Context, key string) (any, error) {
	record, found := m.data[key]
	if !found {
		return nil, nil
	}

	return record.val, nil
}

func (m *memory) del(ctx context.Context, key string) error {
	m.Lock()
	delete(m.data, key)
	m.Unlock()
	return nil
}

func (m *memory) exists(ctx context.Context, key string) bool {
	_, found := m.data[key]
	return found
}
