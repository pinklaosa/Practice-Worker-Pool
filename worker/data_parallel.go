package worker

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
)

func csvWorker(files <-chan string, result chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for file := range files {
		open, err := os.Open(file)
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

	//producer
	for _, file := range files {
		fileCH <- file.Type().String()
	}
	close(fileCH)

	for range w {
		wg.Add(1)
		go csvWorker(fileCH, result, &wg)
	}
	

}
