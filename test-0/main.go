package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const maxGoroutines = 5

func main() {
	urls := []string{
		"https://golang.org",
		"https://golang.org",
	}
	p := newProcessor()
	results := make(chan result, maxGoroutines)

	if len(urls) <= maxGoroutines {
		for _, url := range urls {
			go func(url string) {
				results <- p.process(url)
			}(url)
		}
	} else {
		work := make(chan string)
		go func() {
			for _, url := range urls {
				work <- url
			}
			close(work)
		}()
		for i := 0; i < maxGoroutines; i++ {
			go func() {
				for url := range work {
					results <- p.process(url)
				}
			}()
		}
	}
	total := 0
	for range urls {
		res := <-results
		total += res.count
		if res.err != nil {
			fmt.Printf("Error request to %s: %v\n", res.url, res.err)
		} else {
			fmt.Printf("Count for %s: %d\n", res.url, res.count)
		}
	}
	fmt.Printf("Total: %d\n", total)
}

type result struct {
	url   string
	count int
	err   error
}

type processor struct {
	client  *http.Client
	pattern []byte
}

func newProcessor() *processor {
	return &processor{
		client: &http.Client{
			Timeout: time.Second * 5,
		},
		pattern: []byte("Go"),
	}
}

func (p *processor) process(url string) result {
	resp, err := p.client.Get(url)
	if err != nil {
		return result{
			url:   url,
			count: 0,
			err:   err,
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result{
			url:   url,
			count: 0,
			err:   err,
		}
	}
	return result{
		url:   url,
		count: bytes.Count(body, p.pattern),
		err:   nil,
	}
}
