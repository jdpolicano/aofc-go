package main

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"

	collections "github.com/jdpolicano/aof-go/internal"
)

func main() {
	file, err := os.ReadFile("./cmd/day2/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(file), "\n")
	lines = collections.FilterSlice(lines, func(l string) bool { return len(l) > 0 })
	asNums := collections.MapSlice(lines, func(line string) []int {
		trimmed := strings.Trim(line, " ")
		items := strings.Split(trimmed, " ")
		filtered := collections.FilterSlice(items, func(l string) bool { return len(l) > 0 })
		return collections.MapSlice(filtered, func(s string) int {
			n, err := strconv.Atoi(s)
			if err != nil {
				log.Fatal(err)
			}
			return n
		})
	})
	safeReports := 0
	for i := range asNums {
		if len(asNums[i]) > 1 {
			allRows := getAllPossibleRows(asNums[i])
			if slices.ContainsFunc(allRows, func(r []int) bool { return testRow(r) }) {
				safeReports++
			}
		}
	}
	fmt.Println(safeReports)
}

func isStrictIncreasing(row []int) bool {
	return row[0] > row[1]
}

func isStrictDecreasing(row []int) bool {
	return row[0] < row[1]
}

func getAllPossibleRows(row []int) [][]int {
	rows := make([][]int, 0, len(row)*len(row)+1)
	for i := range row {
		// build up an array excluding the current element
		r := make([]int, 0, len(row))
		for j := range row {
			if j == i {
				continue
			}
			r = append(r, row[j])
		}
		rows = append(rows, r)
	}
	return append(rows, row)
}

func testRow(row []int) bool {
	if isStrictIncreasing(row) {
		return isSafe(row, func(a, b int) int { return a - b })
	}

	if isStrictDecreasing(row) {
		return isSafe(row, func(a, b int) int { return b - a })
	}

	return false
}

func isSafe(n []int, f func(i, j int) int) bool {
	left, right := 0, 1
	for right < len(n) {
		diff := f(n[left], n[right])
		if !inRange(1, 3, diff) {
			return false
		}
		// we're safe
		left = right
		right++
	}
	return true
}

func inRange(i, j, k int) bool {
	return i <= k && k <= j
}
