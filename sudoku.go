// (C) Copyright 2017, Heiko Koehler

package sudoku

import (
	"encoding/json"
	"fmt"
	"math/rand"
)

const Size = 9
const BoxSize = 3

const (
	EASY   = 20
	MEDIUM = 40
	HARD   = 50
	EVIL   = 55
)

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
	b.numSolutions = 0
	b.maxNumSolutions = n
	b.solve(0)
	return b.numSolutions
}

// Deep copy
func (b *Board) Copy() *Board {
	b2 := NewBoard()
	for r := 0; r < Size; r++ {
		copy(b2.board[r][:], b.board[r][:])
	}
	return b2
}

// Punch n holes in board
// Try to find only single solution to puzzle.
func (b *Board) Punch(n int) {
	rand.Seed(0)
	for {
		b1 := b.Copy()
		// punch n holes
		total := 0
		for total < n {
			// pick random square and value
			r := rand.Int() % Size
			c := rand.Int() % Size
			if b1.board[r][c] == 0 {
				continue
			}
			b1.board[r][c] = 0
			total++
		}
		b2 := b1.Copy()
		numSolutions := b2.Solve(2)
		switch numSolutions {
		case 0:
			continue
		case 1:
			*b = *b1
			return
		case 2:
			// try punchin holes again and solve it
			continue
		}
	}
}

// Shuffle numbers to create equal puzzle
func (b *Board) Shuffle(seed int64) {
	var mapping [Size + 1]uint8

	rand.Seed(seed)
	// map i to i by default
	for i := uint8(1); i <= Size; i++ {
		mapping[i] = i
	}
	// switch numbers randomly
	for i := 0; i < 10; i++ {
		i1 := rand.Int()%Size + 1
		i2 := rand.Int()%Size + 1
		mapping[i1], mapping[i2] = mapping[i2], mapping[i1]
	}
	for r := 0; r < Size; r++ {
		for c := 0; c < Size; c++ {
			b.board[r][c] = mapping[b.board[r][c]]
		}
	}
}

// Generate board with given number of filled squares
func (b *Board) Generate(seed int64) {
	total := 0
	rand.Seed(seed)
	for total < 5 {
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

// Sudoku comprises multiple puzzzles with same solution
type Sudoku struct {
	Solution *Board
	Easy     *Board
	Medium   *Board
	Hard     *Board
	Evil     *Board
}

// NewSudoku generates puzzle at various difficutly levels with same solution.
func NewSudoku(seed int64, shuffle int64) *Sudoku {
	b := NewBoard()
	b.Generate(seed)
	b.Shuffle(shuffle)
	s := b.Copy()
	b.Punch(EASY)
	easy := b.Copy()
	b.Punch(MEDIUM - EASY)
	medium := b.Copy()
	b.Punch(HARD - MEDIUM)
	hard := b.Copy()
	//b.Punch(EVIL - HARD)
	//evil := b.Copy()
	return &Sudoku{
		Solution: s,
		Easy:     easy,
		Medium:   medium,
		Hard:     hard,
		//Evil:     evil,
	}
}
