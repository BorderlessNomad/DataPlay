package main

import (
	"fmt"
)

type YN struct {
	year int
	num  int
}

type TY struct {
	table string
	ynums []YN
}

type CT struct {
	col   string
	table string
}

var err error

func DateScrapeA() []CT {
	var ct = []CT{}

	row, _ := DB.Raw("SELECT COLUMN_NAME, TABLE_NAME FROM INFORMATION_SCHEMA.COLUMNS WHERE COLUMN_NAME LIKE '%date%' AND DATA_TYPE = 'date'").Rows()

	for row.Next() {
		var COLUMN_NAME string
		var TABLE_NAME string
		err = row.Scan(&COLUMN_NAME, &TABLE_NAME)
		var tmp CT
		tmp.col = COLUMN_NAME
		tmp.table = TABLE_NAME
		ct = append(ct, tmp)
	}
	return ct
}

func DateScrapeB(ct []CT) {
	var ty = []TY{}

	for _, v := range ct {
		var yn = []YN{}
		if v.table != "imp0f69b88c01ac0bf7b21d2fb62857991974118d7f_f798e15cc6de53e4cc3" && v.table != "impb23bc6857c05bbc137b0538fc84aaab9cd7bad85_63bf6d694c646f2d205" {
			row, _ := DB.Raw(fmt.Sprintf("SELECT EXTRACT(YEAR FROM %s), COUNT(*) AS NUM FROM %s GROUP BY EXTRACT(YEAR FROM %s)", v.col, v.table, v.col)).Rows()

			for row.Next() {
				var DATE_PART int
				var NUM int
				err = row.Scan(&DATE_PART, &NUM)
				var tmp YN
				tmp.year = DATE_PART
				tmp.num = NUM
				yn = append(yn, tmp)
			}
			var tmp2 TY
			tmp2.table = v.table
			tmp2.ynums = yn
			ty = append(ty, tmp2)
		}
	}

	var yArr [67]int

	for _, v := range ty {
		for _, w := range v.ynums {
			year := w.year
			cur := w.num

			switch year {
			case 1950:
				yArr[0] += cur
			case 1951:
				yArr[1] += cur
			case 1952:
				yArr[2] += cur
			case 1953:
				yArr[3] += cur
			case 1954:
				yArr[4] += cur
			case 1955:
				yArr[5] += cur
			case 1956:
				yArr[6] += cur
			case 1957:
				yArr[7] += cur
			case 1958:
				yArr[8] += cur
			case 1959:
				yArr[9] += cur
			case 1960:
				yArr[10] += cur
			case 1961:
				yArr[11] += cur
			case 1962:
				yArr[12] += cur
			case 1963:
				yArr[13] += cur
			case 1964:
				yArr[14] += cur
			case 1965:
				yArr[15] += cur
			case 1966:
				yArr[16] += cur
			case 1967:
				yArr[17] += cur
			case 1968:
				yArr[18] += cur
			case 1969:
				yArr[19] += cur
			case 1970:
				yArr[20] += cur
			case 1971:
				yArr[21] += cur
			case 1972:
				yArr[22] += cur
			case 1973:
				yArr[23] += cur
			case 1974:
				yArr[24] += cur
			case 1975:
				yArr[25] += cur
			case 1976:
				yArr[26] += cur
			case 1977:
				yArr[27] += cur
			case 1978:
				yArr[28] += cur
			case 1979:
				yArr[29] += cur
			case 1980:
				yArr[30] += cur
			case 1981:
				yArr[31] += cur
			case 1982:
				yArr[32] += cur
			case 1983:
				yArr[33] += cur
			case 1984:
				yArr[34] += cur
			case 1985:
				yArr[35] += cur
			case 1986:
				yArr[36] += cur
			case 1987:
				yArr[37] += cur
			case 1988:
				yArr[38] += cur
			case 1989:
				yArr[39] += cur
			case 1990:
				yArr[40] += cur
			case 1991:
				yArr[41] += cur
			case 1992:
				yArr[42] += cur
			case 1993:
				yArr[43] += cur
			case 1994:
				yArr[44] += cur
			case 1995:
				yArr[45] += cur
			case 1996:
				yArr[46] += cur
			case 1997:
				yArr[47] += cur
			case 1998:
				yArr[48] += cur
			case 1999:
				yArr[49] += cur
			case 2000:
				yArr[50] += cur
			case 2001:
				yArr[51] += cur
			case 2002:
				yArr[52] += cur
			case 2003:
				yArr[53] += cur
			case 2004:
				yArr[54] += cur
			case 2005:
				yArr[55] += cur
			case 2006:
				yArr[56] += cur
			case 2007:
				yArr[57] += cur
			case 2008:
				yArr[58] += cur
			case 2009:
				yArr[59] += cur
			case 2010:
				yArr[60] += cur
			case 2011:
				yArr[61] += cur
			case 2012:
				yArr[62] += cur
			case 2013:
				yArr[63] += cur
			case 2014:
				yArr[64] += cur
			case 2015:
				yArr[65] += cur
			default:
				yArr[66] += cur
			}
		}
	}
}
