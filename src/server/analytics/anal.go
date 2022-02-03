package analytics

import (
	"fmt"
	"sort"
)

type Anal struct{}

func (a Anal) Show(values []float32) {
	fmt.Printf("Median (%v): %v ms \n", len(values), a.calcMedian(values))
	fmt.Printf("Avg (%v)   : %v ms \n", len(values), a.calcAvg(values))
	fmt.Printf("Min (%v)   : %v ms \n", len(values), a.calcMin(values))
	fmt.Printf("Max (%v)   : %v ms \n", len(values), a.calcMax(values))
}

func (a Anal) calcMedian(n []float32) float32 {
	sort.Slice(n, func(i, j int) bool { return n[i] < n[j] })
	mNumber := len(n) / 2

	if len(n)%2 != 0 {
		return n[mNumber]
	}

	return (n[mNumber-1] + n[mNumber]) / 2
}

func (a Anal) calcAvg(n []float32) float32 {
	var sum float32 = 0

	for _, t := range n {
		sum += t
	}

	return sum / float32(len(n))
}

func (a Anal) calcMin(n []float32) float32 {
	var min float32 = 1000

	for _, t := range n {
		if t < min {
			min = t
		}
	}

	return min
}

func (a Anal) calcMax(n []float32) float32 {
	var max float32 = 0

	for _, t := range n {
		if t > max {
			max = t
		}
	}

	return max
}
