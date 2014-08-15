package main

import (
	"math"
	"time"
)

// given a small fraction of ratings there is a strong (95%) chance that the "real", final positive rating will be this value
// eg: gives expected (not necessarily current as there may have only been a few votes so far) value of positive ratings / total ratings
func RankValidations(valid int, invalid int) float64 {
	pos := float64(valid)
	tot := float64(valid + invalid)

	if tot == 0 {
		return 0
	}

	z := 1.96
	phat := pos / tot
	result := (phat + z*z/(2*tot) - z*math.Sqrt((phat*(1-phat)+z*z/(4*tot))/tot)) / (1 + z*z/tot)
	return result
}

// increment user validated total for chart and rerank, id = 0 for new validation
func ValidateChart(chartid int, correlated bool, json []byte, originid int, uid int) {
	val := Validated{}
	vld := Validation{}

	if chartid == 0 {
		val.DiscoveredBy = uid
		val.DateDiscovered = time.Now()
		val.Correlated = correlated
		val.Rating = RankValidations(1, 0)
		val.Valid = 1
		val.Invalid = 0
		val.Json = json
		val.OriginId = originid

		err := DB.Save(&val).Error
		check(err)

	} else {
		err := DB.Where("patternid= ?", chartid).Find(&val).Error
		check(err)

		val.Valid++
		val.Rating = RankValidations(val.Valid, val.Invalid)

		err = DB.Save(&val).Error
		check(err)

		vld.PatternId = chartid
		vld.ValidatedBy = uid
		vld.ValidationType = "chart"

		err = DB.Save(&vld).Error
		check(err)
	}
}

// increment user invalidated total for chart and rerank, id = 0 for new validation
func InvalidateChart(chartid int, correlated bool, json []byte, originid int, uid int) {
	val := Validated{}
	vld := Validation{}

	if chartid == 0 {
		val.DiscoveredBy = uid
		val.DateDiscovered = time.Now()
		val.Correlated = correlated
		val.Rating = 0
		val.Valid = 0
		val.Invalid = 1
		val.Json = json
		val.OriginId = originid

		err := DB.Save(&val).Error
		check(err)

	} else {
		err := DB.Where("patternid= ?", chartid).Find(&val).Error
		check(err)

		val.Invalid++
		val.Rating = RankValidations(val.Valid, val.Invalid)

		err = DB.Save(&val).Error
		check(err)

		vld.PatternId = chartid
		vld.ValidatedBy = uid
		vld.ValidationType = "chart"

		err = DB.Save(&vld).Error
		check(err)
	}
}

// increment user validated total for observation and rerank, id = 0 for new observation
func ValidateObservation(obsid int, text string, patternid int, uid int, coordinates string) {
	obs := Observation{}
	vld := Validation{}

	if obsid == 0 {
		obs.Text = text
		obs.PatternId = patternid
		obs.DiscoveredBy = uid
		obs.Coordinates = coordinates
		obs.Rating = RankValidations(1, 0)
		obs.Valid = 1
		obs.Invalid = 0
		obs.DateDiscovered = time.Now()

		err := DB.Save(&obs).Error
		check(err)

	} else {
		err := DB.Where("patternid= ?", obsid).Find(&obs).Error
		check(err)

		obs.Valid++
		obs.Rating = RankValidations(obs.Valid, obs.Invalid)

		err = DB.Save(&obs).Error
		check(err)

		vld.PatternId = obsid
		vld.ValidatedBy = uid
		vld.ValidationType = "observation"

		err = DB.Save(&vld).Error
		check(err)
	}
}

// increment user Invalidated total for correlation, id = 0 for new observation
func InvalidateObservation(obsid int, text string, patternid int, uid int, coordinates string) {
	obs := Observation{}
	vld := Validation{}

	if obsid == 0 {
		obs.Text = text
		obs.PatternId = patternid
		obs.DiscoveredBy = uid
		obs.Coordinates = coordinates
		obs.Rating = 0
		obs.Valid = 0
		obs.Invalid = 1
		obs.DateDiscovered = time.Now()

		err := DB.Save(&obs).Error
		check(err)

	} else {
		err := DB.Where("patternid= ?", obsid).Find(&obs).Error
		check(err)

		obs.Invalid++
		obs.Rating = RankValidations(obs.Valid, obs.Invalid)

		err = DB.Save(&obs).Error
		check(err)

		vld.PatternId = obsid
		vld.ValidatedBy = uid
		vld.ValidationType = "observation"

		err = DB.Save(&vld).Error
		check(err)
	}
}
