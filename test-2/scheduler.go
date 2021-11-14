package main

import (
	"sync"
	"time"
)

const batchDuration = time.Millisecond * 10

const minTimeout = time.Second
const maxTimeout = time.Second * 30
const accuracy = time.Millisecond * 100

type Scheduler struct {
	rp       *RecordProcessor
	records  []batch
	zeroTime time.Time
}

func NewScheduler(rp *RecordProcessor) *Scheduler {
	n := int(maxTimeout/batchDuration) + 100
	return &Scheduler{
		rp:      rp,
		records: make([]batch, n),
	}
}

func (s *Scheduler) Run() {
	s.zeroTime = time.Now()
	s.runReader()
}

func (s *Scheduler) write(r Record) {
	d := r.Timeout - time.Since(r.CreatedAt)
	if d <= accuracy {
		s.rp.Process(&r)
		return
	}
	i := int(r.CreatedAt.Add(r.Timeout).Sub(s.zeroTime)/batchDuration+2) % len(s.records)
	s.records[i].Add(r)
}

func (s *Scheduler) runReader() {
	go func() {
		cur := int64(0)
		for {
			t := s.zeroTime.Add(batchDuration * time.Duration(cur))
			if time.Now().Sub(t) > 0 {
				recs, unlock := s.records[cur%int64(len(s.records))].Get()
				for i := range recs {
					s.rp.Process(&recs[i])
				}
				unlock()
				cur++
			} else {
				time.Sleep(batchDuration)
			}
		}
	}()
}

type batch struct {
	mu      sync.Mutex
	records []Record
}

func (r *batch) Get() ([]Record, func()) {
	r.mu.Lock()
	return r.records, func() {
		r.records = r.records[:0]
		r.mu.Unlock()
	}
}

func (r *batch) Add(record Record) {
	r.mu.Lock()
	r.records = append(r.records, record)
	r.mu.Unlock()
}
