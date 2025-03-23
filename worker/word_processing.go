package worker

import (
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
)

func logWord(m *sync.Map) {
	fmt.Println("Word count results: ")
	m.Range(func(key, value interface{}) bool {
		fmt.Println(key.(string) + ": " + fmt.Sprint(value.(int)))
		return true
	})
}

func countWord(words <-chan string, m *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	for word := range words {
		if v, ok := m.Load(word); ok {
			m.Store(word, v.(int)+1)
		}
	}
}

func splitWord(producer chan string, words chan<- string, m *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	var parts []string
	for line := range producer {
		parts = append(parts, strings.Fields(line)...)
	}

	for _, v := range parts {
		words <- v
	}

	slices.Sort(parts)
	unqiue := slices.Compact(parts)
	for _, w := range unqiue {
		m.Store(w, 1)
	}

}

func produceLine(lines []string, producer chan<- string) {
	for _, v := range lines {
		producer <- v
	}
	close(producer)
}

func WordProcess() {
	var wg sync.WaitGroup
	var m sync.Map
	file, _ := os.ReadFile("./assets/content.txt")
	lines := strings.Split(string(file), "\n")

	for i, v := range lines {
		fmt.Println("Worker " + fmt.Sprint(i+1) + ": " + v)
	}

	producer := make(chan string, len(lines))
	words := make(chan string)

	worker := 5
	for range worker {
		wg.Add(1)
		go countWord(words, &m, &wg)
	}
	go func() {
		wg.Wait()
		close(words)
	}()

	go produceLine(lines, producer)

	for range len(lines) {
		wg.Add(1)
		go splitWord(producer, words, &m, &wg)
	}

	logWord(&m)
}
