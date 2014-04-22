package main

import (
	"github.com/skelterjohn/go.matrix" // daa59528eefd43623a4c8e36373a86f9eef870a2
)

var degree = 2

func GetPolyResults(xGiven []float64, yGiven []float64) []float64 {
	m := len(yGiven)
	if m != len(xGiven) {
		return []float64{0, 0, 0} // Send it back, There is nothing sane here.
	}
	if m < 5 {
		// Prevent the processing of really small datasets, This is becauase there
		// appears to be a bug in the libary that will trigger a crash in the go.matrix
		// if some (small) amount of values are entered. I don't know why this happens
		// (Otherwise I would have fixed it) but the URL for the github issue is:
		// https://github.com/skelterjohn/go.matrix/issues/11
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
		Logger.Println(err)
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
	Logger.Println(c)
	return c
}
