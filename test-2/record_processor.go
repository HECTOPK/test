package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func NewRecordProcessor() *RecordProcessor {
	rp := &RecordProcessor{}
	go func() {
		prevCount := int64(0)
		prevEarly := int64(0)
		prevLate := int64(0)
		for range time.Tick(time.Second) {
			count := atomic.LoadInt64(&rp.count)
			early := atomic.LoadInt64(&rp.early)
			late := atomic.LoadInt64(&rp.late)
			fmt.Printf("%d records processed for the last second (%d too early, %d too late)\n",
				count-prevCount, early-prevEarly, late-prevLate)
			prevCount = count
			prevEarly = early
			prevLate = late
		}
	}()
	return rp
}

type RecordProcessor struct {
	sum   int64
	early int64
	late  int64
	count int64
}

func (rp *RecordProcessor) Process(r *Record) {
	t := time.Now()
	d := t.Sub(r.CreatedAt.Add(r.Timeout))
	if d > time.Millisecond*100 {
		atomic.AddInt64(&rp.late, 1)
	} else if d < 0 {
		atomic.AddInt64(&rp.early, 1)
	}
	atomic.AddInt64(&rp.count, 1)
	atomic.AddInt64(&rp.sum, int64(r.Value))
}

func (rp *RecordProcessor) PrintResult() {
	fmt.Println("Sum:", atomic.LoadInt64(&rp.sum))
}
