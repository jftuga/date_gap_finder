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
	"github.com/jftuga/date_gap_finder/fileOps"
	"github.com/spf13/cobra"
)

type insertOptions struct {
	columnInserts []string
}

var allInsertOptions insertOptions

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
	Use:   "insert",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

func insertOneFile(fname string) [][]string {
	var dummy [][]string
	return dummy
}

/*
func insertOneFile2(fname string) [][]string {
	debug := false
	allMissingDates, layout := searchOneFile(fname)
	if debug {
		fmt.Println("layout:", layout)
	}
	if len(allMissingDates) == 0 {
		return nil
	}

	input := fileOps.CsvOpenRead(fname)
	allRecords, err := input.ReadAll()
	if err != nil {
		log.Fatalf("Can not read file: '%s'\n%s\n", fname, err)
	}

	m := 0
	outputFmt := "L LTS dddd"
	var augmentedData [][]string

	for i, current := range allRecords {
		if allRootOptions.HasHeader && i == 0 {
			continue
		}

		currentDateTime, err := goment.New(current[allRootOptions.Column])
		if err != nil {
			log.Fatalf("Can initialize goment struct for: '%s'\n%s",current, err)
		}

		if m < len(allMissingDates) && currentDateTime.IsSameOrBefore(allMissingDates[m]) {
			if debug {
				fmt.Printf("IsSameOrBefore: %s - %s\n", current, currentDateTime.Format(outputFmt))
			}
			augmentedData = append(augmentedData, current)
			continue
		}

		if debug && m < len(allMissingDates) {
			fmt.Println("Newer: ", allMissingDates[m].Format(outputFmt))
		}

		fmt.Println("iiiiiiiiiii:", i, len(allMissingDates), m)
		if m == len(allMissingDates) {
			fmt.Println("DONE")
			continue
		}
		missedDate := allMissingDates[m].ToTime().Format(layout)
		missingRecord := make(map [int]string)
		missingRecord[allRootOptions.Column] = missedDate
		for _, column := range allInsertOptions.columnInserts {
			if debug {
				fmt.Println("kv:",column)
			}
			col, val := shared.GetKeyVal(column)
			missingRecord[col] = val
		}
		if debug {
			fmt.Println("missingRecord:", missingRecord)
		}
		keys, last := shared.SortIntMapByKey(missingRecord)
		if debug {
			fmt.Println("keys, last:", keys, last)
		}
		var newRow []string
		for i=0; i <= last; i++{
			if val, ok := missingRecord[i]; ok {
				newRow = append(newRow, val)
				if debug {
					fmt.Println("appending:", val)
				}
			} else {
				newRow = append(newRow, "")
			}
		}
		newRow[allRootOptions.Column] = missedDate
		augmentedData = append(augmentedData, newRow)
		augmentedData = append(augmentedData, current)
		m += 1

		if debug {
			fmt.Println("-----------------------------------------------")
		}
	}

	output := fileOps.CsvOpenWriteBuf()
	err = output.WriteAll(augmentedData)
	if err != nil {
		log.Fatalf("Unable to save CSV data: %s\n",err)
	}

	if debug {
		for _, rec := range augmentedData {
			fmt.Println(rec)
		}
	}
	return augmentedData
}
*/