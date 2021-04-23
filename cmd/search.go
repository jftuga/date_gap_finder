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
	"encoding/csv"
	"fmt"
	"github.com/jftuga/date_gap_finder/fileOps"
	"github.com/jftuga/date_gap_finder/shared"
	"github.com/nleeper/goment"
	"github.com/spf13/cobra"
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
	//outputFmt := "L LTS dddd"
	for _, fname := range args {
		//missingDates, _ := searchOneFile(fname)
		searchOneFile(fname)
		/*
		for _,d := range missingDates {
			fmt.Printf("[159] missing: %s\n", d.Format(outputFmt))
		} */
	}
}

func searchOneFile(fname string) /*([]*goment.Goment, string)*/ {
	fileOps.CsvOpenRead(fname)
	input := fileOps.CsvOpenRead(fname)
	searchFromReader(input, fname)
}

func searchFromReader(input *csv.Reader, streamName string)  {
	//debug := allRootOptions.Debug
	outputFmt := "L LTS dddd"

	allRecords, err := input.ReadAll()
	fmt.Println("allRecords:", len(allRecords))
	if err != nil {
		log.Fatalf("Error #89533: Unable to read from stream: '%s'; %s\n", streamName, err)
	}

	var csvDates []goment.Goment
	for i, d := range allRecords {
		if allRootOptions.HasHeader && i == 0 {
			continue
		}
		g, err := goment.New(d[allRootOptions.Column])
		if err != nil {
			log.Fatalf("Error #30425: Invalid data/time: '%s'; %s\n", d[allRootOptions.Column], err)
		}
		csvDates = append(csvDates,*g)
	}

	fmt.Println("csvDates")
	fmt.Println("========")
	for _, d := range csvDates {
		fmt.Println(d.Format(outputFmt))
	}


	f := 0

	if allRootOptions.HasHeader {
		f = 1
	}
	firstRec := allRecords[f]
	first, err := goment.New(firstRec[allRootOptions.Column])
	if err != nil {
		log.Fatalf("Error #30430: Invalid data/time: '%s'; %s\n", firstRec[allRootOptions.Column], err)
	}
	lastRec := allRecords[len(allRecords)-1]
	last, err := goment.New(lastRec[allRootOptions.Column])
	if err != nil {
		log.Fatalf("Error #30435: Invalid data/time: '%s'; %s\n", lastRec[allRootOptions.Column], err)
	}


	var requiredDates []goment.Goment
	durationInSeconds := shared.GetDuration(allRootOptions.Amount, allRootOptions.Period)
	current, _ := goment.New(first)
	for {
		if current.IsAfter(last) {
			break
		}
		requiredDates = append(requiredDates, *current)
		current.Add(durationInSeconds, "seconds")
	}

	fmt.Println()
	fmt.Println("requiredDates")
	fmt.Println("=============")
	for _, d := range requiredDates {
		fmt.Println(d.Format(outputFmt))
	}

	fmt.Println()
	fmt.Println("comparison")
	fmt.Println("==========")
	for _, reqDate := range requiredDates {
		for i, csvDate := range csvDates {
			result := csvDate.IsSameOrBefore(&reqDate)
			fmt.Println("csv:", csvDate.Format(outputFmt), "[sameOrBefore]", "req:", reqDate.Format(outputFmt), "=>", result)
			if result {
				csvDates = shared.RemoveIndex(csvDates,i)
				break
			} else {
				fmt.Println("missing date:", reqDate.Format(outputFmt))
				break
			}
		}
		fmt.Println("---------------")
	}

}
