package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	file, err := os.ReadFile("./cmd/day1/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(file), "\n")
	leftNums := make([]float64, 0, len(lines))
	rightNums := make([]float64, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		nums := strings.Split(line, " ")
		n1, e1 := strconv.Atoi(nums[0])
		if e1 != nil {
			log.Fatal("err converting n1", n1, e1)
		}
		n2, e2 := strconv.Atoi(nums[len(nums)-1])
		if e2 != nil {
			log.Fatal("err converting n2", n2, e2)
		}
		leftNums = append(leftNums, float64(n1))
		rightNums = append(rightNums, float64(n2))
	}
	if len(leftNums) != len(rightNums) {
		log.Fatal("left and right number columns are not the same length")
	}
	occurances := getOccurancesMap(rightNums)

	similiarityScore := 0
	for _, num := range leftNums {
		o, exists := occurances[num]
		if exists {
			similiarityScore += o * int(num)
		}
	}
	fmt.Printf("%d\n", int(similiarityScore))
}

func getOccurancesMap(nums []float64) map[float64]int {
	m := make(map[float64]int, len(nums))
	for _, num := range nums {
		if occurances, exists := m[num]; exists {
			m[num] = occurances + 1
		} else {
			m[num] = 1
		}
	}
	return m
}
