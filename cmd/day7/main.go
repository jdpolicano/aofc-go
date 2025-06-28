package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

type Location [2]int

func (n Location) String() string {
	return fmt.Sprintf("{ row: %d, col: %d }", n.Row(), n.Col())
}

func (n Location) Row() int {
	return n[0]
}

func (n Location) Col() int {
	return n[1]
}

func (n Location) Antinodes(other Location) [2]Location {
	r1, c1, r2, c2 := n.Row(), n.Col(), other.Row(), other.Col()
	rDiff, cDiff := r1-r2, c1-c2
	return [2]Location{
		n.Offset(rDiff, cDiff),
		other.Offset(-rDiff, -cDiff),
	}
}

func (n Location) Offset(rDiff, cDiff int) Location {
	return Location{n.Row() + rDiff, n.Col() + cDiff}
}

func (n Location) AllAntinodes(other Location, inBounds func(Location) bool) []Location {
	antis := make([]Location, 0, 32)
	r1, c1, r2, c2 := n.Row(), n.Col(), other.Row(), other.Col()
	rDiff, cDiff := r1-r2, c1-c2
	// first calulate all of the antis going "up"
	curr := n.Offset(rDiff, cDiff)
	for inBounds(curr) {
		antis = append(antis, curr)
		curr = curr.Offset(rDiff, cDiff)
	}

	// then calulate all of the antis going "down"
	curr = other.Offset(-rDiff, -cDiff)
	for inBounds(curr) {
		antis = append(antis, curr)
		curr = curr.Offset(-rDiff, -cDiff)
	}
	return antis
}

// safely sets a value into a map[K][]V - if the key doesn't exist it is created.
func set[K comparable, V any](m map[K][]V, key K, value V) {
	if list, exists := m[key]; exists {
		m[key] = append(list, value)
	} else {
		list := make([]V, 0, 1024)
		list = append(list, value)
		m[key] = list
	}
}

func makeInBoundsFn[T any](arr [][]T) func(Location) bool {
	maxR := len(arr)
	maxC := len(arr[0])
	return func(n Location) bool {
		return n.Row() >= 0 && n.Row() < maxR && n.Col() >= 0 && n.Col() < maxC
	}
}

func main() {
	b, err := os.ReadFile("./cmd/day7/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	trimmed := bytes.Trim(b, "\n\r\t ")
	lines := bytes.Split(trimmed, []byte("\n"))
	nodes := make(map[byte][]Location)
	antinodes := make(map[Location][]byte)
	inBounds := makeInBoundsFn(lines)
	for r := range lines {
		for c := range lines[r] {
			// if there isn't a node here, move on
			if lines[r][c] == '.' {
				continue
			}
			// now, check the lines between this node and the others of the same type we have passed.
			if prev, exists := nodes[lines[r][c]]; exists {
				// for each previous node...
				for _, loc := range prev {
					// get the two nodes that are colinear
					antis := Location{r, c}.AllAntinodes(loc, inBounds)
					// if they are inbounds, add them to the list of antinodes at {r,c}
					for _, l := range antis {
						set(antinodes, l, lines[r][c])
					}
				}
			}
			set(nodes, lines[r][c], Location{r, c})
		}
	}
	for k := range antinodes {
		fmt.Println(k, string(antinodes[k]))
	}
	fmt.Println("before including the nodes themsevles, ", len(antinodes))
	for c, occurances := range nodes {
		if len(occurances) > 2 {
			for _, l := range occurances {
				set(antinodes, l, c)
			}
		}
	}
	fmt.Println("final count", len(antinodes))
}
