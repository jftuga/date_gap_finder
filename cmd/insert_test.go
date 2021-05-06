package cmd

import (
	"fmt"
	"github.com/jftuga/date_gap_finder/fileOps"
	"github.com/matryer/is"
	"log"
	"os"
	"testing"
)

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

	allRootOptions.Amount = 1440
	allRootOptions.Unit = "minutes"
	allRootOptions.Column = 0
	allRootOptions.Padding = "1s"

	allInsertOptions.columnInserts = []string {"1,-1", "2,999"}

	csv := insertOneFile(fname)
	iss := is.New(t)
	iss.Equal(len(csv), 7)
	iss.Equal(csv[0], "Date,Errors,Warnings")
	iss.Equal(csv[1], "2021-04-15 06:55:01,0,23")
	iss.Equal(csv[2], "2021-04-15 08:30:26,0,23")
	iss.Equal(csv[3], "2021-04-16 06:55:01,0,23")
	iss.Equal(csv[4], "2021-04-17 06:55:01,-1,999")
	iss.Equal(csv[5], "2021-04-18 06:55:01,-1,999")
	iss.Equal(csv[6], "2021-04-19 06:55:01,0,23")
}

// TestInsert2 - second column contains the date field, insert 4 missing rows all with 999888
func TestInsert2(t *testing.T) {
	fname := "TestInsert2.csv"
	data := "Processed,Date\n5125,2021-04-12\n5197,2021-04-13\n5206,2021-04-14\n5222,2021-04-19\n"
	CreateCSVFile(fname, data)

	allRootOptions.Amount = 1
	allRootOptions.Unit = "days"
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
	allRootOptions.Unit = "minutes"
	allRootOptions.Column = 0
	allRootOptions.CsvDelimiter = ","
	allRootOptions.SkipWeekends = false
	allRootOptions.Padding = "3m"

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


// TestInsert4 - 1 missing date, overwrite file
func TestInsert4(t *testing.T) {
	fname := "TestInsert4.csv"
	data := "Date,Amount\n2021-04-01 18:40:00,318\n2021-04-02 18:40:00,252\n2021-04-06 18:40:00,291\n2021-04-07 18:40:00,274\n2021-04-08 18:40:01,243"
	CreateCSVFile(fname, data)

	var err error
	var origCsvStat, bakCsvStat os.FileInfo
	origCsvStat, err = os.Stat(fname)
	if err != nil {
		log.Fatalln(err)
	}

	allRootOptions.Amount = 24
	allRootOptions.Unit = "hours"
	allRootOptions.Column = 0
	allRootOptions.CsvDelimiter = ","
	allRootOptions.SkipWeekends = true
	allRootOptions.HasNoHeader = false
	allRootOptions.Padding = "2m"

	allInsertOptions.columnInserts = []string{"1,-1"}
	allInsertOptions.allColumnInserts = ""
	allInsertOptions.Overwrite = true

	csv := insertOneFile(fname)
	iss := is.New(t)
	iss.Equal(len(csv), 7)
	iss.Equal(csv[0], "Date,Amount")
	iss.Equal(csv[1], "2021-04-01 18:40:00,318")
	iss.Equal(csv[2], "2021-04-02 18:40:00,252")
	iss.Equal(csv[3], "2021-04-05 18:40:00,-1")
	iss.Equal(csv[4], "2021-04-06 18:40:00,291")
	iss.Equal(csv[5], "2021-04-07 18:40:00,274")
	iss.Equal(csv[6], "2021-04-08 18:40:01,243")

	if ! allInsertOptions.Overwrite {
		return
	}

	result, _ := fileOps.OverwriteCsv(fname, csv)
	iss.True(result)

	bakCsvStat, err = os.Stat(fname)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println("o:", origCsvStat.Size(), "  b:", bakCsvStat.Size())
	iss.True(origCsvStat.Size() < bakCsvStat.Size() )
}

// TestInsert5 - same a 4 but with no header, 1 missing date, overwrite file
func TestInsert5(t *testing.T) {
	fname := "TestInsert5.csv"
	data := "2021-04-01 18:40:00,318\n2021-04-02 18:40:00,252\n2021-04-06 18:40:00,291\n2021-04-07 18:40:00,274\n2021-04-08 18:40:01,243"
	CreateCSVFile(fname, data)

	var err error
	var origCsvStat, bakCsvStat os.FileInfo
	origCsvStat, err = os.Stat(fname)
	if err != nil {
		log.Fatalln(err)
	}

	allRootOptions.Amount = 24
	allRootOptions.Unit = "hours"
	allRootOptions.Column = 0
	allRootOptions.CsvDelimiter = ","
	allRootOptions.SkipWeekends = true
	allRootOptions.HasNoHeader = true
	allInsertOptions.columnInserts = []string{"1,-1"}
	allInsertOptions.allColumnInserts = ""
	allInsertOptions.Overwrite = true
	allRootOptions.Padding = "2s"

	csv := insertOneFile(fname)
	iss := is.New(t)
	iss.Equal(len(csv), 6)
	iss.Equal(csv[0], "2021-04-01 18:40:00,318")
	iss.Equal(csv[1], "2021-04-02 18:40:00,252")
	iss.Equal(csv[2], "2021-04-05 18:40:00,-1")
	iss.Equal(csv[3], "2021-04-06 18:40:00,291")
	iss.Equal(csv[4], "2021-04-07 18:40:00,274")
	iss.Equal(csv[5], "2021-04-08 18:40:01,243")

	if ! allInsertOptions.Overwrite {
		return
	}

	result, _ := fileOps.OverwriteCsv(fname, csv)
	iss.True(result)

	bakCsvStat, err = os.Stat(fname)
	if err != nil {
		log.Fatalln(err)
	}
	//fmt.Println("o:", origCsvStat.Size(), "  b:", bakCsvStat.Size())
	iss.True(origCsvStat.Size() < bakCsvStat.Size() )
}
