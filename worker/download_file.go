package worker

import (
	"fmt"
	"sync"
	"time"
)

func loadingFile(file chan int,finished chan <- int,wg *sync.WaitGroup) {
	defer wg.Done()
	for v := range file {
		fmt.Println("loading file: "+fmt.Sprint(v) +"....")
		time.Sleep(time.Second * 2)
		finished <- v
	}
}

func loadedFile(finished chan int, wg *sync.WaitGroup){
	defer wg.Done()
	for v := range finished {
		fmt.Println("loaded file: "+fmt.Sprint(v))
	}
}

func producerFile(file chan <- int, n int) {
	for v := range n {
		file <- v +1
	}
	close(file)
}

func DownloadFile() {

	countFile := 6
	worker := 3
	var wg sync.WaitGroup

	file := make(chan int, countFile)
	finished := make(chan int, countFile)

	go producerFile(file,countFile)

	for range worker {
		wg.Add(1)
		go loadingFile(file,finished,&wg)
	}
	wg.Wait()
	close(finished)

	for range worker {
		wg.Add(1)
		go loadedFile(finished,&wg)
	}
	wg.Wait()

}