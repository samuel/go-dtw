package dtw

import "math"

var inf = math.Inf(+1)

type Matrix interface {
	Get(x, y int) float64
	Set(x, y int, v float64)
}

type MemoryMatrix struct {
	w        Window
	r        Rect
	values   []float64
	rowIndex []int
}

func NewMatrixFromWindow(w Window) *MemoryMatrix {
	r := w.Rect()
	values := make([]float64, w.Area())
	rowIndex := make([]int, r.MaxY-r.MinY+1)
	offset := 0
	for y := r.MinY; y <= r.MaxY; y++ {
		rowIndex[y-r.MinY] = offset
		offset += w.Range(y).Size()
	}
	return &MemoryMatrix{
		w:        w,
		r:        r,
		values:   values,
		rowIndex: rowIndex,
	}
}

func (m *MemoryMatrix) SetWindow(w Window) {
	m.w = w
	m.r = w.Rect()
	area := w.Area()
	if area > cap(m.values) {
		m.values = make([]float64, area)
	} else {
		m.values = m.values[:area]
	}
	rows := m.r.MaxY - m.r.MinY + 1
	if rows > cap(m.rowIndex) {
		m.rowIndex = make([]int, rows)
	} else {
		m.rowIndex = m.rowIndex[:rows]
	}
	offset := 0
	for y := m.r.MinY; y <= m.r.MaxY; y++ {
		m.rowIndex[y-m.r.MinY] = offset
		offset += w.Range(y).Size()
	}
}

func (m *MemoryMatrix) Get(x, y int) float64 {
	if y < m.r.MinY || y > m.r.MaxY {
		return inf
	}
	r := m.w.Range(y)
	if x < r.Min || x > r.Max {
		return inf
	}
	return m.values[m.rowIndex[y-m.r.MinY]+x-r.Min]
}

func (m *MemoryMatrix) Set(x, y int, v float64) {
	if y < m.r.MinY || y > m.r.MaxY {
		panic("dtw/matrix: attempting to Set out of bounds")
	}
	r := m.w.Range(y)
	if x < r.Min || x > r.Max {
		// Treat this like a slice out of bounds
		panic("dtw/matrix: attempting to Set out of bounds")
	}
	m.values[m.rowIndex[y-m.r.MinY]+x-r.Min] = v
}
