package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jdpolicano/aof-go/internal"
)

func GetDiagnols(data []string, row, col int) []string {
	dirs := [][][]int{
		{{-1, 1}, {0, 0}, {1, -1}},
		{{1, 1}, {0, 0}, {-1, -1}},
	}

	isValid := func(data []string, r, c int) bool {
		return internal.InRange(0, len(data)-1, r) && internal.InRange(0, len(data[r])-1, c)
	}

	if !isValid(data, row, col) {
		return []string{}
	}

	res := make([]string, 0, 2) // at most 2
	for _, diag := range dirs {
		var s string
		for _, d := range diag {
			r, c := row+d[1], col+d[0]
			if !isValid(data, r, c) {
				break
			}
			s += string(data[r][c])
		}
		res = append(res, s)
	}

	return res
}

func GetAllDirections(data []string, row, col, n int) []string {
	dirs := [][]int{
		{0, 1},   // right 0 down 1
		{1, 1},   // right 1 down 1
		{1, 0},   // right 1 up 0
		{1, -1},  // right 1 up 1
		{0, -1},  // right 0 up 1
		{-1, -1}, // left 1 up one
		{-1, 0},  // left 1 up 0
		{-1, 1},  // left 1 down one
	}

	isValid := func(data []string, r, c int) bool {
		return internal.InRange(0, len(data)-1, r) && internal.InRange(0, len(data[r])-1, c)
	}

	if !isValid(data, row, col) {
		return []string{}
	}

	res := make([]string, 0, 8) // at most 8 directions to go
	for _, offset := range dirs {
		var s string
		r, c := row, col
		roff, coff := offset[1], offset[0]
		for range n {
			if !isValid(data, r, c) {
				break
			}
			s += string(data[r][c])
			r, c = r+roff, c+coff
		}
		res = append(res, s)
	}

	return res
}

func main() {
	bytes, err := os.ReadFile("./cmd/day4/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(bytes), "\n")
	isEither := func(s string, compa string, compb string) bool { return s == compa || s == compb }
	count := 0
	for row := range lines {
		for col := range lines[row] {
			if lines[row][col] == 'A' {
				diags := GetDiagnols(lines, row, col)
				if isEither(diags[0], "MAS", "SAM") && isEither(diags[1], "MAS", "SAM") {
					count++
				}
			}
		}
	}
	fmt.Println(count)
}
