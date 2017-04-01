// (C) Copyright 2017, Heiko Koehler

package sudoku

import "testing"

func TestBasic(t *testing.T) {
	b := NewBoard()
	b.Solve(1)
	t.Logf("\n%s\n", b)
	if b.Check() == false {
		t.Fail()
	}
	if b.Solve(1) != 1 {
		t.Fail()
	}
}

func TestGenerator(t *testing.T) {
	for i := 0; i < 100; i++ {
		b := NewBoard()
		b.Generate(int64(i))
		if b.Check() == false {
			t.Fail()
		}
		t.Logf("\n%s\n", b)
	}
}

func TestPunch(t *testing.T) {
	b := NewBoard()
	b.Generate(42)
	t.Logf("Solution:\n%s\n", b)
	// easy puzzle
	b.Punch(20)
	t.Logf("Easy:\n%s\n", b)
	// medium puzzle
	b.Punch(20)
	t.Logf("Medium:\n%s\n", b)
	// hard puzzle
	b.Punch(10)
	t.Logf("Hard:\n%s\n", b)
	// evil puzzle
	b.Punch(5)
	t.Logf("Evil:\n%s\n", b)
}

func TestShuffle(t *testing.T) {
	b := NewBoard()
	b.Generate(13)
	t.Logf("Solution:\n%s\n", b)
	b.Shuffle(42)
	t.Logf("Shuffled solution:\n%s\n", b)
	if b.Check() == false {
		t.Fail()
	}
}

func TestSudoku(t *testing.T) {
	for i := int64(0); i < 10; i++ {
		s := NewSudoku(i, 42)
		t.Logf("Solution:\n%s\n", s.Solution)
		t.Logf("Easy:\n%s\n", s.Easy)
		t.Logf("Medium:\n%s\n", s.Medium)
		t.Logf("Hard:\n%s\n", s.Hard)
		t.Logf("Evil:\n%s\n", s.Evil)
	}
}
