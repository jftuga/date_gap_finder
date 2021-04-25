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
FITNESS FOR A PARTICULAR PURPOSE AND NON-INFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"github.com/jftuga/date_gap_finder/fileOps"
	"github.com/jftuga/date_gap_finder/shared"
	"github.com/spf13/cobra"
	"log"
	"sort"
	"strings"
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
		//fileOps.SaveToCsv(fname, augmentedData)
	}
}

func insertOneFile3(fname string) [][]string {
	var dummy [][]string
	return dummy
}


func insertOneFile(fname string) []string {
	allMissingDates := searchOneFile(fname)
	if len(allMissingDates) == 0 {
		return nil
	}

	input := fileOps.CsvOpenRead(fname)
	allRecords, err := input.ReadAll()
	if err != nil {
		log.Fatalf("Can not read file: '%s'\n%s\n", fname, err)
	}
	fmt.Println("allRecords len:", len(allRecords))
	row := 0
	if allRootOptions.HasHeader {
		row = 1
	}
	layout := allRecords[row][allRootOptions.Column]
	fmt.Println("layout:", layout)

	allCsvDates, allRows := getCsvDates(allRecords)
	fmt.Printf("allCsvDates len: %d, allRows len: %d\n", len(allCsvDates), len(allRows))

	for _, m := range allMissingDates {
		csvStyleDate := shared.ConvertDate(m.ToTime(), layout)
		newRow := createNewRow(csvStyleDate)
		allRecords = append(allRecords, newRow)
	}

	fmt.Println()
	fmt.Println("allRecords")
	fmt.Println("=============")
	for _, rec := range allRecords {
		fmt.Println(rec)
	}

	var csvRecords []string
	var headerRow []string
	for i, rec := range allRecords {
		if allRootOptions.HasHeader && i == 0 {
			headerRow = rec
			continue
		}
		csvRecords = append(csvRecords, strings.Join(rec,","))
	}
	sortRecords(csvRecords)
	if allRootOptions.HasHeader {
		csvRecords = append([]string {strings.Join(headerRow,",")}, csvRecords...)
	}

	fmt.Println()
	fmt.Println("csvRecords")
	fmt.Println("=============")
	for _, rec := range csvRecords {
		fmt.Println(rec)
	}

	return csvRecords
}

func sortRecords(entry []string) {
	sort.Slice(entry, func(i, j int) bool {
		return entry[i] < entry[j]
	})
}

func createNewRow(missedDate string) []string {
	debug := allRootOptions.Debug
	missingRecord := make(map [int]string)
	missingRecord[allRootOptions.Column] = missedDate
	for _, column := range allInsertOptions.columnInserts {
		if debug > 9998 {
			fmt.Println("kv:",column)
		}
		col, val := shared.GetKeyVal(column)
		missingRecord[col] = val
	}
	if debug > 9998 {
		fmt.Println("missingRecord:", missingRecord)
	}
	keys, last := shared.SortIntMapByKey(missingRecord)
	if debug > 9998 {
		fmt.Println("keys, last:", keys, last)
	}
	var newRow []string
	for i:=0; i <= last; i++{
		if val, ok := missingRecord[i]; ok {
			newRow = append(newRow, val)
			if debug > 9998 {
				fmt.Println("appending:", val)
			}
		} else {
			newRow = append(newRow, "")
		}
	}
	newRow[allRootOptions.Column] = missedDate
	return newRow
}