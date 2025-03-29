package worker

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func csvWorker(files <-chan string, dir string, result chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range files {
		open, err := os.Open(dir + file)
		if err != nil {
			continue
		}
		scanner := bufio.NewScanner(open)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			nums := strings.Split(line, ",")
			sum := 0
			for _, num := range nums {
				n, err := strconv.Atoi(num)
				if err != nil {
					continue
				}
				sum += n
			}
			result <- sum
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		open.Close()
	}

}

func collectorResult(total chan<- int, result <-chan int) {
	sum := 0
	for num := range result {
		sum += num
	}
	total <- sum
}

func DataParallel() {
	dir := "./assets/dataparallel/"
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	w := 3
	fileCH := make(chan string, len(files))
	result := make(chan int)
	total := make(chan int)

	//producer
	for _, file := range files {
		fileCH <- file.Name()
	}
	close(fileCH)

	for range w {
		wg.Add(1)
		go csvWorker(fileCH, dir, result, &wg)
	}

	go collectorResult(total, result)

	go func() {
		wg.Wait()
		close(result)
	}()

	sum := <-total
	fmt.Println("sum: " + fmt.Sprint(sum))

	
}
