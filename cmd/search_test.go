package cmd

import (
	"bufio"
	"fmt"
	"github.com/matryer/is"
	"github.com/nleeper/goment"
	"log"
	"os"
	"testing"
)

func CreateCSVFile(fname, data string) {
	file, err := os.Create(fname)
	if err != nil {
		log.Fatalf("Unable to open file for writing: '%s'; %s\n", fname, err)
	}
	w := bufio.NewWriter(file)
	_, err = w.WriteString(data + "\n")
	if err != nil {
		log.Fatalf("Unable to write CSV data to file: '%s'; %s\n", fname, err)
	}
	w.Flush()
	file.Close()
}

func debug(missingDates []goment.Goment, csvStyleDate string) {
	fmt.Println()
	fmt.Println("missingDates")
	fmt.Println("============")
	for _, m := range missingDates {
		fmt.Println(m.ToTime())
	}
	fmt.Println("csvStyleDate:", csvStyleDate)
}

// TestSearch1 - two missing dates in between the first and last dates
func TestSearch1(t *testing.T) {
	fname := "TestSearch1.csv"
	data := "Date,Errors,Warnings\n2021-04-15 06:55:01,0,23\n2021-04-15 08:30:26,0,23\n2021-04-16 06:55:01,0,23\n2021-04-19 06:55:01,0,23"
	CreateCSVFile(fname, data)
	return

	allRootOptions.Amount = 1442
	allRootOptions.Unit = "minutes"
	allRootOptions.Column = 0
	allRootOptions.CsvDelimiter = ","
	allRootOptions.SkipWeekends = false
	allRootOptions.HasNoHeader = false

	missingDates, csvStyleDate := SearchOneFile(fname)
	debug(missingDates, csvStyleDate)
	iss := is.New(t)
	iss.Equal(len(missingDates), 2)
	iss.Equal(missingDates[0].ToTime().String(), "2021-04-17 06:59:01 +0000 UTC")
	iss.Equal(missingDates[1].ToTime().String(), "2021-04-18 07:01:01 +0000 UTC")
	iss.Equal(csvStyleDate, "2021-04-15 06:55:01")
}

// TestSearch2 - second column contains the date field, contains 4 missing dates
func TestSearch2(t *testing.T) {
	fname := "TestSearch2.csv"
	data := "Processed,Date\n5125,2021-04-12\n5197,2021-04-13\n5206,2021-04-14\n5222,2021-04-19\n"
	CreateCSVFile(fname, data)
	return

	allRootOptions.Amount = 1
	allRootOptions.Unit = "days"
	allRootOptions.Column = 1

	missingDates, csvStyleDate := SearchOneFile(fname)
	iss := is.New(t)
	iss.Equal(len(missingDates), 4)
	iss.Equal(missingDates[0].ToTime().String(), "2021-04-15 03:00:00 +0000 UTC")
	iss.Equal(missingDates[1].ToTime().String(), "2021-04-16 04:00:00 +0000 UTC")
	iss.Equal(missingDates[2].ToTime().String(), "2021-04-17 05:00:00 +0000 UTC")
	iss.Equal(missingDates[3].ToTime().String(), "2021-04-18 06:00:00 +0000 UTC")
	iss.Equal(csvStyleDate, "2021-04-12")
}

// TestSearch3 - 5 missing dates, 3 missing data when skipping weekends
func TestSearch3(t *testing.T) {
	fname := "TestSearch3.csv"
	data := "Date,Total\n2021-03-10 18:40:01,317\n2021-03-11 18:40:01,249\n2021-03-15 18:40:04,287\n2021-03-16 18:40:03,320\n2021-03-19 18:40:06,102\n"
	CreateCSVFile(fname, data)
	return

	allRootOptions.Amount = 25
	allRootOptions.Unit = "hours"
	allRootOptions.Column = 0

	missingDates, csvStyleDate := SearchOneFile(fname)
	//debug(missingDates, csvStyleDate)

	iss := is.New(t)
	iss.Equal(len(missingDates), 5)
	iss.Equal(missingDates[0].ToTime().String(), "2021-03-12 20:40:01 +0000 UTC")
	iss.Equal(missingDates[1].ToTime().String(), "2021-03-13 21:40:01 +0000 UTC")
	iss.Equal(missingDates[2].ToTime().String(), "2021-03-14 22:40:01 +0000 UTC")
	iss.Equal(missingDates[3].ToTime().String(), "2021-03-18 01:40:01 +0000 UTC")
	iss.Equal(missingDates[4].ToTime().String(), "2021-03-19 02:40:01 +0000 UTC")
	iss.Equal(csvStyleDate, "2021-03-10 18:40:01")

	allRootOptions.SkipWeekends = true
	missingDates, _ = SearchOneFile(fname)
	iss.Equal(len(missingDates), 3)
	iss.Equal(missingDates[0].ToTime().String(), "2021-03-12 20:40:01 +0000 UTC")
	iss.Equal(missingDates[1].ToTime().String(), "2021-03-18 01:40:01 +0000 UTC")
	iss.Equal(missingDates[2].ToTime().String(), "2021-03-19 02:40:01 +0000 UTC")
}

