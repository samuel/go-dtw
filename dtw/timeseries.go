package dtw

type TimeSeries interface {
	At(index int) []float64
	Len() int
}

type ScalarTimeSeries []float64

func (ts ScalarTimeSeries) At(index int) []float64 {
	return ts[index : index+1]
}

func (ts ScalarTimeSeries) Len() int {
	return len(ts)
}

type InterleavedTimeSearches struct {
	N int
	V []float64
}

func (ts *InterleavedTimeSearches) At(index int) []float64 {
	o := index * ts.N
	return ts.V[o : o+ts.N]
}

func (ts *InterleavedTimeSearches) Len() int {
	return len(ts.V) / ts.N
}

type VectorTimeSeries [][]float64

func (ts VectorTimeSeries) At(index int) []float64 {
	return ts[index]
}

func (ts VectorTimeSeries) Len() int {
	return len(ts)
}

func Downsample(ts TimeSeries, factor int) TimeSeries {
	if ts.Len() == 0 {
		return ts
	}
	switch t := ts.(type) {
	case ScalarTimeSeries:
		m := len(t) / factor
		if len(t)%factor != 0 {
			m++
		}
		out := make([]float64, m)
		sum := 0.0
		for i, v := range t {
			sum += v
			if i%factor == factor-1 {
				out[i/factor] = sum / float64(factor)
				sum = 0.0
			}
		}
		if n := len(t) % factor; n != 0 {
			out[len(out)-1] = sum / float64(n)
		}
		return ScalarTimeSeries(out)
	}
	n := len(ts.At(0))
	m := ts.Len() / factor * n
	if ts.Len()%factor != 0 {
		m += n
	}
	out := make([]float64, m)
	sum := out
	for i := 0; i < ts.Len(); i++ {
		for j, v := range ts.At(i) {
			sum[j] += v
		}
		if i%factor == factor-1 {
			for j := 0; j < n; j++ {
				sum[j] /= float64(factor)
			}
			sum = out[n*(i+1)/factor:]
		}
	}
	if n := ts.Len() % factor; n != 0 {
		for i := 0; i < len(sum); i++ {
			out[i] /= float64(n)
		}
	}
	return &InterleavedTimeSearches{N: n, V: out}
}
