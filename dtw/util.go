package dtw

import (
	"fmt"
	"io"
)

type Point struct {
	X, Y int
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func minI(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxI(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func DisplayPath(wr io.Writer, path []Point) error {
	last := path[len(path)-1]
	width := last.X + 1
	height := last.Y + 1
	grid := make([]byte, width*height)
	for i := 0; i < len(grid); i++ {
		grid[i] = '.'
	}
	for _, p := range path {
		grid[p.Y*width+p.X] = '*'
	}
	for y := 0; y < height; y++ {
		if _, err := fmt.Println(string(grid[y*width : (y+1)*width])); err != nil {
			return err
		}
	}
	return nil
}
