package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type coef int

const (
	P coef = iota
	S
	V
)

type DateVal struct {
	Date  time.Time
	Value float64
}

type FromTo struct {
	From time.Time
	To   time.Time
}

func GetCorrelation(table1 string, valCol1 string, dateCol1 string, c coef) string {
	if table1 == "" || valCol1 == "" || dateCol1 == "" {
		return ""
	}

	m := make(map[string]string)
	m["table1"], m["dateCol1"], m["valCol1"] = table1, dateCol1, valCol1
	cor := Correlation{}
	result, method := "", ""
	var coef []float64
	nameChk := GetRandomNames(m, c)

	if nameChk {

		if c == P {
			err := DB.Model(&cor).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("method = ?", "Pearson").Pluck("coef", &coef).Error
			check(err)
			method = "Pearson"

		} else if c == S {
			err := DB.Model(&cor).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("tbl3 = ?", m["table3"]).Where("col3 = ?", m["valCol3"]).Where("method = ?", "Spurious").Pluck("coef", &coef).Error
			check(err)
			method = "Spurious"

		} else if c == V {
			err := DB.Model(&cor).Where("tbl1 = ?", m["table1"]).Where("col1 = ?", m["valCol1"]).Where("tbl2 = ?", m["table2"]).Where("col2 = ?", m["valCol2"]).Where("method = ?", "Visual").Pluck("coef", &coef).Error
			check(err)
			method = "Visual"

		} else {

			return ""
		}

		if coef == nil {
			cf := GetCoef(m, c)
			correlation := Correlation{
				Tbl1:   m["table1"],
				Col1:   m["valCol1"],
				Tbl2:   m["table2"],
				Col2:   m["valCol2"],
				Tbl3:   m["table3"],
				Col3:   m["valCol3"],
				Method: method,
				Coef:   cf,
			}

			jv, _ := json.Marshal(correlation)
			correlation.Json = string(jv)
			result = string(jv)
			err := DB.Save(&correlation).Error
			check(err)
		}
	}

	return result
}

func GetCoef(m map[string]string, c coef) float64 {
	if len(m) == 0 {
		return 0.0
	}

	var bucketRange []FromTo
	var cf float64

	if c == P {
		x := ExtractDateVal(m["table1"], m["dateCol1"], m["valCol1"])
		y := ExtractDateVal(m["table2"], m["dateCol2"], m["valCol2"])
		fromX, toX, rngX := DetermineRange(x)
		fromY, toY, rngY := DetermineRange(y)

		if rngX == 0 || rngY == 0 {
			return 0
		}

		bucketRange = GetIntersect(&fromX, &toX, &rngX, fromY, toY, rngY)
		xBuckets := FillBuckets(x, bucketRange)
		yBuckets := FillBuckets(y, bucketRange)
		cf = Pearson(xBuckets, yBuckets)

	} else if c == S {
		x := ExtractDateVal(m["table2"], m["dateCol2"], m["valCol2"])
		y := ExtractDateVal(m["table3"], m["dateCol3"], m["valCol3"])
		z := ExtractDateVal(m["table1"], m["dateCol1"], m["valCol1"])
		fromX, toX, rngX := DetermineRange(x)
		fromY, toY, rngY := DetermineRange(y)
		fromZ, toZ, rngZ := DetermineRange(z)

		if rngX == 0 || rngY == 0 || rngZ == 0 {
			return 0
		}

		_ = GetIntersect(&fromX, &toX, &rngX, fromY, toY, rngY)
		bucketRange = GetIntersect(&fromX, &toX, &rngX, fromZ, toZ, rngZ)
		xBuckets := FillBuckets(x, bucketRange)
		yBuckets := FillBuckets(y, bucketRange)
		zBuckets := FillBuckets(z, bucketRange)
		cf = Spurious(xBuckets, yBuckets, zBuckets)

	} else if c == V {
		x := ExtractDateVal(m["table1"], m["dateCol1"], m["valCol1"])
		y := ExtractDateVal(m["table2"], m["dateCol2"], m["valCol2"])
		fromX, toX, rngX := DetermineRange(x)
		fromY, toY, rngY := DetermineRange(y)

		if rngX == 0 || rngY == 0 {
			return 0
		}

		bucketRange = GetIntersect(&fromX, &toX, &rngX, fromY, toY, rngY)
		xBuckets := FillBuckets(x, bucketRange)
		yBuckets := FillBuckets(y, bucketRange)
		n := len(xBuckets)
		var high []float64
		var low []float64

		for i, v := range xBuckets {
			if v > yBuckets[i] {
				high = append(high, v)
				low = append(low, yBuckets[i])
			} else {
				high = append(high, yBuckets[i])
				low = append(low, v)
			}
		}

		highTotal := 0
		lowTotal := 0

		for i := 0; i < n; i++ {
			days := dayNum(bucketRange[i].To) - dayNum(bucketRange[i].From)
			highTotal = highTotal + days*int(high[i])
			lowTotal = lowTotal + days*int(low[i])
		}

		cf = 1 / (float64(highTotal) / float64(lowTotal))

	} else {

		return 0
	}

	return cf
}

