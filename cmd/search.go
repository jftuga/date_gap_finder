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
	"github.com/araddon/dateparse"
	"github.com/jftuga/date_gap_finder/fileOps"
	"github.com/jftuga/date_gap_finder/shared"
	"github.com/nleeper/goment"
	"github.com/spf13/cobra"
	"io"
	"log"
)

/*
type searchOptions struct {
	Column int
	HasHeader bool
	Amount int
	Period string
}

var allSearchOptions searchOptions
*/
// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		//fmt.Println("search called")
		searchAllFiles(args)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	/*
	searchCmd.PersistentFlags().IntVarP(&allSearchOptions.Column, "column", "c", 0, "CSV column number (starts at zero)")
	searchCmd.PersistentFlags().BoolVarP(&allSearchOptions.HasHeader, "header", "H", true, "if CSV file has header line")
	searchCmd.PersistentFlags().IntVarP(&allSearchOptions.Amount, "amount", "a", -1, "a maximum, numeric duration")
	searchCmd.PersistentFlags().StringVarP(&allSearchOptions.Period, "period", "p", "", "period of time, such as: days, hours, minutes")
	*/

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func searchAllFiles(args []string) {
	for _, fname := range args {
		searchOneFile(fname)
	}
}

func searchOneFile(fname string) ([]*goment.Goment, string) {
	debug := false
	outputFmt := "L LTS dddd"
	input := fileOps.CsvOpenRead(fname)
	i := 0
	var previous *goment.Goment
	var layout string
	previous, _ = goment.New("")

	var missingDates []*goment.Goment
	for {
		record, err := input.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Unable to read record from file: '%s'\n%s\n", fname, err)
		}
		if allRootOptions.HasHeader && i == 0 {
			i += 1
			continue
		}
		i += 1
		currentTimeStamp := record[allRootOptions.Column]
		current, err := goment.New(currentTimeStamp)
		if err != nil {
			log.Fatalf("error with timestamp: '%s'\n%s\n", currentTimeStamp,err)
		}

		// get the date layout from the first row of CSV data
		if i == 2 {
			layout, err = dateparse.ParseFormat(currentTimeStamp)
			if err != nil {
				log.Fatalf("Can not parse date time for: '%s'\n%s\n", currentTimeStamp, err)
			}
		}
		hasGap, gapOnWeekday := shared.DatesHaveGaps(previous, current, allRootOptions.Amount, allRootOptions.Period)
		if hasGap && gapOnWeekday {
			if debug {
				fmt.Printf("missing date: '%s'\n", previous.Format(outputFmt))
			}
			missingDates = append(missingDates, previous)
		}
		previous = current
	}
	return missingDates, layout
}
