package main

import (
	"gopkg.in/redis.v3"
	"log"
	"math/rand"
	"time"
)

func NewRedisStorage(redisAddr string, redisPassword string) Storage {

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0, // default db
	})

	pong, err := client.Ping().Result()
	log.Println("redis ping: ", pong, err)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	newRecords := make(chan NewRecord)
	getRecords := make(chan GetRecord)
	editRecords := make(chan EditRecord)

	maxLength := 2048
	expireAfter := time.Hour * 24

	go func() {
		for {
			select {
			// new
			case c := <-newRecords:
				id := genID(r)
				if len(c.Content) > maxLength {
					c.Content = c.Content[:maxLength]
				}
				client.Set("todo-"+id, c.Content, 0)
				c.Result <- id
			// get
			case c := <-getRecords:
				content, err := client.Get("todo-" + c.ID).Result()

				if err == redis.Nil {
					log.Println("redis error:", err)
					c.Result <- &Record{
						Content:   "0x831ab128!",
						Visits:    42,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now()}
					continue
				}

				r := &Record{
					Content:   content,
					Visits:    42,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now()}

				c.Result <- r
			// edit
			case c := <-editRecords:
				content, err := client.Get("todo-" + c.ID).Result()

				if err == redis.Nil {
					c.Result <- &Record{
						Content:   "0x831ab128!",
						Visits:    42,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now()}
					continue
				}
				client.Set("todo-"+c.ID, c.Content, expireAfter)

				r := &Record{
					Content:   content,
					Visits:    42,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now()}
				c.Result <- r
			}
		}
		// get

	}()

	return Storage{
		NewChan:  newRecords,
		GetChan:  getRecords,
		EditChan: editRecords}
}
