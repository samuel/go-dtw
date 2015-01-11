package dtw

import "math"

// WarpDistance is a more memory efficient O(N) way to calculate the warp
// distance since it only needs to keep 2 columns instead of N*M costs O(N^2)
func WarpDistance(ts1, ts2 TimeSeries, distFunc DistanceFunc) float64 {
	n1 := ts1.Len()
	n2 := ts2.Len()
	mem := make([]float64, n1*2)
	lastRow := mem[n1:]
	curRow := mem[:n1]
	curRow[0] = distFunc(ts1.At(0), ts2.At(0))
	for x := 1; x < n1; x++ {
		curRow[x] = curRow[x-1] + distFunc(ts1.At(x), ts2.At(0))
	}
	for y := 1; y < n2; y++ {
		curRow, lastRow = lastRow, curRow
		curRow[0] = lastRow[0] + distFunc(ts1.At(0), ts2.At(y))
		for x := 1; x < n1; x++ {
			minCost := min(curRow[x-1], min(lastRow[x], lastRow[x-1]))
			curRow[x] = minCost + distFunc(ts1.At(x), ts2.At(y))
		}
	}
	return curRow[len(curRow)-1]
}

// DTW performs dynamic time warping on the given timeseries and returns
// the path as well as the minimum cost.
func DTW(ts1, ts2 TimeSeries, distFunc DistanceFunc) ([]Point, float64) {
	n1 := ts1.Len()
	n2 := ts2.Len()
	grid := make([]float64, n1*n2)

	grid[0] = distFunc(ts1.At(0), ts2.At(0))
	for x := 1; x < n1; x++ {
		grid[x] = grid[x-1] + distFunc(ts1.At(x), ts2.At(0))
	}
	off := n1
	for y := 1; y < n2; y++ {
		grid[off] = grid[off-n1] + distFunc(ts1.At(0), ts2.At(y))
		off++
		for x := 1; x < n1; x++ {
			minCost := min(grid[off-1], min(grid[off-n1], grid[off-n1-1]))
			grid[off] = minCost + distFunc(ts1.At(x), ts2.At(y))
			off++
		}
	}

	path := make([]Point, 0)
	x, y := n1-1, n2-1
	path = append(path, Point{X: x, Y: y})
	for x > 0 || y > 0 {
		o := y*n1 + x
		diag := math.Inf(1)
		left := math.Inf(1)
		down := math.Inf(1)
		if x > 0 && y > 0 {
			diag = grid[o-n1-1]
		}
		if x > 0 {
			left = grid[o-1]
		}
		if y > 0 {
			down = grid[o-n1]
		}
		switch {
		case diag <= left && diag <= down:
			x--
			y--
		case left < diag && left < down:
			x--
		case down < diag && down < left:
			y--
		// Move towards the diagnal if all equal
		case x <= y:
			x--
		default:
			y--
		}
		path = append(path, Point{X: x, Y: y})
	}

	// Reverse the path
	for i := 0; i < len(path)/2; i++ {
		j := len(path) - i - 1
		path[i], path[j] = path[j], path[i]
	}

	return path, grid[n1*n2-1]
}

func Constrained(ts1, ts2 TimeSeries, window Window, grid Matrix, distFunc DistanceFunc) ([]Point, float64) {
	rect := window.Rect()
	if grid == nil {
		grid = NewMatrixFromWindow(window)
	}

	r := window.Range(rect.MinY)
	grid.Set(r.Min, rect.MinY, distFunc(ts1.At(r.Min), ts2.At(rect.MinY)))
	for x := r.Min + 1; x <= r.Max; x++ {
		grid.Set(x, rect.MinY, grid.Get(x-1, rect.MinY)+distFunc(ts1.At(x), ts2.At(rect.MinY)))
	}
	for y := rect.MinY + 1; y <= rect.MaxY; y++ {
		r := window.Range(y)
		lastCost := grid.Get(r.Min, y-1)
		if lastCost == inf {
			// lastCost = 0.0
			panic("HMM")
		}
		grid.Set(r.Min, y, lastCost+distFunc(ts1.At(r.Min), ts2.At(y)))
		for x := r.Min + 1; x <= r.Max; x++ {
			minCost := min(grid.Get(x-1, y), min(grid.Get(x, y-1), grid.Get(x-1, y-1)))
			if minCost == inf {
				panic("WTF?")
			}
			grid.Set(x, y, minCost+distFunc(ts1.At(x), ts2.At(y)))
		}
	}

	path := make([]Point, 0)
	x, y := rect.MaxX, rect.MaxY
	path = append(path, Point{X: x, Y: y})
	for x > 0 || y > 0 {
		diag := grid.Get(x-1, y-1)
		left := grid.Get(x-1, y)
		down := grid.Get(x, y-1)
		switch {
		case diag <= left && diag <= down:
			x--
			y--
		case left < diag && left < down:
			x--
		case down < diag && down < left:
			y--
		// Move towards the diagnal if all equal
		case x <= y:
			x--
		default:
			y--
		}
		path = append(path, Point{X: x, Y: y})
	}

	// Reverse the path
	for i := 0; i < len(path)/2; i++ {
		j := len(path) - i - 1
		path[i], path[j] = path[j], path[i]
	}

	return path, grid.Get(rect.MaxX, rect.MaxY)
}
