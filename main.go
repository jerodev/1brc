package main

import (
	_ "embed"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"sync"
)

const (
	BUFF_LEN      = 2 * 1024 * 1024
	ROUTINE_COUNT = 12
)

//go:embed measurements-1000000000.out
var expected string

func main() {
	var err error
	var i, n int

	// Start a number of background processes
	results := map[string]*City{}
	resultChnl := make(chan map[string][]int, ROUTINE_COUNT*2)
	chunkSenderChnl := make(chan []byte, ROUTINE_COUNT*2)
	for range ROUTINE_COUNT {
		go func(rx chan []byte, tx chan map[string][]int) {
			for {
				chunk, more := <-rx
				if !more {
					break
				}

				parseBuffer(tx, chunk)
			}
		}(chunkSenderChnl, resultChnl)
	}

	// Start parsing results coming from the resultChnl
	var wg sync.WaitGroup
	go func() {
		for {
			r, more := <-resultChnl
			if !more {
				break
			}

			for c, t := range r {
				if _, ok := results[c]; ok {
					results[c].Merge(t)
				} else {
					results[c] = NewCity(c, t)
				}
			}

			wg.Done()
		}
	}()

	f, _ := os.Open("measurements-1000000000.txt")
	defer f.Close()

	// Read the file in chunks and send data on the channels
	var leftOver []byte
	var prefix []byte
	var chunk []byte
	buff := make([]byte, BUFF_LEN)
	for {
		n, err = f.Read(buff)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}

			break
		}

		// Full buffer, cut of last part
		if n == BUFF_LEN {
			// Search for the last newline to fill the leftover
			for i = BUFF_LEN - 1; i > 0; i-- {
				if buff[i] == '\n' {
					break
				}
			}

			// Copy prefix from last run leftovers
			prefix = make([]byte, len(leftOver))
			copy(prefix, leftOver)

			chunk = buff[:i]

			// Copy leftovers from buffer
			leftOver = make([]byte, len(buff[i+1:]))
			copy(leftOver, buff[i+1:])
		} else {
			chunk = buff[:n]
			prefix = leftOver
		}

		wg.Add(1)
		chunkSenderChnl <- append(prefix, chunk...)
	}
	close(chunkSenderChnl)

	wg.Wait()
	close(resultChnl)

	// Sort the results
	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	sortedResults := make([]*City, len(results))
	for i, k := range keys {
		sortedResults[i] = results[k]
	}

	// Le special output
	cities := make([]string, 0, len(sortedResults))
	for _, r := range sortedResults {
		cities = append(cities, r.ToString())
	}

	r := "{" + strings.Join(cities, ", ") + "}"
	fmt.Println(r)

	if r == expected {
		fmt.Println("Correct!")
	} else {
		fmt.Println("No no no...")
	}
}

// parseBuffer cuts a chunk in lines and parses each line
func parseBuffer(tx chan map[string][]int, chunk []byte) {
	var city string
	var temp int
	var ok bool

	results := map[string][]int{}

	var start, ptr int
	for ptr = range chunk {
		if chunk[ptr] == '\n' {
			city, temp = parseLine(chunk[start:ptr])

			if _, ok = results[city]; ok {
				results[city] = append(results[city], temp)
			} else {
				results[city] = []int{temp}
			}

			ptr++
			start = ptr
		}
	}

	tx <- results
}

// Parse line returns the name of the city and the temperature as an integer
func parseLine(line []byte) (string, int) {
	var number, ptr int
	ten := 1
	for ptr = len(line) - 1; ptr >= 0; ptr-- {
		if line[ptr] == ';' {
			break
		}

		if line[ptr] >= '0' && line[ptr] <= '9' {
			number += int(line[ptr]-'0') * ten
			ten *= 10

			continue
		}

		if line[ptr] == '-' {
			number *= -1
			ptr--
			break
		}
	}

	return string(line[:ptr]), number
}
