package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"math"
	"math/rand"
	"net/http"
	"time"
)

type cmeth int

const ( //go version of enum
	P cmeth = iota
	S
	V
)

type CorrelationData struct {
	CorrelationId int       `json:"correlationid, omitempty"`
	Method        string    `json:"method"`
	ChartType     string    `json:"type, omitempty"`
	Discovered    bool      `json:"discovered,omitempty"`
	From          string    `json:"from"`
	To            string    `json:"to"`
	Table1        TableData `json:"table1"`
	Table2        TableData `json:"table2"`
	Table3        TableData `json:"table3, omitempty"`
	Coefficient   float64   `json:"coefficient, omitempty"`
}

type TableData struct {
	RelationId string  `json:"relationid, omitempty"`
	ChartType  string  `json:"type, omitempty"`
	Discovered bool    `json:"discovered,omitempty"`
	Title      string  `json:"title"`
	Desc       string  `json:"desc"`
	LabelX     string  `json:"xLabel"`
	LabelXLong string  `json:"xxLabel"`
	LabelY     string  `json:"yLabel,omitempty"`
	LabelYLong string  `json:"yyLabel,omitempty"`
	LabelZ     string  `json:"zLabel,omitempty"`
	LabelZLong string  `json:"zzLabel,omitempty"`
	Values     []XYVal `json:"values,omitempty"`
}

type XYVal struct {
	X     string `json:"x"`
	Y     string `json:"y"`
	Z     string `json:"z,omitempty"`
	Xtype string `json:"-"`
	Ytype string `json:"-"`
	Ztype string `json:"-"`
}

type DateVal struct {
	Date  time.Time
	Value float64
}

type FromTo struct {
	From time.Time
	To   time.Time
}

type TableCols struct {
	tbl1  string
	tbl2  string
	tbl3  string
	dat1  string
	dat2  string
	dat3  string
	val1  string
	val2  string
	val3  string
	ctype string
	chart string
}

/**
 * @details Take in table name and a threshold for the looping and attempt to get correlated tables
 * through a variety of correlation methods.
 * Use 3:1:1 for Spurious to Pearson and Visual as Spurious less likely to find correlations
 *
 * @param string tablename
 * @param int searchDepth
 */
func GenerateCorrelations(tablename string, searchDepth int) {
	rand.Seed(time.Now().UTC().UnixNano())
	cols1 := GetSQLTableSchema(tablename)
	var tablenames []string

	err := DB.Model(OnlineData{}).Order("random()").Pluck("tablename", &tablenames).Error
	if err != nil && err != gorm.RecordNotFound {
		fmt.Println("Database Error while generating correlations.", err)
		return
	} else if err == gorm.RecordNotFound {
		fmt.Println("No tables for generating correlations.", err)
		return
	}

	tableCols := make([]TableCols, searchDepth)

	for i := 0; i < searchDepth; i++ {
		attemptCorrelation := true

		tableCols[i].tbl1 = tablename
		tableCols[i].tbl2 = tablenames[rand.Intn(len(tablenames))]

		cols2 := GetSQLTableSchema(tableCols[i].tbl2)

		tableCols[i].dat1 = RandomDateColumn(cols1)
		tableCols[i].val1 = RandomValueColumn(cols1)

		tableCols[i].dat2 = RandomDateColumn(cols2)
		tableCols[i].val2 = RandomValueColumn(cols2)

		if tableCols[i].tbl1 == tableCols[i].tbl2 || tableCols[i].tbl1 == "" || tableCols[i].tbl2 == "" || tableCols[i].val1 == "" || tableCols[i].val2 == "" || tableCols[i].dat1 == "" || tableCols[i].dat2 == "" {
			attemptCorrelation = false
		}

		tableCols[i].ctype = RandomMethod()

		if tableCols[i].ctype == "Spurious" {
			tableCols[i].tbl3 = tablenames[rand.Intn(len(tablenames))]

			cols3 := GetSQLTableSchema(tableCols[i].tbl3)

			tableCols[i].dat3 = RandomDateColumn(cols3)
			tableCols[i].val3 = RandomValueColumn(cols3)
			if tableCols[i].tbl1 == tableCols[i].tbl3 || tableCols[i].tbl2 == tableCols[i].tbl3 || tableCols[i].tbl3 == "" || tableCols[i].val3 == "" || tableCols[i].dat3 == "" {
				attemptCorrelation = false
			}
		}

		if attemptCorrelation {
			go AttemptCorrelation(tableCols[i])
		}
	}
}

