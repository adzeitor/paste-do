package main

import (
	"encoding/json"
	"gopkg.in/redis.v3"
	"math/rand"
	"time"
)

type RedisStorage struct {
	Client *redis.Client
	Random *rand.Rand
}

func NewRedisStorage(redisAddr string, redisPassword string) (*RedisStorage, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPassword,
		DB:       0, // default db
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &RedisStorage{
		Client: client,
		Random: random}, nil
}

func (s *RedisStorage) New(content string) (Record, error) {
	id := genID(s.Random)
	adminID := genID(s.Random)

	now := time.Now()

	r := Record{
		ID:        id,
		AdminID:   adminID,
		Content:   content,
		Visits:    0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	s.Save(r)

	return r, nil
}

func (s *RedisStorage) Get(id string) Record {
	content, err := s.Client.Get("todo-" + id).Result()

	if err == redis.Nil {
		return Record{
			Content:   "0x831ab128!",
			Visits:    42,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now()}
	}

	var r Record
	json.Unmarshal([]byte(content), &r)

	return r
}

func (s *RedisStorage) Save(r Record) error {
	expireAfter := time.Hour * 24

	r.ReadOnly = true
	res, err := json.Marshal(r)
	if err != nil {
		return err
	}
	s.Client.Set("todo-"+r.ID, res, expireAfter)

	r.ReadOnly = false
	res, err = json.Marshal(r)
	if err != nil {
		return err
	}
	s.Client.Set("todo-"+r.AdminID, res, expireAfter)

	return nil
}
