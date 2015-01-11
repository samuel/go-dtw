package dtw

import (
	"bufio"
	"io"
)

type Range struct {
	Min, Max int
}

func (r Range) Size() int {
	return r.Max - r.Min + 1
}

type Rect struct {
	MinX, MaxX int
	MinY, MaxY int
}

type Window interface {
	Range(y int) Range
	Rect() Rect
	Area() int
}

type SakoeChibaBandWindow struct {
	Size  int
	Width int
}

func (w *SakoeChibaBandWindow) Range(y int) Range {
	return Range{Min: maxI(y-w.Width*2, 0), Max: minI(y+w.Width*2, w.Size-1)}
}

func (w *SakoeChibaBandWindow) Rect() Rect {
	return Rect{MinX: 0, MaxX: w.Size - 1, MinY: 0, MaxY: w.Size - 1}
}

func (w *SakoeChibaBandWindow) Area() int {
	area := 0
	for y := 0; y < w.Size; y++ {
		area += w.Range(y).Size()
	}
	return area
}

type NullWindow struct {
	Width, Height int
}

func (w *NullWindow) Range(y int) Range {
	return Range{Min: 0, Max: w.Width - 1}
}

func (w *NullWindow) Rect() Rect {
	return Rect{MinX: 0, MinY: 0, MaxX: w.Width - 1, MaxY: w.Height - 1}
}

func (w *NullWindow) Area() int {
	return w.Width * w.Height
}

type CustomWindow struct {
	width  int
	height int
	rows   []Range
}

func NewCustomWindow(width, height int) *CustomWindow {
	rows := make([]Range, height)
	for i := 0; i < height; i++ {
		rows[i] = Range{Min: 0, Max: -1}
	}
	return &CustomWindow{
		width:  width,
		height: height,
		rows:   rows,
	}
}

func (w *CustomWindow) Range(y int) Range {
	return w.rows[y]
}

func (w *CustomWindow) Rect() Rect {
	return Rect{MinX: 0, MinY: 0, MaxX: w.width - 1, MaxY: w.height - 1}
}

func (w *CustomWindow) Area() int {
	// return w.area
	area := 0
	for y := 0; y < w.height; y++ {
		area += w.rows[y].Size()
	}
	return area
}

func (w *CustomWindow) Mark(x, y int) {
	r := w.rows[y]
	if r.Max == -1 {
		r = Range{Min: x, Max: x}
	} else {
		if x < r.Min {
			r.Min = x
		}
		if x > r.Max {
			r.Max = x
		}
	}
	w.rows[y] = r
}

func (w *CustomWindow) Expand(radius int) {
	if radius <= 0 {
		return
	}

	rows := make([]Range, w.height)
	copy(rows, w.rows)
	for y, r := range w.rows {
		r.Min -= radius
		if r.Min < 0 {
			r.Min = 0
		}
		r.Max += radius
		if r.Max >= w.width {
			r.Max = w.width - 1
		}
		rows[y] = r

		for i := -radius; i <= radius; i++ {
			if i == 0 {
				continue
			}
			y2 := i + y
			if y2 < 0 || y2 >= w.height {
				continue
			}
			r2 := rows[y2]
			if r2.Min > r.Min {
				r2.Min = r.Min
			}
			if r2.Max < r.Max {
				r2.Max = r.Max
			}
			rows[y2] = r2
		}
	}
	w.rows = rows
}

func DisplayWindow(wr io.Writer, win Window) error {
	r := win.Rect()
	bw := bufio.NewWriter(wr)
	for y := r.MinY; y <= r.MaxY; y++ {
		rng := win.Range(y)
		for x := r.MinX; x <= r.MaxX; x++ {
			c := byte(' ')
			if x < rng.Min || x > rng.Max {
				c = '#'
			}
			if err := bw.WriteByte(c); err != nil {
				return err
			}
		}
		if err := bw.WriteByte('\n'); err != nil {
			return err
		}
	}
	return bw.Flush()
}
