package worker

import (
	"fmt"
	"sync"
)

func workerPool(values chan int, result chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range values {
		result <- (v + 1) * (v + 1)
	}
}

func producer(values chan<- int, n int) {
	for v := range n {
		values <- v
	}
	close(values)
}

func MathPow() {
	n := 10
	w := 5

	var wg sync.WaitGroup
	values := make(chan int, n)
	result := make(chan int, n)


	for range w {
		wg.Add(1)
		go workerPool(values, result, &wg)
	}

	go producer(values, n)
	
	wg.Wait()
	close(result)

	var sum int
	for v := range result {
		sum += v
	}
	
	fmt.Println(sum)
}
