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
	"github.com/nleeper/goment"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search CSV files for missing dates",
	Long: `CSV dates are assumed to be sorted from oldest to newest within the file.`,
	Run: func(cmd *cobra.Command, args []string) {
		total := searchAllFiles(args)
		os.Exit(total)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}

func searchAllFiles(args []string) int {
	total := 0
	for _, fname := range args {
		missingDates, csvStyleDate := SearchOneFile(fname)
		for _,d := range missingDates {
			missingFormatted := ConvertDate(d.ToTime(), csvStyleDate)
			fmt.Println(missingFormatted)
		}
		total += len(missingDates)
	}
	return total
}

func SearchOneFile(fname string) ([]goment.Goment, string) {
	debugLevel := allRootOptions.Debug
	// why? fileOps.CsvOpenRead(fname)
	input, file := fileOps.CsvOpenRead(fname)
	var r []rune
	if allRootOptions.CsvDelimiter == `\t` {
		r = []rune{'\t'}
	} else {
		r = []rune(allRootOptions.CsvDelimiter)
	}
	input.Comma = r[0]
	csvDates, requiredDates, csvStyleDate := getCsvAndRequiredDates(input, fname)
	file.Close()

	if debugLevel > 98 {
		fmt.Println("csvDates")
		fmt.Println("========")
		for _, d := range csvDates {
			fmt.Println(d.Format(dateOutputFmt))
		}
	}

	if debugLevel > 98  {
		fmt.Println()
		fmt.Println("requiredDates")
		fmt.Println("=============")
		for _, d := range requiredDates {
			fmt.Println(d.Format(dateOutputFmt))
		}
	}

	return findMissingDates(csvDates, requiredDates), csvStyleDate
}

func isSameOrBefore(csvDate, reqDate goment.Goment) bool {
	return csvDate.IsSameOrBefore(&reqDate)
}

func findMissingDates(csvDates, requiredDates []goment.Goment) []goment.Goment {
	debugLevel := allRootOptions.Debug

	maxTimeDiff := GetDuration(allRootOptions.Amount, allRootOptions.Period)
	seenDates := make(map [string]bool)
	for _, reqDate := range requiredDates {
		for _, csvDate := range csvDates {
			if isSameOrBefore(csvDate, reqDate) {
				key := reqDate.Format(dateOutputFmt)
				// compare the time duration difference
				diff := GetTimeDifference(csvDate,reqDate)
				if diff.Seconds() < maxTimeDiff.Seconds() {
					seenDates[key] = true
				}

			}
		}
	}

	if debugLevel > 98 {
		fmt.Println()
		fmt.Println("seenDates")
		fmt.Println("============")
		for key := range seenDates {
			fmt.Println(key)
		}
	}

	if debugLevel > 98 {
		fmt.Println()
		fmt.Println("MissingDates")
		fmt.Println("============")
	}

	var allSkipDaysLower string
	if len(allRootOptions.SkipDays) > 0 {
		allSkipDaysLower = strings.ToLower(allRootOptions.SkipDays)
	}
	var allMissingDates []goment.Goment
	for _, reqDate := range requiredDates {
		toCheck := reqDate.Format(dateOutputFmt)
		if allRootOptions.SkipWeekends && (reqDate.Format("dddd") == "Saturday" || reqDate.Format("dddd") == "Sunday") {
			if debugLevel > 98 {
				fmt.Println("skipping weekend:", reqDate.Format(dateOutputFmt))
			}
			continue
		}
		skipDayLower := strings.ToLower(reqDate.Format("dddd"))
		if len(allRootOptions.SkipDays) > 0 && strings.Index(allSkipDaysLower,skipDayLower) >= 0 {
			if debugLevel > 98 {
				fmt.Println("skipping day:", reqDate.Format(dateOutputFmt))
			}
			continue
		}
		_, ok := seenDates[toCheck]
		if !ok {
			allMissingDates = append(allMissingDates, reqDate)
			if debugLevel > 98 {
				fmt.Printf("missing date: %s\n", toCheck)
			}
		}
	}

	if debugLevel > 98 {
		fmt.Println()
		fmt.Println("==========================================================")
		fmt.Println()
	}

	return allMissingDates
}

func getCsvDates(allRecords [][]string) ([]goment.Goment, map[string][]string) {
	// build csvDates
	var csvDates []goment.Goment
	allRows := make(map [string][]string)
	for i, d := range allRecords {
		if allRootOptions.HasHeader && i == 0 {
			continue
		}
		g, err := goment.New(d[allRootOptions.Column])
		if err != nil {
			log.Fatalf("Error #30425: Invalid data/time: '%s'; %s\n", d[allRootOptions.Column], err)
		}
		csvDates = append(csvDates,*g)
		allRows[g.Format(dateOutputFmt)] = d
	}
	return csvDates, allRows
}

func getCsvAndRequiredDates(input *csv.Reader, streamName string) ([]goment.Goment, []goment.Goment, string) {
	allRecords, err := input.ReadAll()
	if err != nil {
		log.Fatalf("Error #89533: Unable to read from stream: '%s'; %s\n", streamName, err)
	}
	if debugLevel > 98 {
		fmt.Println("allRecords length:", len(allRecords))
	}

	csvDates, _ := getCsvDates(allRecords)

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
	if err != nil {
		log.Fatalf("Error #30435: Invalid data/time: '%s'; %s\n", lastRec[allRootOptions.Column], err)
	}

	layout := firstRec[allRootOptions.Column]
	csvStyleDate := ConvertDate(first.ToTime(), layout)

	var requiredDates []goment.Goment
	durationInSeconds := GetDurationInSeconds(allRootOptions.Amount, allRootOptions.Period)
	current, _ := goment.New(first)
	for {
		if current.IsAfter(last) {
			break
		}
		requiredDates = append(requiredDates, *current)
		current.Add(durationInSeconds, "seconds")
	}
	if current.IsBetween(&last) {
		requiredDates = append(requiredDates, *current)
	}

	return csvDates, requiredDates, csvStyleDate
}
