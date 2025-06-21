package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
)

type Equation struct {
	sum   int64
	parts []int64
}

func NewEquation() *Equation {
	return &Equation{
		sum:   0,
		parts: make([]int64, 0),
	}
}

// parses a line of input into an Equation object
// the line is expected to be in the format "<sum>: <part1> <part2> ... <partN>"
func ParseEquation(line string) *Equation {
	parts := bytes.Split([]byte(line), []byte(": "))
	if len(parts) != 2 {
		log.Fatalf("Invalid equation format: %s", line)
	}

	eq := NewEquation()
	sum, err := strconv.ParseInt(string(parts[0]), 10, 64)
	if err != nil {
		log.Fatalf("Invalid sum in equation: %s", line)
	}
	eq.sum = sum
	partStrings := bytes.Split(parts[1], []byte(" "))
	for _, partStr := range partStrings {
		part, err := strconv.ParseInt(string(partStr), 10, 64)
		if err != nil {
			log.Fatalf("Invalid part in equation: %s", line)
		}
		eq.PushPart(part)
	}

	return eq
}

// will try to concatenate the top two parts in the list together and continue computation from there
func tryConcat(acc, target int64, parts []int64) bool {
	if len(parts) < 1 {
		return false
	}
	newPart := fmt.Sprintf("%d%d", acc, parts[0])
	newPartInt, err := strconv.ParseInt(newPart, 10, 64)
	if err != nil {
		log.Fatalf("Invalid concatenation result: %s", newPart)
	}
	return backtrack(newPartInt, target, parts[1:])
}

func backtrack(acc, target int64, parts []int64) bool {
	if len(parts) == 0 {
		return acc == target
	}

	if acc > target {
		return false
	}

	return backtrack(acc+parts[0], target, parts[1:]) ||
		backtrack(acc*parts[0], target, parts[1:]) ||
		tryConcat(acc, target, parts)
}

func (e *Equation) IsSolvable() bool {
	return backtrack(e.parts[0], e.sum, e.parts[1:])
}

func (e *Equation) PushPart(part int64) {
	e.parts = append(e.parts, part)
}

func (e *Equation) String() string {
	return fmt.Sprintf("Sum: %d, Parts: %v", e.sum, e.parts)
}

func main() {
	b, err := os.ReadFile("./cmd/day7/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	trimmed := bytes.Trim(b, "\n\r\t ")
	lines := bytes.Split(trimmed, []byte("\n"))
	sum := int64(0)
	wg := &sync.WaitGroup{}
	for _, line := range lines {
		wg.Add(1)
		go runLine(line, &sum, wg)
	}
	wg.Wait()
	fmt.Println("Total sum of solvable equations:", atomic.LoadInt64(&sum))
}

func runLine(line []byte, sum *int64, wg *sync.WaitGroup) {
	equation := ParseEquation(string(line))
	if equation.IsSolvable() {
		atomic.AddInt64(sum, equation.sum)
	}
	wg.Done()
}
