package main

import (
	"errors"
	"time"
)

func PasteEdit(storage Storage, id string, content string) error {
	r := storage.Get(id)

	if r.ReadOnly {
		return errors.New("read only link")
	}

	r.Content   = content

	if len(r.Content) > 1024*1024 {
		return errors.New("Paste is too long")
	}

	r.UpdatedAt = time.Now()

	storage.Save(r)


	return nil
}
