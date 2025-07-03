package main

import (
	"bytes"
	"log"
	"os"
)

type Slicer[T any] interface {
	Slice(size int) T
}

type Block struct {
	id   int
	data []byte
	free []byte
	next *Block
	prev *Block
}

func main() {
	fs, err := os.ReadFile("./cmd/day8/test_input.txt")
	if err != nil {
		log.Fatal(err)
	}
	fs = bytes.Trim(fs, "\n\r\t ")
	// left, right := 0, len(blocks)-1
	// for left < right {
	// 	b1, b2 := &blocks[left], &blocks[right]
	// 	if len(b2.data) > b1.free {
	// 		// if b2 is bigger - we fill b1
	// 		b1.data = append(b1.data, b2.data[0:b1.free]...)
	// 		b2.data = b2.data[b1.free:]
	// 		b1.free = 0
	// 	} else {
	// 		// else we consume the entire block of b2
	// 		b1.data = append(b1.data, b2.data...)
	// 		b1.free -= len(b2.data)
	// 		b2.data = nil
	// 	}
	// 	if b1.free == 0 {
	// 		left++
	// 	}
	// 	if len(b2.data) == 0 {
	// 		right--
	// 	}
	// }

	// final := make([]int, 0, 1024)
	// for _, block := range blocks {
	// 	final = append(final, block.data...)
	// }

	// checksum := 0
	// for i, n := range final {
	// 	checksum += i * n
	// }

	// fmt.Println("final layout", final)
	// fmt.Println("checksum", checksum)
}
