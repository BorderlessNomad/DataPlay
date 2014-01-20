package api

import (
	"fmt"
	"github.com/skelterjohn/go.matrix"
)

// var xGiven = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
// var yGiven = []float64{1, 6, 17, 34, 57, 86, 121, 162, 209, 262, 321}
var degree = 2

func GetPolyResults(xGiven []float64, yGiven []float64) []float64 {
	m := len(yGiven)
	if m != len(xGiven) {
		return []float64{0, 0, 0} // Send it back, There is nothing sane here.
	}
	n := degree + 1
	y := matrix.MakeDenseMatrix(yGiven, m, 1)
	x := matrix.Zeros(m, n)
	for i := 0; i < m; i++ {
		ip := float64(1)
		for j := 0; j < n; j++ {
			x.Set(i, j, ip)
			ip *= xGiven[i]
		}
	}

	q, r := x.QR()
	qty, err := q.Transpose().Times(y)
	if err != nil {
		fmt.Println(err)
		return []float64{0, 0, 0}
	}
	c := make([]float64, n)
	for i := n - 1; i >= 0; i-- {
		c[i] = qty.Get(i, 0)
		for j := i + 1; j < n; j++ {
			c[i] -= c[j] * r.Get(i, j)
		}
		c[i] /= r.Get(i, i)
	}
	fmt.Println(c)
	return c
}
