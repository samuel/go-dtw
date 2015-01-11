package dtw

import "testing"

func TestFastDTW(t *testing.T) {
	ts1, err := readSeries("trace0.csv")
	if err != nil {
		t.Fatal(err)
	}
	ts2, err := readSeries("trace1.csv")
	if err != nil {
		t.Fatal(err)
	}

	_, optimalDist := DTW(ts1, ts2, EuclideanDistance)

	for _, rad := range []int{1, 5, 10, 20, 30} {
		_, approxDist := FastDTW(ts1, ts2, rad, EuclideanDistance)
		t.Logf("Radius %d: %f (%.2f%% error)", rad, approxDist, 100.0*(approxDist-optimalDist)/optimalDist)
	}
}

func TestFastDTW1000(t *testing.T) {
	ts1 := sineWave(1000, 1, 0.5)
	ts2 := sineWave(1500, 1, 0.6)

	_, optimalDist := DTW(ts1, ts2, EuclideanDistance)

	for _, rad := range []int{1, 5, 10, 20, 30} {
		_, approxDist := FastDTW(ts1, ts2, rad, EuclideanDistance)
		t.Logf("Radius %d: %f (%.2f%% error)", rad, approxDist, 100.0*(approxDist-optimalDist)/optimalDist)
	}
}

func benchmarkFastDTW(b *testing.B, size, r int) {
	ts1 := sineWave(size, 1, 0.5)
	ts2 := sineWave(size+size/2, 1, 0.6)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		FastDTW(ts1, ts2, r, ManhattanDistance)
	}
}

func BenchmarkFastDTW1000Radius1(b *testing.B) {
	benchmarkFastDTW(b, 1000, 1)
}

func BenchmarkFastDTW1000Radius5(b *testing.B) {
	benchmarkFastDTW(b, 1000, 5)
}

func BenchmarkFastDTW1000Radius10(b *testing.B) {
	benchmarkFastDTW(b, 1000, 10)
}

func BenchmarkFastDTW1000Radius20(b *testing.B) {
	benchmarkFastDTW(b, 1000, 20)
}

func BenchmarkFastDTW10000Radius20(b *testing.B) {
	benchmarkFastDTW(b, 10000, 20)
}
