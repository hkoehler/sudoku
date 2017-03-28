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
