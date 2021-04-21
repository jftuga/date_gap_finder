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
		insertOneFile(fname)
	}
}

func insertOneFile(fname string) {
	allMissingDates, layout := searchOneFile(fname)
	fmt.Println("layout:", layout)
	if len(allMissingDates) == 0 {
		return
	}

	input := fileOps.CsvOpenRead(fname)
	allRecords, err := input.ReadAll()
	if err != nil {
		log.Fatalf("Can not read file: '%s'\n%s\n", fname, err)
	}

	m := 0
	i := 0
	outputFmt := "L LTS dddd"
	var augmentedData [][]string

	for _, current := range allRecords {
		if allRootOptions.HasHeader && i == 0 {
			i += 1
			continue
		}
		i += 1

		fmt.Println(i, current)
		if i == 5 {
			fmt.Println("debug:", current)
		}

		currentDateTime, err := goment.New(current[allRootOptions.Column])
		if err != nil {
			log.Fatalf("Can initialize goment struct for: '%s'\n%s",current, err)
		}
		if currentDateTime.IsSameOrBefore(allMissingDates[m]) {
			fmt.Printf("IsSameOrBefore: %s - %s\n", current, currentDateTime.Format((outputFmt)))
			augmentedData = append(augmentedData, current)
			continue
		}
		//fmt.Println("Newer: ", currentDateTime.Format(outputFmt))
		fmt.Println("Newer: ", allMissingDates[m].Format(outputFmt))
		missedDate := allMissingDates[m].ToTime().Format(layout)
		missingRecord := []string {missedDate}
		augmentedData = append(augmentedData, missingRecord)
		augmentedData = append(augmentedData, current)
		m += 1

		fmt.Println("-----------------------------------------------")
	}

	fmt.Println()
	fmt.Println("=========================================================")
	fmt.Println()

	output := fileOps.CsvOpenWriteBuf()
	err = output.WriteAll(augmentedData)
	if err != nil {
		log.Fatalf("Unable to save CSV data: %s\n",err)
	}
	for _, rec := range augmentedData {
		fmt.Println(rec)
	}
}
