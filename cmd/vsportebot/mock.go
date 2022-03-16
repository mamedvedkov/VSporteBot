package main

import (
	"context"
	"sync"
)

type registerStore struct {
	data map[int]string
	mu   sync.RWMutex
}

func newRegisterStore() *registerStore {
	return &registerStore{
		data: make(map[int]string),
		mu:   sync.RWMutex{},
	}
}

func (r *registerStore) IsRegistred(ctx context.Context, id int) (role string, ok bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	role, ok = r.data[id]
	return
}

func (r *registerStore) Register(ctx context.Context, id int, role string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[id] = role
}