func GetRandomNames(m map[string]string, c coef) bool {
	allNames := true

	m["table2"] = RandomTableName()
	guid2 := NameToGuid(m["table2"])
	columnNames2 := FetchTableCols(guid2)
	m["valCol2"] = RandomAmountColumn(columnNames2)
	m["dateCol2"] = RandomDateColumn(columnNames2)

	if m["table1"] == "" || m["table2"] == "" || m["valCol1"] == "" || m["valCol2"] == "" || m["dateCol1"] == "" || m["dateCol2"] == "" {
		allNames = false
	}

	if c == S {
		m["table3"] = RandomTableName()
		guid3 := NameToGuid(m["table3"])
		columnNames3 := FetchTableCols(guid3)
		m["valCol3"] = RandomAmountColumn(columnNames3)
		m["dateCol3"] = RandomDateColumn(columnNames3)
		if m["table3"] == "" || m["valCol3"] == "" || m["dateCol3"] == "" {
			allNames = false
		}
	}

	return allNames
}

func GetIntersect(pFromX *time.Time, pToX *time.Time, pRngX *int, fromY time.Time, toY time.Time, rngY int) []FromTo {
	var bucketRange []FromTo
	fromX, toX, rngX := *pFromX, *pToX, *pRngX

	if rngX <= rngY && (fromX == fromY && toX == toY || fromX.After(fromY) && toX.Before(toY)) {
		bucketRange = CreateBuckets(fromX, toX, rngX)
	} else if rngY < rngX && fromY.After(fromX) && toY.Before(toX) {
		bucketRange = CreateBuckets(fromY, toY, rngY)
		*pFromX = fromY
		*pToX = toY
		*pRngX = rngY
	} else if fromX.Before(fromY) && toX.Before(fromY) || fromX.After(toY) && toX.After(toY) {
		return nil
	} else if fromX.Before(fromY) {
		rngYX := dayNum(toX) - dayNum(fromY)
		bucketRange = CreateBuckets(fromY, toX, rngYX)
		*pFromX = fromY
		*pRngX = rngYX
	} else {
		rngXY := dayNum(toY) - dayNum(fromX)
		bucketRange = CreateBuckets(fromX, toY, rngXY)
		*pToX = toY
		*pRngX = rngXY
	}

	return bucketRange
}

func RandomAmountColumn(cols []ColType) string {
	if cols == nil {
		return ""
	}

	rand.Seed(time.Now().UTC().UnixNano())
	columns := make([]string, 0)

	for i, _ := range cols {
		if (cols[i].Sqltype == "numeric" || cols[i].Sqltype == "float" || cols[i].Sqltype == "integer") && cols[i].Name != "transaction_number" {
			columns = append(columns, cols[i].Name)
		}
	}

	n := len(columns)

	if n > 0 {
		x := rand.Intn(n)
		return columns[x]
	} else {
		return ""
	}
}

func RandomDateColumn(cols []ColType) string {
	if cols == nil {
		return ""
	}

	rand.Seed(time.Now().UTC().UnixNano())
	columns := make([]string, 0)

	for _, v := range cols {
		isDate, _ := regexp.MatchString("date", strings.ToLower(v.Name)) //find a column of date type
		if isDate {
			columns = append(columns, v.Name)
		}
	}

	n := len(columns)

	if n > 0 {
		x := rand.Intn(n)
		return columns[x]
	} else {
		return ""
	}
}

func RandomTableName() string {
	var name []string
	err := DB.Table("priv_onlinedata").Order("random()").Limit(1).Pluck("tablename", &name).Error
	if err != nil && err != gorm.RecordNotFound {
		return ""
	}
	return name[0]
}

func NameToGuid(tablename string) string {
	var guid []string
	err := DB.Table("priv_onlinedata").Where("tablename = ?", tablename).Pluck("guid", &guid).Error
	if err != nil && err != gorm.RecordNotFound {
		return ""
	}
	return guid[0]
}

