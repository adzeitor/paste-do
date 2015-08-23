package main

import (
	"math/rand"
	"time"
	"fmt"
)



func NewMemoryStorage(prefix string) Storage {
	// FIXME: ...
	r := rand.New(rand.NewSource(42))

	m := make(map[string](*Record))

	newRecords := make(chan NewRecord)
	getRecords := make(chan GetRecord)
	editRecords := make(chan EditRecord)

	go func() {
		for {
			select {
			// new
			case c :=  <-newRecords:
				id := genID(r)
				now := time.Now()
				m[id] = &Record{
					Content: c.Content,
					Visits:  0,
					CreatedAt: now,
					UpdatedAt: now}
				c.Result <- id
			// get
			case c := <-getRecords:
				val, ok := m[c.ID]
				if ok == false {
					c.Result <- &Record{
						Content: "0x831ab128!",
						Visits:  42,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now()}
					continue
				}
				val.Visits += 1
				c.Result <- val
			// edit
			case c := <-editRecords:
				fmt.Println(c)
				val, ok := m[c.ID]
				if ok == false {
					c.Result <- &Record{
						Content: "0x831ab128!",
						Visits:  42,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now()}
					continue
				}
				val.Content = c.Content
				c.Result <- val
			}
		}
		// get

	}()

	return Storage{
		NewChan:newRecords,
		GetChan:getRecords,
		EditChan:editRecords}
}
