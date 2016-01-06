package main

import (
	"math/rand"
	"time"
	"sync"
)

type MemoryStorage struct {
	Items  map[string]Record
	Random *rand.Rand
	sync.Mutex
}

func NewMemoryStorage() *MemoryStorage {

	m := make(map[string]Record)

	return &MemoryStorage{
		Random: rand.New(rand.NewSource(42)),
		Items:  m,
	}
}

func (s *MemoryStorage) New(content string) (Record,error) {
	id := genID(s.Random)
	// FIXME: use independent random generators
	adminID := genID(s.Random)
	now := time.Now()

	r := Record{
		ID: id,
		AdminID: adminID,
		Content: content,
		CreatedAt: now,
		UpdatedAt:now,
		ReadOnly: true,
	}

	s.Items[id] = r

	r.ReadOnly = false
	s.Items[adminID] = r

	return r,nil
}


func (s *MemoryStorage) Get(id string) Record {
	return s.Items[id]
}

func (s *MemoryStorage) Save(r Record) error {
	s.Lock()

	// read only link
	r.ReadOnly = true
	s.Items[r.ID] = r
	// admin link
	r.ReadOnly = false
	s.Items[r.AdminID] = r

	s.Unlock()
	return nil
}
