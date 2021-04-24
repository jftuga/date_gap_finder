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
	"strings"
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
	outputFmt := "L LTS dddd"
	for _, fname := range args {
		missingDates := searchOneFile(fname)
		for i,d := range missingDates {
			fmt.Printf("[%d] missing: %s\n", i+1, d.Format(outputFmt))
		}
	}
}

func searchOneFile(fname string) []*goment.Goment {
	debug := allRootOptions.Debug
	outputFmt := "L LTS dddd"

	fileOps.CsvOpenRead(fname)
	input := fileOps.CsvOpenRead(fname)
	csvDates, requiredDates := getCsvAndRequiredDates(input, fname)

	if debug > 98 {
		fmt.Println("csvDates")
		fmt.Println("========")
		for _, d := range csvDates {
			fmt.Println(d.Format(outputFmt))
		}
	}

	if debug > 98  {
		fmt.Println()
		fmt.Println("requiredDates")
		fmt.Println("=============")
		for _, d := range requiredDates {
			fmt.Println(d.Format(outputFmt))
		}
	}

	findMissingDates(csvDates, requiredDates)

	return nil
}

func isSameOrBefore(csvDate, reqDate goment.Goment) bool {
	return csvDate.IsSameOrBefore(&reqDate)
}

func findMissingDates(csvDates, requiredDates []goment.Goment) {
	debug := allRootOptions.Debug
	outputFmt := "L LTS dddd"

	maxDiff := shared.GetDuration(allRootOptions.Amount, allRootOptions.Period)
	fmt.Println("maxDiff:", maxDiff)
	// reqDate=key; csvDate(s)=val
	allCsvBeforeRequired := make(map [string][]string)
	seenCsvDates := make(map [string]bool)
	for _, reqDate := range requiredDates {
		fmt.Println("checking:", reqDate.Format(outputFmt))
		for _, csvDate := range csvDates {
			if isSameOrBefore(csvDate, reqDate) {
				key := reqDate.Format(outputFmt)
				val := csvDate.Format(outputFmt)
				// FIXME:
				diff := shared.GetTimeDifference(csvDate,reqDate)
				fmt.Println("diff:", diff)
				allCsvBeforeRequired[key] = append(allCsvBeforeRequired[key], diff.String())
				if diff.Seconds() < maxDiff.Seconds() {
					fmt.Println("xxxxx:", key, val)
					seenCsvDates[key] = true
				}

			}
		}
	}

	if debug > 98 {
		fmt.Println()
		fmt.Println("allCsvBeforeRequired")
		fmt.Println("====================")
		for key, reqDate := range allCsvBeforeRequired {
			fmt.Printf("%32s => %s\n", key, strings.Join(reqDate, "; "))
		}
	}

	if debug > 98 {
		fmt.Println()
		fmt.Println("seenCsvDates")
		fmt.Println("============")
		for key := range seenCsvDates {
			fmt.Println(key)
		}
	}

	fmt.Println()
	fmt.Println("MissingDates")
	fmt.Println("============")
	for _, reqDate := range requiredDates {
		toCheck := reqDate.Format(outputFmt)
		//fmt.Println()
		//fmt.Println("toCheck:", toCheck)
		_, ok := seenCsvDates[toCheck]
		if !ok {
			fmt.Printf("missing date: %s\n", toCheck)
		}
	}
}

func getCsvAndRequiredDates(input *csv.Reader, streamName string) ([]goment.Goment, []goment.Goment) {
	debug := allRootOptions.Debug

	allRecords, err := input.ReadAll()
	if err != nil {
		log.Fatalf("Error #89533: Unable to read from stream: '%s'; %s\n", streamName, err)
	}
	if debug > 98 {
		fmt.Println("allRecords:", len(allRecords))
	}

	// build csvDates
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

	// build requiredDates
	f := 0
	if allRootOptions.HasHeader {
		f = 1
	}
	if f >= len(allRecords) {
		log.Fatalf("Error #98450: CSV file only contains '%d' records: '%s'\n", len(allRecords), streamName)
	}
	firstRec := allRecords[f]
	first, err := goment.New(firstRec[allRootOptions.Column])
	if err != nil {
		log.Fatalf("Error #30430: Invalid data/time: '%s'; %s\n", firstRec[allRootOptions.Column], err)
	}
	lastRec := allRecords[len(allRecords)-1]
	last, err := goment.New(lastRec[allRootOptions.Column])
	//fmt.Println("last:", last.Format("L LTS dddd"))

	if err != nil {
		log.Fatalf("Error #30435: Invalid data/time: '%s'; %s\n", lastRec[allRootOptions.Column], err)
	}

	var requiredDates []goment.Goment
	durationInSeconds := shared.GetDurationInSeconds(allRootOptions.Amount, allRootOptions.Period)
	current, _ := goment.New(first)
	for {
		if current.IsAfter(last) {
			break
		}
		requiredDates = append(requiredDates, *current)
		current.Add(durationInSeconds, "seconds")
	}
	requiredDates = append(requiredDates, *current)

	return csvDates, requiredDates
}

func checkCSVDate(csvDates []goment.Goment, reqDate *goment.Goment) []*goment.Goment {
	debug := allRootOptions.Debug
	outputFmt := "L LTS dddd"

	var allMissingDates []*goment.Goment
	for i, csvDate := range csvDates {
		result := csvDate.IsSameOrBefore(reqDate)
		if debug > 98  {
			fmt.Println("csv:", csvDate.Format(outputFmt), "[sameOrBefore]", "req:", reqDate.Format(outputFmt), "=>", result)
		}
		if result {
			fmt.Println("Removing from CSV:", csvDates[i].Format(outputFmt))
			csvDates = shared.RemoveIndex(csvDates,i)
			break
		} else {
			if debug > 98  {
				fmt.Println("missing date:", reqDate.Format(outputFmt))
			}
			allMissingDates = append(allMissingDates, reqDate)
			break
		}
	}
	if debug > 98  {
		fmt.Println("---------------")
	}
	return allMissingDates
}
