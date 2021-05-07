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
	"github.com/spf13/cobra"
	"log"
	"sort"
	"strings"
)

type insertOptions struct {
	columnInserts []string
	allColumnInserts string
	Overwrite bool
	MaxBackupFiles int
}

var allInsertOptions insertOptions

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "insert missing CSV entries",
	Long: `CSV dates are assumed to be sorted from oldest to newest within the file.
Multiple -r options can be used.  Each -r option is comma-delimited with
the column number first and the value to insert (into that column) second.`,
	Run: func(cmd *cobra.Command, args []string) {
		insertAllFiles(args)
	},
}

func init() {
	rootCmd.AddCommand(insertCmd)
	insertCmd.Flags().StringArrayVarP(&allInsertOptions.columnInserts, "record", "r", []string{}, "insert record with missing data; format: col#,value")
	insertCmd.Flags().StringVarP(&allInsertOptions.allColumnInserts, "allRecords", "R", "", "insert data to all columns of a missing row")
	insertCmd.PersistentFlags().BoolVarP(&allInsertOptions.Overwrite, "overwrite", "O", false, "overwrite existing CSV file; original file saved as .bak")
	insertCmd.Flags().IntVarP(&allInsertOptions.MaxBackupFiles, "max", "m", -1,"max number of backup files to save; -1=save all")
}

func insertAllFiles(args []string) {
	for _, fname := range args {
		augmentedData := insertOneFile(fname)
		if augmentedData == nil {
			return
		}
		if allInsertOptions.Overwrite {
			ok, _ := fileOps.OverwriteCsv(fname, augmentedData)
			if ok && allInsertOptions.MaxBackupFiles > -1 {
				fileOps.RemoveOldBackups(fname, allInsertOptions.MaxBackupFiles)
			}
			continue
		}
		for _, aug := range augmentedData {
			fmt.Println(aug)
		}
	}
}

func insertOneFile(fname string) []string {
	debug := allRootOptions.Debug
	allMissingDates, _ := SearchOneFile(fname)
	if len(allMissingDates) == 0 {
		return nil
	}

	input, file := fileOps.CsvOpenRead(fname)
	var r []rune
	if allRootOptions.TabDelimiter {
		allRootOptions.CsvDelimiter = "\\t"
	}
	if allRootOptions.CsvDelimiter == `\t` {
		r = []rune{'\t'}
	} else {
		r = []rune(allRootOptions.CsvDelimiter)
	}
	input.Comma = r[0]
	allRecords, err := input.ReadAll()
	if err != nil {
		log.Fatalf("Can not read file: '%s'\n%s\n", fname, err)
	}
	if debug > 9998 {
		fmt.Println("allRecords len:", len(allRecords))
	}
	err = file.Close()
	if err != nil {
		log.Fatalf("Error #15265: Unable to close CSV file: '%s'; %s\n", fname, err)
		return nil
	}

	row := 0
	if !allRootOptions.HasNoHeader {
		row = 1
	}
	layout := allRecords[row][allRootOptions.Column]
	numOfColumns := len(allRecords[row])
	if debug > 9998 {
		fmt.Println("layout:", layout)
	}

	allCsvDates, allRows := getCsvDates(allRecords)
	if debug > 9998 {
		fmt.Printf("allCsvDates len: %d, allRows len: %d\n", len(allCsvDates), len(allRows))
	}

	for _, m := range allMissingDates {
		csvStyleDate := ConvertDate(m.ToTime(), layout)
		newRow := createNewRow(csvStyleDate, numOfColumns)
		allRecords = append(allRecords, newRow)
	}

	if debug > 9998 {
		fmt.Println()
		fmt.Println("allRecords")
		fmt.Println("=============")
		for _, rec := range allRecords {
			fmt.Println(rec)
		}
	}

	sortedRecords := allRecords
	sortRecords(sortedRecords)
	if debug > 9998 {
		fmt.Println()
		fmt.Println("sortedRecords")
		fmt.Println("=============")
		for _, rec := range sortedRecords {
			fmt.Println(rec)
		}
	}

	var csvRecords []string
	var headerRow []string
	delimiter := allRootOptions.CsvDelimiter
	if delimiter == "\\t" {
		delimiter = fmt.Sprintf("%c", 0x09)
	}
	for i, rec := range sortedRecords {
		if !allRootOptions.HasNoHeader && i == len(sortedRecords)-1 {
			headerRow = rec
			//fmt.Println("headerRow:", headerRow)
			continue
		}
		csvRecords = append(csvRecords, strings.Join(rec,delimiter))
	}
	if !allRootOptions.HasNoHeader {
		csvRecords = append([]string {strings.Join(headerRow,delimiter)}, csvRecords...)
	}

	if debug > 9998 {
		fmt.Println()
		fmt.Println("csvRecords")
		fmt.Println("=============")
		for _, rec := range csvRecords {
			fmt.Println(rec)
		}
	}

	if debug > 9998 {
		fmt.Println()
		fmt.Println("==========================================================")
		fmt.Println()
	}
	return csvRecords
}

func sortRecords(entry [][]string) {
	sort.SliceStable(entry, func(i, j int) bool {
		return entry[i][allRootOptions.Column] < entry[j][allRootOptions.Column]
	})
}

func createNewRow(missedDate string, numOfColumns int) []string {
	debug := allRootOptions.Debug
	missingRecord := make(map [int]string)
	missingRecord[allRootOptions.Column] = missedDate
	for _, column := range allInsertOptions.columnInserts {
		if debug > 9998 {
			fmt.Println("kv:",column)
		}
		col, val := GetKeyVal(column)
		missingRecord[col] = val
	}
	if debug > 9998 {
		fmt.Println("missingRecord:", missingRecord)
	}
	keys, last := SortIntMapByKey(missingRecord)
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

	if len(newRow)+1 <= numOfColumns {
		for i := len(newRow); i < numOfColumns; i++ {
			newRow = append(newRow, "")
			if debug > 9998 {
				fmt.Printf("[%d] adding column to newRow\n", i)
			}
		}
	}

	if len(allInsertOptions.allColumnInserts) > 0 {
		for i := 0; i < numOfColumns; i ++ {
			if newRow[i] == "" {
				newRow[i] = allInsertOptions.allColumnInserts
			}
		}
	}

	return newRow
}
