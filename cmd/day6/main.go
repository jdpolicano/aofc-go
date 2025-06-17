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

func (c Coordinate) move(dir int) Coordinate {
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
		}
	}
	return Coordinate{-1, -1}
}

type Distinct struct {
	dir int
	c   Coordinate
}

type Board struct {
	data       [][]byte
	coordinate Coordinate
	direction  int
	visited    map[Distinct]bool
	done       bool
	isCycle    bool
}

func NewBoard(b [][]byte) *Board {
	startPos := getStartPos(b)
	if isValid(b, startPos) {
		return &Board{b, startPos, Up, make(map[Distinct]bool, 1024), false, false}
	}
	return nil
}

func (b *Board) Run() {
	for !b.done {
		b.move()
	}
}

func (b *Board) move() {
	key := Distinct{b.direction, b.coordinate}
	if b.visited[key] {
		b.done = true
		b.isCycle = true
		return
	}
	b.visited[key] = true
	next := b.coordinate.move(b.direction)
	value := b.get(next)
	switch value {
	case OutBounds:
		{
			b.done = true
			b.coordinate = next
			return
		}
	case Guard, Empty:
		{
			b.data[next.row][next.col] = Visited
			b.coordinate = next
			return
		}
	case Visited:
		{
			b.coordinate = next
			return
		}
	case Obstacle:
		{
			b.turn()
			return
		}
	}
}

func (b *Board) get(c Coordinate) byte {
	if !isValid(b.data, c) {
		return OutBounds
	}
	return b.data[c.row][c.col]
}

func (b *Board) turn() {
	switch b.direction {
	case Up:
		{
			b.direction = Right
		}
	case Right:
		{
			b.direction = Down
		}
	case Down:
		{
			b.direction = Left
		}
	case Left:
		{
			b.direction = Up
		}
	}
}

func getStartPos(b [][]byte) Coordinate {
	for i := range len(b) {
		for j := range len(b[i]) {
			if b[i][j] == Guard {
				return NewCoordinate(i, j)
			}
		}
	}
	return NewCoordinate(-1, -1)
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
	waysToCycle := 0
	for i := range lines {
		for j := range lines[i] {
			if lines[i][j] == Guard {
				continue
			}
			cp := deepCopyBytes(lines)
			cp[i][j] = Obstacle
			board := NewBoard(cp)
			board.Run()
			if board.isCycle {
				waysToCycle++
			}
		}
	}
	fmt.Println(waysToCycle)
	return
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
