package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	processor := NewRecordProcessor()
	s := NewScheduler(processor)

	gens := make([]*RecordGenerator, 10)
	for i := range gens {
		gens[i] = NewRecordGenerator(s, time.Second/10000)
	}

	fmt.Println("run scheduler")
	s.Run()

	fmt.Println("run generators")
	for _, g := range gens {
		g.Run()
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	processor.PrintResult()
}
