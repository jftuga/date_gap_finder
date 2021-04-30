package cmd

import (
	"bufio"
	"fmt"
	"github.com/matryer/is"
	"log"
	"os"
	"testing"
)

func TestSearch1(t *testing.T) {
	fname := "test1.csv"
	data := "Date,Errors,Warnings\n2021-04-15 06:55:01,0,23\n2021-04-15 08:30:26,0,23\n2021-04-16 06:55:01,0,23\n2021-04-19 06:55:01,0,23"

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
	allRootOptions.Amount = 1442
	allRootOptions.Period = "minutes"
	fmt.Println(allRootOptions)

	missingDates, csvStyleDate := SearchOneFile(fname)
	iss := is.New(t)
	iss.Equal(csvStyleDate, "2021-04-15 06:55:01)")
	for _, m := range missingDates {
		fmt.Println(m.ToTime())
	}


}
