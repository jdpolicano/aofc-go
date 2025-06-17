package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	file, err := os.ReadFile("./cmd/day3/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	re := regexp.MustCompile("do\\(\\)|don't\\(\\)|mul\\((\\d{1,3}),(\\d{1,3})\\)")
	sums := 0
	apply := true
	for _, match := range re.FindAllSubmatch(file, -1) {
		if string(match[0]) == "do()" {
			apply = true
			continue
		}
		if string(match[0]) == "don't()" {
			apply = false
			continue
		}
		if !apply {
			continue
		}
		n1, e1 := strconv.Atoi(string(match[1]))
		if e1 != nil {
			log.Fatalf("regex failed to parse first number correctly from %s", match[0])
		}
		n2, e2 := strconv.Atoi(string(match[2]))
		if e2 != nil {
			log.Fatalf("regex failed to parse second number correctly from %s", match[0])
		}
		prod := n1 * n2
		sums += prod
		fmt.Printf("full string \"%s\"\nresult = %d\n", string(match[0]), prod)
	}
	fmt.Println(sums)
}
