package main

import (
	"math/rand"
	"time"
)

var letterRunes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

type Record struct {
	ID        int
	CreatedAt time.Time
	Name      string
	Value     int
	Timeout   time.Duration
}

func GenerateRandomRecord() Record {
	randString := func(n int) string {
		b := make([]byte, n)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		return string(b)
	}

	return Record{
		ID:        rand.Intn(100000000),
		CreatedAt: time.Now(),
		Name:      randString(16),
		Value:     rand.Intn(10),
	}
}
