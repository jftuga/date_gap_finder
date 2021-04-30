package cmd

import (
	"fmt"
	"github.com/matryer/is"
	"testing"
)

// TODO: test -O switch and ensure backup files exist

func debugArr(allRows []string) {
	fmt.Println()
	fmt.Println("csv records")
	fmt.Println("===========")
	for _, row := range allRows {
		fmt.Println(row)
	}
}

// TestInsert1 - two missing dates, insert -1,99 for missing ros
func TestInsert1(t *testing.T) {
	fname := "TestInsert1.csv"
	data := "Date,Errors,Warnings\n2021-04-15 06:55:01,0,23\n2021-04-15 08:30:26,0,23\n2021-04-16 06:55:01,0,23\n2021-04-19 06:55:01,0,23"
	CreateCSVFile(fname, data)

	allRootOptions.Amount = 1442
	allRootOptions.Period = "minutes"
	allRootOptions.Column = 0

	allInsertOptions.columnInserts = []string {"1,-1", "2,999"}

	csv := insertOneFile(fname)
	iss := is.New(t)
	iss.Equal(len(csv), 7)
	iss.Equal(csv[0], "Date,Errors,Warnings")
	iss.Equal(csv[1], "2021-04-15 06:55:01,0,23")
	iss.Equal(csv[2], "2021-04-15 08:30:26,0,23")
	iss.Equal(csv[3], "2021-04-16 06:55:01,0,23")
	iss.Equal(csv[4], "2021-04-17 06:59:01,-1,999")
	iss.Equal(csv[5], "2021-04-18 07:01:01,-1,999")
	iss.Equal(csv[6], "2021-04-19 06:55:01,0,23")
}

// TestInsert2 - second column contains the date field, insert 4 missing rows all with 999888
func TestInsert2(t *testing.T) {
	fname := "TestInsert2.csv"
	data := "Processed,Date\n5125,2021-04-12\n5197,2021-04-13\n5206,2021-04-14\n5222,2021-04-19\n"
	CreateCSVFile(fname, data)

	allRootOptions.Amount = 1
	allRootOptions.Period = "days"
	allRootOptions.Column = 1

	allInsertOptions.columnInserts = []string {}
	allInsertOptions.allColumnInserts = "999888"

	csv := insertOneFile(fname)
	iss := is.New(t)
	iss.Equal(len(csv), 9)
	iss.Equal(csv[0], "Processed,Date")
	iss.Equal(csv[1], "5125,2021-04-12")
	iss.Equal(csv[2], "5197,2021-04-13")
	iss.Equal(csv[3], "5206,2021-04-14")
	iss.Equal(csv[4], "999888,2021-04-15")
	iss.Equal(csv[5], "999888,2021-04-16")
	iss.Equal(csv[6], "999888,2021-04-17")
	iss.Equal(csv[7], "999888,2021-04-18")
	iss.Equal(csv[8], "5222,2021-04-19")
}

// TestInsert3 - 3 missing dates, tab-delimited file, insert col 2 with -1, all other columns with 999
func TestInsert3(t *testing.T) {
	fname := "TestInsert3.csv"
	data := "Date,Errors,Warnings,N1,N2\n2021-04-15 06:55:01,0,23,15,62\n2021-04-15 08:30:26,0,23,15,62\n2021-04-16 06:55:01,0,23,15,62\n2021-04-19 06:55:01,0,23,15,62"
	CreateCSVFile(fname, data)

	allRootOptions.Amount = 1442
	allRootOptions.Period = "minutes"
	allRootOptions.Column = 0
	allRootOptions.CsvDelimiter = ","
	allRootOptions.SkipWeekends = false

	allInsertOptions.columnInserts = []string{"2,-1"}
	allInsertOptions.allColumnInserts = "999"

	csv := insertOneFile(fname)
	iss := is.New(t)
	iss.Equal(len(csv), 7)
	iss.Equal(csv[0], "Date,Errors,Warnings,N1,N2")
	iss.Equal(csv[1], "2021-04-15 06:55:01,0,23,15,62")
	iss.Equal(csv[2], "2021-04-15 08:30:26,0,23,15,62")
	iss.Equal(csv[3], "2021-04-16 06:55:01,0,23,15,62")
	iss.Equal(csv[4], "2021-04-17 06:59:01,999,-1,999,999")
	iss.Equal(csv[5], "2021-04-18 07:01:01,999,-1,999,999")
	iss.Equal(csv[6], "2021-04-19 06:55:01,0,23,15,62")
}
