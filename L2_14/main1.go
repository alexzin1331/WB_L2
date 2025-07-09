package main

import (
	"fmt"
	"sync"
	"time"
)

func orCycle(channels ...<-chan interface{}) <-chan interface{} {
	if len(channels) == 0 {
		return nil
	}

	orDone := make(chan interface{})
	var once sync.Once //lock + unlock

	for _, ch := range channels {
		go func(c <-chan interface{}) {
			select {
			case <-c:
				once.Do(func() {
					close(orDone)
				})
			case <-orDone:
				//канал уже закрыла какая-то горутина
			}
		}(ch)
	}

	return orDone
}

func sig1(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {
	start := time.Now()
	<-orCycle(
		sig1(2*time.Hour),
		sig1(5*time.Minute),
		sig1(1*time.Second),
		sig1(1*time.Hour),
		sig1(1*time.Minute),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}
