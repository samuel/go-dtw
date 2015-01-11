package dtw

import "math"

type DistanceFunc func(v1, v2 []float64) float64

func EuclideanDistance(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		return math.NaN()
	}
	if len(v1) == 1 {
		return math.Abs(v1[0] - v2[0])
	}
	sum := 0.0
	for i, a := range v1 {
		d := (a - v2[i])
		sum += d * d
	}
	return math.Sqrt(sum)
}

func ManhattanDistance(v1, v2 []float64) float64 {
	if len(v1) != len(v2) {
		return math.NaN()
	}
	sum := 0.0
	for i, a := range v1 {
		sum += math.Abs(a - v2[i])
	}
	return sum
}
