package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	BUFF_LEN      = 24_000
	LEFTOVER_LEN  = 128
	ROUTINE_COUNT = 4
)

func main() {
	var err error
	var i int

	// Start a number of background processes
	results := map[string]*City{}
	resultChnl := make(chan map[string][]int, ROUTINE_COUNT*2)
	chunkSenderChnl := make(chan []byte, ROUTINE_COUNT*2)
	for range ROUTINE_COUNT {
		go (func(rx chan []byte, tx chan map[string][]int) {
			for {
				chunk, more := <-rx
				if !more {
					break
				}

				parseBuffer(tx, chunk)
			}
		})(chunkSenderChnl, resultChnl)
	}

	// Start parsing results comming from the resultChnl
	var wg sync.WaitGroup
	go (func() {
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
	})()

	f, _ := os.Open("measurements-1000000000.txt")
	defer f.Close()

	// Read the file in chunks and send data on the channels
	var leftOver []byte
	var prefix []byte
	var chunk []byte
	buff := make([]byte, BUFF_LEN)
	for {
		_, err = f.Read(buff)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}

			break
		}

		// Full buffer, cut of last part
		if len(buff) == BUFF_LEN {
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
			chunk = buff
			prefix = leftOver
		}

		wg.Add(1)
		chunkSenderChnl <- append(prefix, chunk...)
	}
	close(chunkSenderChnl)

	wg.Wait()
	close(resultChnl)

	fmt.Println(results)
}

// parseBuffer cuts a chunk in lines and parses each line
func parseBuffer(tx chan map[string][]int, chunk []byte) {
	var city string
	var temp int
	var ok bool

	line := make([]byte, 0, LEFTOVER_LEN)
	results := map[string][]int{}

	for _, c := range chunk {
		if c == '\n' {
			if len(line) > 0 {
				if !bytes.ContainsRune(line, ';') {
					panic(string(chunk))
				}

				city, temp = parseLine(line)

				if _, ok = results[city]; ok {
					results[city] = append(results[city], temp)
				} else {
					results[city] = []int{temp}
				}

				line = make([]byte, 0, LEFTOVER_LEN)
				continue
			}
		}

		line = append(line, c)
	}

	tx <- results
}

// Parse line returns the name of the city and the temperature as an integer
func parseLine(line []byte) (string, int) {
	parts := bytes.SplitN(line, []byte{';'}, 2)
	if len(parts) != 2 {
		panic(string(line))
	}

	// Super simple number parser
	var number int
	var gotdot bool
	for _, c := range parts[1] {
		if c == '-' {
			continue
		}

		if c == '.' {
			gotdot = true
			continue
		}

		number = number*10 + int(c-'0')
	}

	if parts[1][0] == '-' {
		number *= -1
	}

	if !gotdot {
		number *= 10
	}

	return string(parts[0]), number
}