// TestSearch4 - 3 missing dates, tab-delimited file
func TestSearch4(t *testing.T) {
	fname := "TestSearch4.csv"
	data := "Date\tCount\n2021-04-01 18:40:00\t318\n2021-04-02 18:40:00\t252\n2021-04-06 18:40:00\t291\n2021-04-07 18:40:00\t274\n2021-04-08 18:40:01\t243"
	CreateCSVFile(fname, data)
	return

	allRootOptions.Amount = 25
	allRootOptions.Unit = "hours"
	allRootOptions.Column = 0
	allRootOptions.CsvDelimiter = "\t"
	allRootOptions.SkipWeekends = false

	missingDates, csvStyleDate := SearchOneFile(fname)

	iss := is.New(t)
	iss.Equal(len(missingDates), 3)
	iss.Equal(missingDates[0].ToTime().String(), "2021-04-03 20:40:00 +0000 UTC")
	iss.Equal(missingDates[1].ToTime().String(), "2021-04-04 21:40:00 +0000 UTC")
	iss.Equal(missingDates[2].ToTime().String(), "2021-04-05 22:40:00 +0000 UTC")
	iss.Equal(csvStyleDate, "2021-04-01 18:40:00")
}

// TestSearch5 - No header, 1 missing date
func TestSearch5(t *testing.T) {
	fname := "TestSearch5.csv"
	data := "2021-04-01 18:40:00,318\n2021-04-02 18:40:00,252\n2021-04-06 18:40:00,291\n2021-04-07 18:40:00,274\n2021-04-08 18:40:01,243"
	CreateCSVFile(fname, data)
	return

	allRootOptions.Amount = 25
	allRootOptions.Unit = "hours"
	allRootOptions.Column = 0
	allRootOptions.CsvDelimiter = ","
	allRootOptions.SkipWeekends = true
	allRootOptions.HasNoHeader = true

	missingDates, csvStyleDate := SearchOneFile(fname)

	iss := is.New(t)
	iss.Equal(len(missingDates), 1)
	iss.Equal(missingDates[0].ToTime().String(), "2021-04-05 22:40:00 +0000 UTC")
	iss.Equal(csvStyleDate, "2021-04-01 18:40:00")
}

// TestSearch6 - skipDays (Wednesday,Thursday), 1 missing date
func TestSearch6(t *testing.T) {
	fname := "TestSearch6.csv"
	// this has embedded tab characters
	data := "Date	Errors	Warnings\n2021-04-13 06:55:01	0	23\n2021-04-16 06:55:01	0	23\n2021-04-17 06:55:01	0	19\n2021-04-18 06:55:01	0	18\n2021-04-19 06:55:01	0	23\n2021-04-22 06:55:01	0	17\n"
	CreateCSVFile(fname, data)
	return

	allRootOptions.Amount = 86400
	allRootOptions.Unit = "seconds"
	allRootOptions.Column = 0
	allRootOptions.CsvDelimiter = "\t"
	allRootOptions.SkipWeekends = false
	allRootOptions.HasNoHeader = false
	allRootOptions.SkipDays = "Wednesday,Thursday"

	missingDates, csvStyleDate := SearchOneFile(fname)
	iss := is.New(t)
	iss.Equal(len(missingDates), 1)
	iss.Equal(missingDates[0].ToTime().String(), "2021-04-20 06:55:01 +0000 UTC")
	iss.Equal(csvStyleDate, "2021-04-13 06:55:01")
}

// TestSearch7 - no date gaps
func TestSearch7(t *testing.T) {
	fname := "TestSearch7.csv"
	data := "2021-04-15 06:55:01,0,23\n2021-04-15 08:30:26,0,23\n2021-04-16 06:55:01,0,23\n2021-04-19 06:55:01,0,23\n"
	CreateCSVFile(fname, data)
	return

	allRootOptions.Amount = 1
	allRootOptions.Unit = "days"
	allRootOptions.Column = 0
	allRootOptions.CsvDelimiter = ","
	allRootOptions.SkipWeekends = true
	allRootOptions.HasNoHeader = true

	missingDates, csvStyleDate := SearchOneFile(fname)
	iss := is.New(t)
	iss.Equal(len(missingDates), 0)
	iss.Equal(csvStyleDate,"2021-04-15 06:55:01")
}
