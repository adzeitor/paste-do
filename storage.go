package main

import (
	"math/rand"
	"time"
)

type Record struct {
	ID      string
	AdminID string
	// FIXME: Reader
	Content   string
	Visits    uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	ReadOnly  bool
	// TODO: think about  lat and lng
}

type Storage interface {
	New(string) (Record, error)
	Get(string) Record
	Save(Record) error
}

func genID(rand *rand.Rand) string {
	const permitted = "abcdefghijklmnopqrstuvwxyz0123456789"
	const size = 16

	id := new([size]uint8)
	for i := 0; i < size; i++ {
		r := rand.Intn(len(permitted))
		id[i] = permitted[r]
	}
	return string(id[:size])
}
