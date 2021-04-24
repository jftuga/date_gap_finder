/*
Copyright © 2021 John Taylor

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/jftuga/date_gap_finder/fileOps"
	"github.com/nleeper/goment"
	"github.com/spf13/cobra"
	"log"
)

type insertOptions struct {
	columnInserts []string
}

var allInsertOptions insertOptions

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "insert missing CSV entries",
	Long: `CSV dates are assumed to be oldest to newest within the file.
Multiple -r options can be used.  Each -r option is comma-delimited with
the column number first and the value to insert (into that column) second.`,
	Run: func(cmd *cobra.Command, args []string) {
		insertAllFiles(args)
		/*
		fmt.Println("insert called")
		for _, v := range allInsertOptions.columnInserts {
			fmt.Println(v)
		}
		*/
	},
}

func init() {
	rootCmd.AddCommand(insertCmd)
	insertCmd.Flags().StringArrayVarP(&allInsertOptions.columnInserts, "record", "r", []string{}, "insert record with missing date")
}

func insertAllFiles(args []string) {
	for _, fname := range args {
		augmentedData := insertOneFile(fname)
		if augmentedData == nil {
			return
		}
		fileOps.SaveToCsv(fname, augmentedData)
	}
}

func insertOneFile3(fname string) [][]string {
	var dummy [][]string
	return dummy
}


func insertOneFile(fname string) [][]string {
	allMissingDates := searchOneFile(fname)
	if len(allMissingDates) == 0 {
		return nil
	}

	input := fileOps.CsvOpenRead(fname)
	allRecords, err := input.ReadAll()
	if err != nil {
		log.Fatalf("Can not read file: '%s'\n%s\n", fname, err)
	}

	allCsvDates := getCsvDates(allRecords)
	var augmentedData [][]string

	m := 0
	for _, csvDate := range allCsvDates {
		trunc := allMissingDates[m][:18]
		missing, err := goment.New(trunc) //FIXME
		if err != nil {
			log.Fatalf("Error #69805: Invalid date/time: '%s'; %s\n", trunc, err)
		}
		if isSameOrBefore(csvDate, *missing) {
			continue
		}
		fmt.Println("do something with:", csvDate.Format(dateOutputFmt), missing.Format(dateOutputFmt), m)
		m += 1
	}




	return augmentedData
}
