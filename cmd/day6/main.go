package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

const (
	Left = iota
	Right
	Up
	Down
	OutBounds
	Obstacle = '#'
	Empty    = '.'
	Visited  = 'X'
	Guard    = '^'
)

type Coordinate struct {
	row int
	col int
}

func NewCoordinate(row, col int) Coordinate {
	return Coordinate{row, col}
}

func (c Coordinate) get(data [][]byte) byte {
	if !isValid(data, c) {
		return OutBounds
	}
	return data[c.row][c.col]
}

func (c Coordinate) up() Coordinate {
	return Coordinate{c.row - 1, c.col}
}

func (c Coordinate) down() Coordinate {
	return Coordinate{c.row + 1, c.col}
}

func (c Coordinate) left() Coordinate {
	return Coordinate{c.row, c.col - 1}
}

func (c Coordinate) right() Coordinate {
	return Coordinate{c.row, c.col + 1}
}

func (c Coordinate) dist(other Coordinate) int {
	diff := 0
	if c.row == other.row {
		diff = c.col - other.col
	} else {
		diff = c.row - other.row
	}
	if diff < 0 {
		return -diff
	}
	return diff
}

func (c Coordinate) move(dir int) Coordinate {
	var none Coordinate
	switch dir {
	case Up:
		{
			return c.up()
		}
	case Down:
		{
			return c.down()
		}
	case Left:
		{
			return c.left()
		}
	case Right:
		{
			return c.right()
		}
	default:
		{
			log.Fatal("unknown direction", dir)
			// unreachable
			return none
		}
	}
}

type JumpTable struct {
	tbl  [][][]Coordinate // three dimensions, ([row][col][destination if you go one of four directions]
	data [][]byte
}

func BuildJumpTable(grid [][]byte) *JumpTable {
	tbl := make([][][]Coordinate, len(grid))
	visited := make([][][]bool, len(grid))
	for i := range grid {
		tbl[i] = make([][]Coordinate, len(grid[i]))
		visited[i] = make([][]bool, len(grid[i]))
		for j := range grid[i] {
			tbl[i][j] = make([]Coordinate, 4)
			visited[i][j] = make([]bool, 4)
		}
	}
	jmp := &JumpTable{tbl: tbl, data: grid}
	for i := range len(tbl) {
		for j := range len(tbl[i]) {
			co := NewCoordinate(i, j)
			jmp.Update(co, Up, visited)
			jmp.Update(co, Right, visited)
			jmp.Update(co, Down, visited)
			jmp.Update(co, Left, visited)
		}
	}
	return jmp.tbl
}

func (jmp *JumpTable) Update(start Coordinate, dir int, visited [][][]bool) Coordinate {
	if visited[start.row][start.col][dir] {
		return jmp.tbl[start.row][start.col][dir]
	}
	next := start.move(dir)
	if !isValid(jmp.data, next) {
		visited[start.row][start.col][dir] = true
		jmp.tbl[start.row][start.col][dir] = next
		return next
	}
	if next.get(jmp.data) == Obstacle {
		visited[start.row][start.col][dir] = true
		jmp.tbl[start.row][start.col][dir] = start
		return start
	}
	visited[start.row][start.col][dir] = true
	jmp.tbl[start.row][start.col][dir] = jmp.Update(next, dir, visited)
	return jmp.tbl[start.row][start.col][dir]
}

func (jmp *JumpTable) get(start Coordinate, dir int) Coordinate {
	return jmp.tbl[start.row][start.col][dir]
}

func turn(direction int) int {
	switch direction {
	case Up:
		{
			return Right
		}
	case Right:
		{
			return Down
		}
	case Down:
		{
			return Left
		}
	case Left:
		{
			return Up
		}
	}
	log.Fatal("direction should be an enum")
	return -1
}

func getStartPos(b [][]byte) Coordinate {
	for i := range len(b) {
		for j := range len(b[i]) {
			if b[i][j] == Guard {
				return NewCoordinate(i, j)
			}
		}
	}
	log.Fatal("getStartPos() No guard found")
	return NewCoordinate(0, 0)
}

func isValid(b [][]byte, c Coordinate) bool {
	n, m := len(b), len(b[0])
	return c.row >= 0 && c.row < n && c.col >= 0 && c.col < m
}

func main() {
	b, err := os.ReadFile("./cmd/day6/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	trimmed := bytes.Trim(b, "\n\r\t ")
	lines := bytes.Split(trimmed, []byte("\n"))
	jmp := BuildJumpTable(lines)
	curr := getStartPos(jmp.data)
	dist := 0
	direction := Up
	quit := false
	for !quit {
		next := jmp.get(curr, direction)
		nextNeighbor := next.move(direction)
		if !isValid(jmp.data, nextNeighbor) {
			break
		}
		dist += curr.dist(next)
		curr = next
		direction = turn(direction)
	}
	fmt.Println(dist)
}

func deepCopyBytes(original [][]byte) [][]byte {
	if original == nil {
		return nil
	}
	cp := make([][]byte, len(original))
	for i, inner := range original {
		cp[i] = make([]byte, len(inner))
		copy(cp[i], inner)
	}
	return cp
}

// lines := [][]byte{
// 	{'.', '.', '.', '#', '.', '.', '#', '.', '.', '.'},
// 	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
// 	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
// 	{'.', '.', '.', '#', '.', '.', '.', '.', '.', '.'},
// 	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
// 	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
// 	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
// 	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
// 	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
// 	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
// }
