package main

import (
	"math"
	"sort"
)

type Ordinal struct {
	Original int     //keep original order for ordinal data
	DataVal  float64 // value of the data
	RankVal  float64 // calculated rank value
}

/**
 * @brief calculates Pearson coefficient
 * @details calculates Pearson product-moment correlation coefficient for two interval/ratio data sets of equal size
 *
 * @param float64 arrays of x & y values
 * @return correlation value
 */
func Pearson(x []float64, y []float64) float64 {
	sumx, sumy, sumxSq, sumySq, pSum := 0.0, 0.0, 0.0, 0.0, 0.0
	n := float64(len(x))
	if n == 0 || n != float64(len(y)) {
		return 0
	}

	for i, _ := range x {
		sumx += x[i]
		sumy += y[i]
		sumxSq += math.Pow(x[i], 2)
		sumySq += math.Pow(y[i], 2)
		pSum += x[i] * y[i]
	}

	num := pSum - (sumx * sumy / n)
	den := math.Pow(((sumxSq - math.Pow(sumx, 2)/n) * (sumySq - math.Pow(sumy, 2)/n)), 0.5)

	if den == 0 {
		return 0
	}
	return num / den
}

/**
 * @brief calculates the correlation between two data sets with a common divisor
 * @details
 *
 * @param float64 arrays of x, y & z values
 * @return correlation value
 */
func Spurious(x []float64, y []float64, z []float64) float64 {
	vx := math.Pow(Variation(x), 2)
	vy := math.Pow(Variation(y), 2)
	v1z := math.Pow(1/Variation(z), 2)
	num := v1z * Sgn(Mean(x)) * Sgn(Mean(y))
	den := math.Sqrt((vx*(1+v1z) + v1z) * (vy*(1+v1z) + v1z))
	r := num / den

	if math.IsNaN(r) {
		return 0
	}
	return r
}

/**
 * @brief calculates Spearman's rank correlation coefficient for ranked, ordinal data]
 * @details
 *
 * @param float64 arrays of x & y values
 * @return correlation value
 */
func Spearman(x []float64, y []float64) float64 {
	n := len(x)
	rx := make([]Ordinal, n)
	ry := make([]Ordinal, n)

	// copy data values and original place value of ordered data into array
	for i, _ := range x {
		rx[i].DataVal = x[i]
		rx[i].Original = i
		ry[i].DataVal = y[i]
		ry[i].Original = i
	}

	// now we have original place values of data stored, we can temporarily reorder data in order to be ranked easily
	sort.Sort(sort.Reverse(ByData(rx)))
	sort.Sort(sort.Reverse(ByData(ry)))
	tiedx := Rank(rx)
	tiedy := Rank(ry)

	// put data back in original order
	sort.Sort(ByOrig(rx))
	sort.Sort(ByOrig(ry))

	p := 0.0

	if tiedx == false && tiedy == false { // this is the Pearson formula for use on distinct data values

		sumSqDiff, num, den := 0.0, 0.0, 0.0
		for i := 0; i < n; i++ {
			sumSqDiff += math.Pow(rx[i].RankVal-ry[i].RankVal, 2)
		}
		num = 6 * sumSqDiff
		den = float64(n) * (math.Pow(float64(n), 2) - 1)
		p = 1 - num/den

	} else { // with tied data values pass ranking values to Pearson
		px := make([]float64, n, n)
		py := make([]float64, n, n)
		for i, _ := range rx {
			px[i] = rx[i].RankVal
			py[i] = ry[i].RankVal
		}
		p = Pearson(px, py)
	}
	return p
}

/**
 * @brief calculates the coeficient of variation
 * @details calculates the relative variability (the ratio of the standard deviation to the mean)
 *
 * @param array of float values
 * @return variation value
 */
func Variation(x []float64) float64 {
	standDev := StandDev(x)
	mean := Mean(x)
	return standDev / mean
}

/**
 * @brief calculates the population standard deviation
 * @details (not the sample standard deviation as we are not interested in extrapolating)
 *
 * @param array of float values
 * @return standard deviation value
 */
func StandDev(x []float64) float64 {
	sumx := 0.0
	n := float64(len(x))
	mean := Mean(x)
	for _, v := range x {
		sumx += math.Pow((v - mean), 2)
	}
	return math.Sqrt(sumx / n)
}

/**
 * @brief calculates the mean average
 * @details
 *
 * @param float64 array of values
 * @return mean of values
 */
func Mean(x []float64) float64 {
	n := float64(len(x))
	sumx := 0.0
	for _, v := range x {
		sumx += v
	}
	return sumx / n
}

/**
 * @brief returns sign of value
 * @details
 *
 * @param float64 value
 * @return sign multiplier
 */
func Sgn(a float64) float64 {

	switch {
	case a < 0:
		return -1
	case a > 0:
		return 1
	default:
		return 0
	}
}

/**
 * @brief Determines the value of the rank and wherever data values are tied gives them all the same average rank value
 * @details
 *
 * @param ordinal struct
 * @return true if there are tied values
 */
func Rank(o []Ordinal) bool {
	n := cap(o)
	tied := false

	// determine rank values
	for i := 0; i < n; {
		rank := 1.0 + float64(i)
		j := i
		for j < n-1 && o[j].DataVal == o[j+1].DataVal { // discover if current examined value is duplicated in next value
			tied = true // flag when values are repeated
			j++
			rank += float64(j) + 1
		}
		rank /= float64(j-i) + 1
		o[i].RankVal = rank
		for i < n-1 && o[i].DataVal == o[i+1].DataVal { // give all duplicate values the same average rank
			i++
			o[i].RankVal = rank
		}
		i++
	}
	return tied
}

//Implement sort.Interface for Ordinal struct type based on Original and DataVal fields
type ByOrig []Ordinal
type ByData []Ordinal

func (a ByOrig) Len() int           { return len(a) }
func (a ByOrig) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByOrig) Less(i, j int) bool { return a[i].Original < a[j].Original }
func (a ByData) Len() int           { return len(a) }
func (a ByData) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByData) Less(i, j int) bool { return a[i].DataVal < a[j].DataVal }