func ExtractDateVal(tablename string, dateCol string, valCol string) []DateVal {
	if tablename == "" || dateCol == "" || valCol == "" {
		return nil
	}

	var dates []time.Time
	var amounts []float64

	d := "DELETE FROM " + tablename + " WHERE " + dateCol + " = '0001-01-01 BC'" ////////TEMP FIX TO GET RID OF INVALID VALUES IN GOV DATA
	DB.Exec(d)

	err = DB.Table(tablename).Pluck(dateCol, &dates).Error
	if err != nil && err != gorm.RecordNotFound {
		check(err)
	}

	err = DB.Table(tablename).Pluck(valCol, &amounts).Error
	if err != nil && err != gorm.RecordNotFound {
		check(err)
	}

	result := make([]DateVal, len(dates))

	for i, v := range dates {
		result[i].Date = v
	}

	for i, v := range amounts {
		result[i].Value = v
	}

	return result
}

func DetermineRange(Dates []DateVal) (time.Time, time.Time, int) {
	var fromDate time.Time
	var toDate time.Time

	if Dates == nil {
		return toDate, fromDate, 0
	}

	dVal, high, low := 0, 0, 100000000

	for _, v := range Dates {
		dVal = dayNum(v.Date)
		if dVal > high {
			high = dVal
			toDate = v.Date
		}
		if dVal < low {
			low = dVal
			fromDate = v.Date
		}
	}

	rng := dayNum(toDate) - dayNum(fromDate)
	return fromDate, toDate, rng
}

func CreateBuckets(fromDate time.Time, toDate time.Time, rng int) []FromTo {
	if rng == 0 {
		return nil
	}

	lim, max := 10, 0

	if rng >= lim { /// no more than 10 buckets
		max = lim
	} else {
		max = rng
	}

	days, bNum := Steps(rng, max) // get days between dates and number of buckets
	extra := 0
	if days > 1 {
		extra = 1
	}
	result := make([]FromTo, bNum+extra)
	date := fromDate // set starting date

	for i := 0; i < bNum; i++ {
		result[i].From = date                   // current date becomes from date
		result[i].To = date.AddDate(0, 0, days) // step amount to to date
		date = result[i].To
	}
	if extra == 1 {
		result[bNum].From = date
		result[bNum].To = toDate.AddDate(0, 0, 1)
	}
	return result
}

func FillBuckets(dateVal []DateVal, bucketRange []FromTo) []float64 {
	if dateVal == nil || bucketRange == nil {
		return nil
	}

	buckets := make([]float64, len(bucketRange))

	for _, v := range dateVal {
		for j, w := range bucketRange {
			if v.Between(w.From, w.To) {
				buckets[j] += float64(v.Value)
				break
			}
		}
	}
	return buckets
}

func (d DateVal) Between(from time.Time, to time.Time) bool {
	if d.Date == from || (d.Date.After(from) && d.Date.Before(to)) {
		return true
	}
	return false
}

func dayNum(d time.Time) int {
	var date time.Time
	var days int

	for i := 1900; i < d.Year(); i++ {
		date = time.Date(i, 12, 31, 0, 0, 0, 0, time.UTC)
		days += date.YearDay()
	}

	days += d.YearDay()
	return days
}

func Steps(a int, b int) (int, int) {
	stepNum := math.Ceil(float64(a) / float64(b))
	bucketNum := a / int(stepNum)
	return int(stepNum), bucketNum
}

func LabelGen(dv []FromTo) []string {
	result := make([]string, len(dv))

	for i, v := range dv {
		if v.From != v.To {
			result[i] = strconv.Itoa(v.From.Day()) + " " + strings.ToUpper(v.From.Month().String()[0:3]) + " " + strconv.Itoa(v.From.Year()) + " to " + strconv.Itoa(v.To.Day()) + " " + strings.ToUpper(v.To.Month().String()[0:3]) + " " + strconv.Itoa(v.To.Year())
		} else {
			result[i] = strconv.Itoa(v.From.Day()) + " " + strings.ToUpper(v.From.Month().String()[0:3]) + " " + strconv.Itoa(v.From.Year())
		}
	}

	return result
}

func daysInMonth(m time.Month, year int) int {
	return time.Date(year, m+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func daysInYear(y int) int {
	d1 := time.Date(y, 1, 1, 0, 0, 0, 0, time.UTC)
	d2 := time.Date(y+1, 1, 1, 0, 0, 0, 0, time.UTC)
	return int(d2.Sub(d1) / (24 * time.Hour))
}