/**
 * Take in table name and a correlation type, then get some random apt columns from it and
 * generate more random tables and columns and check for any  pre-existing correlations on
 * that combination.
 *
 * If a correlation for the generated tables combination doesn't exist, attempt to calculate
 * a new correlation coefficient and save the new correlation.
 *
 * @param TableCols tableCols
 * @return *appError
 */
func AttemptCorrelation(tableCols TableCols) *appError {
	c := P
	if tableCols.ctype == "Pearson" {
		c = P
	} else if tableCols.ctype == "Spurious" {
		c = S
	} else {
		c = V
	}

	correlation := Correlation{}
	ctype := ""
	query := DB.Where("tbl1 = ?", tableCols.tbl1).Where("col1 = ?", tableCols.val1).Where("tbl2 = ?", tableCols.tbl2).Where("col2 = ?", tableCols.val2)
	if c == P {
		ctype = "Pearson"
		query = query.Where("method = ?", ctype)
	} else if c == S {
		ctype = "Spurious"
		query = query.Where("tbl3 = ?", tableCols.tbl3).Where("col3 = ?", tableCols.val3).Where("method = ?", ctype)
	} else if c == V {
		ctype = "Visual"
		query = query.Where("method = ?", ctype)
	}

	err := query.Find(&correlation).Error
	if err != nil && err != gorm.RecordNotFound {
		return &appError{err, "Database query failed (DateCol - " + ctype + ").", http.StatusInternalServerError}
	} else if err == gorm.RecordNotFound {
		return &appError{err, "Unable to find correlation (DateCol - " + ctype + ").", http.StatusNotFound}
	}

	if correlation.Json == nil { // if no correlation exists then generate one
		cd := new(CorrelationData)
		cf, errCF := CalculateCoefficient(tableCols, c, cd)
		if errCF != nil {
			return errCF
		}

		if cf != 0 { //Save the various permutations of the correlation if one is generated
			if tableCols.ctype == "Pearson" {
				tableCols.chart = "bar"
				err := SaveCorrelation(tableCols, c, cf, cd) // save everything to the correlation table
				if err != nil {
					return err
				}

				tableCols.chart = "column"
				err = SaveCorrelation(tableCols, c, cf, cd) // save everything to the correlation table
				if err != nil {
					return err
				}

				tableCols.chart = "line"
				err = SaveCorrelation(tableCols, c, cf, cd) // save everything to the correlation table
				if err != nil {
					return err
				}

				tableCols.chart = "scatter"
				err = SaveCorrelation(tableCols, c, cf, cd) // save everything to the correlation table
				if err != nil {
					return err
				}
			} else if tableCols.ctype == "Spurious" {
				tableCols.chart = "line"
				err := SaveCorrelation(tableCols, c, cf, cd) // save everything to the correlation table
				if err != nil {
					return err
				}

				tableCols.chart = "scatter"
				err = SaveCorrelation(tableCols, c, cf, cd) // save everything to the correlation table
				if err != nil {
					return err
				}
			} else if tableCols.ctype == "Visual" {
				tableCols.chart = "bar"
				err := SaveCorrelation(tableCols, c, cf, cd) // save everything to the correlation table
				if err != nil {
					return err
				}

				tableCols.chart = "column"
				err = SaveCorrelation(tableCols, c, cf, cd) // save everything to the correlation table
				if err != nil {
					return err
				}

				tableCols.chart = "line"
				err = SaveCorrelation(tableCols, c, cf, cd) // save everything to the correlation table
				if err != nil {
					return err
				}

				tableCols.chart = "scatter"
				err = SaveCorrelation(tableCols, c, cf, cd) // save everything to the correlation table
				if err != nil {
					return err
				}
			} else {
				tableCols.chart = "unknown"
				err := SaveCorrelation(tableCols, c, cf, cd) // save everything to the correlation table
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Determine if two sets of dates overlap - X values are referenced so they can be altered in place and passed back
// again when used with Spurious correlation which covers the intersect between three data sets
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
		rngYX := DayNum(toX) - DayNum(fromY)
		bucketRange = CreateBuckets(fromY, toX, rngYX)
		*pFromX = fromY
		*pRngX = rngYX
	} else {
		rngXY := DayNum(toY) - DayNum(fromX)
		bucketRange = CreateBuckets(fromX, toY, rngXY)
		*pToX = toY
		*pRngX = rngXY
	}

	return bucketRange
}

// Generate a correlation coefficient (if data allows), based on the requested correlation type
// Uses discrete buckets in order to normalise data for calculating coefficient but stores the entire range of x,y,z values in the correlation data struct
func CalculateCoefficient(tableCols TableCols, c cmeth, cd *CorrelationData) (float64, *appError) {
	var hasVals bool
	var bucketRange []FromTo
	var xBuckets, yBuckets, zBuckets []float64
	var cf float64

	x, errX := ExtractDateVal(tableCols.tbl1, tableCols.dat1, tableCols.val1)
	if errX != nil {
		return 0.0, errX
	}

	y, errY := ExtractDateVal(tableCols.tbl2, tableCols.dat2, tableCols.val2)
	if errY != nil {
		return 0.0, errY
	}

	fromX, toX, rngX := DetermineRange(x)
	fromY, toY, rngY := DetermineRange(y)

	if rngX == 0 || rngY == 0 {
		return 0.0, nil
	}

	bucketRange = GetIntersect(&fromX, &toX, &rngX, fromY, toY, rngY)
	if bucketRange == nil {
		return 0, nil
	}

	xBuckets = FillBuckets(x, bucketRange)
	yBuckets = FillBuckets(y, bucketRange)

	if MostlyEmpty(xBuckets) || MostlyEmpty(yBuckets) {
		return 0.0, nil
	}

	l := len(bucketRange) - 1

	(*cd).Table1.Values, hasVals = GetValues(x, bucketRange[0].From, bucketRange[l].To)
	if !hasVals {
		return 0.0, nil
	}

	(*cd).Table2.Values, hasVals = GetValues(y, bucketRange[0].From, bucketRange[l].To)
	if !hasVals {
		return 0.0, nil
	}

	(*cd).From = (bucketRange[0].From.String()[0:10])
	(*cd).To = (bucketRange[l].To.AddDate(0, 0, -1).String()[0:10])

	if c == P {
		cf = Pearson(xBuckets, yBuckets)
	} else if c == V {
		cf = Visual(xBuckets, yBuckets, bucketRange)
	} else if c == S {
		z, errZ := ExtractDateVal(tableCols.tbl3, tableCols.dat3, tableCols.val3)
		if errZ != nil {
			return 0.0, errZ
		}

		fromZ, toZ, rngZ := DetermineRange(z)

		if rngZ == 0 {
			return 0.0, nil
		}

		//from X, toX and rngX now equal full from, to and rng of x and y from last iteration so just get intersect of those with Z
		bucketRange = GetIntersect(&fromX, &toX, &rngX, fromZ, toZ, rngZ)
		if bucketRange == nil {
			return 0.0, nil
		}

		l := len(bucketRange) - 1
		xBuckets = FillBuckets(x, bucketRange)
		yBuckets = FillBuckets(y, bucketRange)
		zBuckets = FillBuckets(z, bucketRange)

		if MostlyEmpty(xBuckets) || MostlyEmpty(yBuckets) || MostlyEmpty(zBuckets) {
			return 0.0, nil
		}

		(*cd).Table1.Values, hasVals = GetValues(x, bucketRange[0].From, bucketRange[l].To)
		if !hasVals {
			return 0.0, nil
		} else {
			(*cd).Table2.Values, hasVals = GetValues(y, bucketRange[0].From, bucketRange[l].To)
		}

		if !hasVals {
			return 0.0, nil
		}

		(*cd).Table3.Values, hasVals = GetValues(z, bucketRange[0].From, bucketRange[l].To)

		if !hasVals {
			return 0.0, nil
		}

		(*cd).From = (bucketRange[0].From.String()[0:10])
		(*cd).To = (bucketRange[l].To.String()[0:10])

		cf = Spurious(yBuckets, zBuckets, xBuckets) // order is table2 = x arg , table3 = y arg, table1 = z arg so that we get correlation of 2 random tables against underlying table
	} else {
		return 0.0, nil
	}
	return cf, nil
}

//Create a json string containing all the data needed for generating a graph and then insert this and all the other correlation info into the correlations table
func SaveCorrelation(tableCols TableCols, c cmeth, cf float64, cd *CorrelationData) *appError {
	ind1, ind2, ind3 := Index{}, Index{}, Index{}

	guid1, _ := GetGuid(tableCols.tbl1)
	guid2, _ := GetGuid(tableCols.tbl2)
	guid3, _ := GetGuid(tableCols.tbl3)

	err1 := DB.Model(&ind1).Where("guid= ?", guid1).Find(&ind1).Error
	if err1 != nil {
		return &appError{err1, "Database query failed (guid index1)", http.StatusServiceUnavailable}
	}

	err2 := DB.Model(&ind2).Where("guid= ?", guid2).Find(&ind2).Error
	if err2 != nil {
		return &appError{err2, "Database query failed (guid index2)", http.StatusServiceUnavailable}
	}

	if c == S {
		err3 := DB.Model(&ind3).Where("guid= ?", guid3).Find(&ind3).Error
		if err3 != nil {
			return &appError{err3, "Database query failed (guid index3)", http.StatusServiceUnavailable}
		}
	}

	(*cd).Method = tableCols.ctype
	(*cd).ChartType = tableCols.chart
	(*cd).Table1.Title = SanitizeString(ind1.Title)
	(*cd).Table2.Title = SanitizeString(ind2.Title)
	(*cd).Table1.Desc = SanitizeString(ind1.Notes)
	(*cd).Table2.Desc = SanitizeString(ind2.Notes)
	(*cd).Table1.LabelX = tableCols.dat1
	(*cd).Table2.LabelX = tableCols.dat2
	(*cd).Table1.LabelY = tableCols.val1
	(*cd).Table2.LabelY = tableCols.val2

	if c == S {
		(*cd).Table3.Title = SanitizeString(ind3.Title)
		(*cd).Table3.Desc = SanitizeString(ind3.Notes)
		(*cd).Table3.LabelX = tableCols.dat3
		(*cd).Table3.LabelY = tableCols.val3
	}

	jv, _ := json.Marshal(*cd)

	correlation := Correlation{
		Tbl1:    tableCols.tbl1,
		Col1:    tableCols.val1,
		Tbl2:    tableCols.tbl2,
		Col2:    tableCols.val2,
		Tbl3:    tableCols.tbl3,
		Col3:    tableCols.val3,
		Method:  tableCols.ctype,
		Coef:    cf,
		Json:    jv,
		Abscoef: math.Abs(cf), //absolute value for ranking as highly negative correlations are also interesting
	}

	err4 := DB.Save(&correlation).Error
	if err4 != nil {
		return &appError{err4, "Database query failed (save correlation)", http.StatusServiceUnavailable}
	}

	return nil
}
