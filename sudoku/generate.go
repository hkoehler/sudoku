// (C) Copyright 2017, Heiko Koehler

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hkoehler/sudoku"
)

func createOutputDirectory(path string) error {
	return os.MkdirAll(path, 0766)
}

func writeSudoku(path string, seed int64, s *sudoku.Sudoku) error {
	outputFile := filepath.Join(path, fmt.Sprintf("%d", seed))

	if buf, err := s.Write(); err != nil {
		return err
	} else if f, err := os.OpenFile(outputFile, os.O_CREATE|os.O_RDWR, 0666); err != nil {
		return err
	} else if _, err := f.Write(buf); err != nil {
		return err
	}
	return nil
}

func generate(path string, startSeed, endSeed int64) error {
	for seed := startSeed; seed < endSeed; seed++ {
		s := sudoku.NewSudoku(seed, seed)
		if err := writeSudoku(path, seed, s); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	var start, end int64
	var path string

	flag.Int64Var(&start, "start", 0, "first seed")
	flag.Int64Var(&end, "end", 0, "last seed")
	flag.StringVar(&path, "path", "", "output directory")
	flag.Parse()

	if err := createOutputDirectory(path); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
	if err := generate(path, start, end); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
