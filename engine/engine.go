package engine

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

func Hello() {
	fmt.Println("Hello from Engine!")
}

func Run(r io.Reader, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	go func(s *bufio.Scanner) {
		for s.Scan() {
			line := s.Text()
			if line == "exit" || line == "quit" {
				break
			}
			fmt.Fprintf(w, "You entered: %s\n", line)
		}
		err := s.Err()
		if err != nil {
			fmt.Fprintf(w, "Readrg error: %v\n", err.Error())
		}
	}(scanner)
	round := 1
	for {
		fmt.Fprintf(w, "Begin round %d\n", round)
		round++
		time.Sleep(time.Second * 6)
	}
}
