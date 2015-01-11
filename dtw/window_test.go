package dtw

import (
	"math/rand"
	"os"
	"testing"
)

func TestWindow(t *testing.T) {
	w := &SakoeChibaBandWindow{
		Size:  20,
		Width: 2,
	}
	DisplayWindow(os.Stdout, w)
}

func TestCustomWindow(t *testing.T) {
	rand.Seed(2)
	w := NewCustomWindow(20, 20)
	x, y := 0, 0
	w.Mark(0, 0)
	for x != 19 || y != 19 {
		if x < 19 {
			x += rand.Int() & 1
		}
		if y < 19 {
			y += rand.Int() & 1
		}
		w.Mark(x, y)
	}
	w.Expand(2)
	DisplayWindow(os.Stdout, w)
}
