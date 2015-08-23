package main


import (
	"math/rand"
	"time"
)

type NewRecord struct {
	Content string
	Result  chan string
}

type GetRecord struct {
	ID     string
	Result chan *Record
}

type EditRecord struct {
	ID     string
	Content string
	Result chan *Record
}


type Record struct {
	Content string
	Visits  uint64
	CreatedAt time.Time
	UpdatedAt time.Time
	// add lat and lng
}

type Storage struct {
	NewChan chan NewRecord
	GetChan chan GetRecord
	EditChan chan EditRecord
}

func genID(rand *rand.Rand) string {
	const permitted = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	const size = 16

	id := new([size]uint8)
	for i := 0; i < size; i++ {
		r := rand.Intn(len(permitted))
		id[i] = permitted[r]
	}
	return string(id[:size])
}

func (s *Storage) Get (id string) *Record {
	ch := make(chan *Record)

	s.GetChan <- GetRecord{
		ID:id,
		Result:ch}
	r := <- ch
	return r
}


func (s *Storage) New (content string)string {
	ch := make(chan string)

	s.NewChan <- NewRecord {
		Content: content,
		Result: ch}
	r := <- ch
	return r
}

func (s *Storage) Edit (id string, content string) *Record {
	ch := make(chan *Record)

	s.EditChan <- EditRecord {
		ID: id,
		Content: content,
		Result: ch}
	r := <- ch
	return r
}
