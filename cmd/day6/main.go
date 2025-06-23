package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
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

type Direction int

func (d Direction) Turn() Direction {
	switch d {
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
	panic("direction should be an enum")
}

func (d Direction) Opposite() Direction {
	switch d {
	case Up:
		{
			return Down
		}
	case Right:
		{
			return Left
		}
	case Down:
		{
			return Up
		}
	case Left:
		{
			return Right
		}
	}
	panic("direction should be an enum")
}

type Coordinate struct {
	row int
	col int
}

func NewCoordinate(row, col int) Coordinate {
	return Coordinate{row, col}
}

func (c Coordinate) String() string {
	return fmt.Sprintf("{row: %d, col: %d}", c.row, c.col)
}

func (c Coordinate) Get(data [][]byte) byte {
	if !isValid(data, c) {
		return OutBounds
	}
	return data[c.row][c.col]
}

func (c Coordinate) Up() Coordinate {
	return Coordinate{c.row - 1, c.col}
}

func (c Coordinate) Down() Coordinate {
	return Coordinate{c.row + 1, c.col}
}

func (c Coordinate) Left() Coordinate {
	return Coordinate{c.row, c.col - 1}
}

func (c Coordinate) Right() Coordinate {
	return Coordinate{c.row, c.col + 1}
}

func (c Coordinate) Dist(other Coordinate) int {
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

func (c Coordinate) Move(dir Direction) Coordinate {
	switch dir {
	case Up:
		{
			return c.Up()
		}
	case Down:
		{
			return c.Down()
		}
	case Left:
		{
			return c.Left()
		}
	case Right:
		{
			return c.Right()
		}
	default:
		{
			panic(fmt.Sprintf("unknown direction %d", dir))
		}
	}
}

type ChangeRecord struct {
	pos    Coordinate
	origin Coordinate
	dir    Direction
}

type JumpTable [][][]Coordinate

func BuildJumpTable(grid [][]byte) JumpTable {
	jmp := make(JumpTable, len(grid))
	visited := make([][][]bool, len(grid))
	for i := range grid {
		jmp[i] = make([][]Coordinate, len(grid[i]))
		visited[i] = make([][]bool, len(grid[i]))
		for j := range grid[i] {
			jmp[i][j] = make([]Coordinate, 4)
			visited[i][j] = make([]bool, 4)
		}
	}
	for i := range len(jmp) {
		for j := range len(jmp[i]) {
			co := NewCoordinate(i, j)
			jmp.UpdateDirection(grid, co, Up, visited)
			jmp.UpdateDirection(grid, co, Right, visited)
			jmp.UpdateDirection(grid, co, Down, visited)
			jmp.UpdateDirection(grid, co, Left, visited)
		}
	}
	return jmp
}

func (jmp JumpTable) UpdateDirection(data [][]byte, start Coordinate, dir Direction, visited [][][]bool) Coordinate {
	if visited[start.row][start.col][dir] {
		return jmp[start.row][start.col][dir]
	}
	next := start.Move(dir)
	if !isValid(data, next) {
		visited[start.row][start.col][dir] = true
		jmp[start.row][start.col][dir] = next
		return next
	}
	if next.Get(data) == Obstacle {
		visited[start.row][start.col][dir] = true
		jmp[start.row][start.col][dir] = start
		return start
	}
	visited[start.row][start.col][dir] = true
	jmp[start.row][start.col][dir] = jmp.UpdateDirection(data, next, dir, visited)
	return jmp[start.row][start.col][dir]
}

func (jmp JumpTable) AddObstacle(grid [][]byte, center Coordinate) []ChangeRecord {
    changeSet := make([]ChangeRecord, 0, 256)

    // We're going to pretend grid[center] = Obstacle here,
    // but we don't need to modify grid itself if we update jmp.
    for _, dir := range []Direction{Up, Right, Down, Left} {
        neighbor := center.Move(dir)
        if !isValid(grid, neighbor) || grid[neighbor.row][neighbor.col] == Obstacle {
            // nothing to update in this direction
            continue
        }
        opposite := dir.Opposite()
        // walk out from the neighbor until you hit an existing obstacle or boundary
        for c := neighbor; isValid(grid, c) && grid[c.row][c.col] != Obstacle; c = c.Move(dir) {
            // record the old jump-target so we can restore later
            old := jmp[c.row][c.col][opposite]
            // patch it to point at 'neighbor' (the cell just before our new obstacle)
            jmp[c.row][c.col][opposite] = neighbor
            changeSet = append(changeSet, ChangeRecord{pos: c, origin: old, dir: opposite})
        }
    }

    return changeSet
}

func (jmp JumpTable) Restore(c []ChangeRecord) {
	for _, change := range c {
		jmp.Set(change.pos, change.origin, change.dir)
	}
}

func (jmp JumpTable) Set(pos, val Coordinate, dir Direction) Coordinate {
	if !isValid(jmp, pos) {
		return Coordinate{}
	}
	origin := jmp[pos.row][pos.col][dir]
	jmp[pos.row][pos.col][dir] = val
	return origin
}

func (jmp JumpTable) Get(start Coordinate, dir Direction) Coordinate {
	if !isValid(jmp, start) {
		return start
	}
	return jmp[start.row][start.col][dir]
}

func (jmp JumpTable) PathFrom(start Coordinate, dir Direction) (Coordinate, []Coordinate) {
	end := jmp.Get(start, dir)
	path := make([]Coordinate, 0, start.Dist(end))
	for end != start && isValid(jmp, start) {
		path = append(path, start)
		start = start.Move(dir)
	}
	return end, path
}

type Simulator struct {
	pos  Coordinate
	grid [][]byte
	jmp  JumpTable
	path []Coordinate
}

func NewSimulator(grid [][]byte) *Simulator {
	jmp := BuildJumpTable(grid)
	pos := getStartPos(grid)
	path := make([]Coordinate, 0, 8192)
	return &Simulator{pos, grid, jmp, path}
}

func (sim *Simulator) RunFullUnsafe() {
	start := sim.pos
	dir := Direction(Up) // this is always the default
	for isValid(sim.grid, start) {
		next, path := sim.jmp.PathFrom(start, dir)
		sim.path = append(sim.path, path...)
		start = next
		dir = dir.Turn()
	}
	return
}

func (sim *Simulator) RunFastUnsafe() {
	start := sim.pos
	sim.path = append(sim.path, start)
	dir := Direction(Up) // this is always the default
	for isValid(sim.grid, start) {
		next := sim.jmp.Get(start, dir)
		start = next
		dir = dir.Turn()
		sim.path = append(sim.path, start)
	}
	return
}

func (sim *Simulator) CountPossibleCyclesNaive() int {
	cnt := 0
	for i := range sim.grid {
		for j := range sim.grid[i] {
			co := Coordinate{i, j}
			if co == sim.pos || sim.grid[i][j] == Obstacle {
				continue
			}
			changeSet := sim.jmp.AddObstacle(co)
			if !sim.Escapes() {
				cnt++
			}
			sim.jmp.Restore(changeSet)
		}
	}

	return cnt
}

func (sim *Simulator) Escapes() bool {
	slowDir, fastDir := Direction(Up), Direction(Up) // this is always the default
	slow, fast := sim.pos, sim.pos
	for isValid(sim.grid, fast) {
		slow = sim.jmp.Get(slow, slowDir)
		slowDir = slowDir.Turn()
		fast = sim.jmp.Get(fast, fastDir)
		fastDir = fastDir.Turn()
		fast = sim.jmp.Get(fast, fastDir)
		fastDir = fastDir.Turn()
		if slow == fast && slowDir == fastDir {
			return false
		}
	}
	return true
}

func getStartPos(b [][]byte) Coordinate {
	for i := range len(b) {
		for j := range len(b[i]) {
			if b[i][j] == Guard {
				return NewCoordinate(i, j)
			}
		}
	}
	panic("getStartPos() No guard found")
}

func isValid[T any](b [][]T, c Coordinate) bool {
	n, m := len(b), len(b[0])
	return c.row >= 0 && c.row < n && c.col >= 0 && c.col < m
}

func main() {
	b, err := os.ReadFile("./cmd/day6/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	begin := time.Now()
	trimmed := bytes.Trim(b, "\n\r\t ")
	lines := bytes.Split(trimmed, []byte("\n"))
	sim := NewSimulator(lines)
	answer := sim.CountPossibleCyclesNaive()
	fmt.Println(answer)
	fmt.Println(time.Now().Sub(begin))
}

var testLines = [][]byte{
	{'.', '.', '.', '#', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '#', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '#', '^', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
	{'.', '.', '.', '.', '.', '.', '.', '.', '.', '.'},
}
