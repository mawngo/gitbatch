package util

import "sync"

func SplitParallel[T any](maxParallel int, data []T, handler func(T)) {
	var ch = make(chan T, maxParallel+1)
	var wg sync.WaitGroup

	// This starts maxParallel number of goroutines that wait for something to do
	wg.Add(maxParallel)
	for i := 0; i < maxParallel; i++ {
		go func() {
			for {
				a, ok := <-ch
				if !ok { // if there is nothing to do and the channel has been closed then end the goroutine
					wg.Done()
					return
				}
				handler(a)
			}
		}()
	}

	// Now the jobs can be added to the channel, which is used as a queue
	for _, item := range data {
		ch <- item
	}

	close(ch) // This tells the goroutines there's nothing else to do
	wg.Wait() // Wait for the threads to finish
}
