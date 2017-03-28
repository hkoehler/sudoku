// (C) Copyright 2017, Heiko Koehler

package sudoku

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

const Size = 9
const BoxSize = 3

type Board struct {
	// state of board
	board [Size][Size]uint8
	// how many solutions have been found
	numSolutions int
	// how many solutions to find maximum
	maxNumSolutions int
}

func NewBoard() *Board {
	return &Board{}
}

func (b *Board) Read(data []byte) error {
	return json.Unmarshal(data, &b.board)
}

func (b *Board) Write() ([]byte, error) {
	return json.Marshal(&b.board)
}

func (b *Board) String() (s string) {
	for r := 0; r < Size; r++ {
		s += fmt.Sprintf("%v\n", b.board[r])
	}
	return
}

// bit mask of invalid/taken values
func (b *Board) InvalidValues(row, col int) (mask uint32) {
	for r := 0; r < Size; r++ {
		mask |= (1 << b.board[r][col])
	}
	for c := 0; c < Size; c++ {
		mask |= (1 << b.board[row][c])
	}
	boxRow := (row / BoxSize) * BoxSize
	boxCol := (col / BoxSize) * BoxSize
	for r := boxRow; r < boxRow+BoxSize; r++ {
		for c := boxCol; c < boxCol+BoxSize; c++ {
			mask |= (1 << b.board[r][c])
		}
	}
	return
}

func (b *Board) Check() bool {
	var expMask uint32

	for i := 1; i <= Size; i++ {
		expMask |= 1 << uint8(i)
	}
	for r := 0; r < Size; r++ {
		for c := 0; c < Size; c++ {
			mask := b.InvalidValues(r, c)
			if mask != expMask {
				return false
			}
		}
	}
	return true
}

func (b *Board) solve(i int) bool {
	if i == Size*Size {
		b.numSolutions++
		if b.numSolutions == b.maxNumSolutions {
			return true
		}
		return false
	}
	row := i / Size
	col := i % Size
	if b.board[row][col] > 0 {
		return b.solve(i + 1)
	}

	mask := b.InvalidValues(row, col)
	for val := uint8(1); val <= uint8(Size); val++ {
		if mask&(1<<val) == 0 {
			b.board[row][col] = val
			if b.solve(i + 1) {
				return true
			}
			b.board[row][col] = 0
		}
	}
	return false
}

// DFS search
func (b *Board) Solve(n int) int {
	b.maxNumSolutions = n
	b.solve(0)
	return b.numSolutions
}

// Generate board with given number of filled squares
func (b *Board) Generate(seed int64) {
	total := 0
	rand.Seed(seed)
	for total < 10 {
		// pick random square and value
		r := rand.Int() % Size
		c := rand.Int() % Size
		v := rand.Int()%Size + 1
		if b.board[r][c] > 0 {
			continue
		}
		mask := b.InvalidValues(r, c)
		if mask&(1<<uint8(v)) != 0 {
			continue
		}
		b.board[r][c] = uint8(v)
		total++
	}
	b.Solve(1)
}
