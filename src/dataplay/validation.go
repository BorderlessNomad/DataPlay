package main

// import (
// 	"math"
// )

// given a small fraction of ratings there is a strong (95%) chance that the "real", final positive rating will be this value
// eg: gives expected (not necessarily current as there may have only been a few votes so far) value of positive ratings / total ratings
// func Ranking(id int) float64 {
// 	cor := Correlation{}
// 	err := DB.Where("id= ?", id).Find(&cor).Error
// 	check(err)
// 	pos := float64(cor.Valid)
// 	tot := float64(cor.Valid + cor.Invalid)

// 	if tot == 0 {
// 		return 0
// 	}

// 	z := 1.96
// 	phat := pos / tot
// 	cor.Rating = (phat + z*z/(2*tot) - z*math.Sqrt((phat*(1-phat)+z*z/(4*tot))/tot)) / (1 + z*z/tot)
// 	err = DB.Save(&cor).Error
// 	check(err)
// 	return cor.Rating
// }

// // increment user validated total for correlation
// func ValidateChart(id int) {
// 	cor := Correlation{}
// 	err := DB.Where("id= ?", id).Find(&cor).Error
// 	check(err)
// 	err = DB.Model(&cor).Update("valid", cor.Valid+1).Error
// 	check(err)
// }

// // increment user Invalided total for correlation
// func InvalidateChart(id int) {
// 	cor := Correlation{}
// 	err := DB.Where("id= ?", id).Find(&cor).Error
// 	check(err)
// 	err = DB.Model(&cor).Update("invalid", cor.Invalid+1).Error
// 	check(err)
// }

// // increment user validated total for correlation
// func ValidateObservation(id int) {
// 	cor := Correlation{}
// 	err := DB.Where("id= ?", id).Find(&cor).Error
// 	check(err)
// 	err = DB.Model(&cor).Update("valid", cor.Valid+1).Error
// 	check(err)
// }

// // increment user Invalided total for correlation
// func InvalidateObservation(id int) {
// 	cor := Correlation{}
// 	err := DB.Where("id= ?", id).Find(&cor).Error
// 	check(err)
// 	err = DB.Model(&cor).Update("invalid", cor.Invalid+1).Error
// 	check(err)
// }
