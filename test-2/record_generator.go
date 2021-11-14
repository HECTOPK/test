package main

import (
	"math/rand"
	"time"
)

type RecordGenerator struct {
	scheduler *Scheduler
	tick      time.Duration
}

func NewRecordGenerator(s *Scheduler, tick time.Duration) *RecordGenerator {
	return &RecordGenerator{
		scheduler: s,
		tick:      tick,
	}
}

func (rp *RecordGenerator) Run() {
	go func() {
		zeroTime := time.Now()
		i := 0
		for {
			t := zeroTime.Add(rp.tick * time.Duration(i))
			if time.Now().Sub(t) > 0 {
				record := GenerateRandomRecord()
				record.Timeout = time.Duration(rand.Intn(int(maxTimeout-minTimeout))) + minTimeout
				rp.scheduler.write(record)
				i++
			} else {
				time.Sleep(rp.tick)
			}
		}
	}()
}
