package dtw

func FastDTW(ts1, ts2 TimeSeries, searchRadius int, distFunc DistanceFunc) ([]Point, float64) {
	if searchRadius < 0 {
		searchRadius = 0
	}
	minTSSize := searchRadius + 2

	if ts1.Len() <= minTSSize || ts2.Len() <= minTSSize {
		return DTW(ts1, ts2, distFunc)
	}

	factor := 2
	shrunk1 := Downsample(ts1, factor)
	shrunk2 := Downsample(ts2, factor)

	path, _ := FastDTW(shrunk1, shrunk2, searchRadius, distFunc)
	window := ExpandedResWindow(ts1, ts2, factor, path, searchRadius)
	return Constrained(ts1, ts2, window, nil, distFunc)
}

func ExpandedResWindow(ts1, ts2 TimeSeries, factor int, path []Point, searchRadius int) Window {
	w := NewCustomWindow(ts1.Len(), ts2.Len())

	pt := path[0]
	lastWarped := Point{X: ts1.Len(), Y: ts2.Len()}

	for _, warped := range path {
		if warped.X > lastWarped.X {
			pt.X += factor
		}
		if warped.Y > lastWarped.Y {
			pt.Y += factor
		}
		if warped.X > lastWarped.X && warped.Y > lastWarped.Y {
			w.Mark(pt.X-1, pt.Y)
			w.Mark(pt.X, pt.Y-1)
		}
		for y := 0; y < factor; y++ {
			// Only need 2 points since the window fills out the row
			// w.Mark(pt.X, pt.Y+y)
			// w.Mark(pt.X+factor-1, pt.Y+y)

			if y := pt.Y + y; y < ts2.Len() {
				w.Mark(pt.X, y)
				if x := pt.X + factor - 1; x < ts1.Len() {
					w.Mark(x, y)
				} else {
					w.Mark(x-1, y)
				}
			}
		}
		lastWarped = warped
	}

	w.Expand(searchRadius)

	return w
}
