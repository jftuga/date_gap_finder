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
	"time"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search CSV files for missing dates",
	Long: `CSV dates are assumed to be sorted from oldest to newest within the file.`,
	Run: func(cmd *cobra.Command, args []string) {
		//p := profile.Start(profile.CPUProfile, profile.ProfilePath("."))
		//p := profile.Start(profile.MemProfile, profile.MemProfileRate(512), profile.ProfilePath("."))
		total := searchAllFiles(args)
		//p.Stop()
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

func show(req, csv, missing []goment.Goment, c, r int) {
	DisplayTable(missing,"seendate", false, -1)
	DisplayTable(csv,"csvdate", false, c)
	DisplayTable(req,"reqdate", true, r)
	fmt.Println("======================================================")
}

func getExtendedTime(originalTime goment.Goment) *goment.Goment {
	if len(allRootOptions.Padding) == 0 {
		return &originalTime
	}

	var err error
	padTime, _ := time.ParseDuration("0s")
	padTime, err = time.ParseDuration(allRootOptions.Padding)
	if err != nil {
		log.Fatalf("Error #23606: unable to convert to time.Duration: '%s, %s'\n", allRootOptions.Padding, err)
	}

	newTime := originalTime.ToTime()
	newTime = newTime.Add(padTime)
	g, err := goment.New(newTime)
	if err != nil {
		log.Fatalf("Error #23819: unable to create goment Time: '%s, %s'\n", newTime, err)
	}
	return g
}

func GetPaddingRange(g goment.Goment) (goment.Goment, goment.Goment){
	padTime, err := time.ParseDuration(allRootOptions.Padding)
	if err != nil {
		log.Fatalf("Error #90985: unable to create time duration for: %s; %s\n", allRootOptions.Padding, err)
	}
	a, err := goment.New(g.ToTime())
	if err != nil {
		log.Fatalf("Error #90990: unable to create time duration for: %s; %s\n", g.Format(dateOutputFmt), err)
	}
	b, err := goment.New(g.ToTime())
	if err != nil {
		log.Fatalf("Error #90995: unable to create time duration for: %s; %s\n", g.Format(dateOutputFmt), err)
	}
	a.Subtract(padTime)
	b.Add(padTime)
	return *a, *b
}

func ExcludeDates(reqDate []goment.Goment) []goment.Goment {
	if allRootOptions.SkipWeekends == false && len(allRootOptions.SkipDays) == 0 {
		return reqDate
	}

	var allSkipDaysLower string
	if len(allRootOptions.SkipDays) > 0 {
		allSkipDaysLower = strings.ToLower(allRootOptions.SkipDays)
	}

	var included []goment.Goment
	for _, req := range reqDate {
		if allRootOptions.SkipWeekends && (req.Format("dddd") == "Saturday" || req.Format("dddd") == "Sunday") {
			continue
		}
		skipDayLower := strings.ToLower(req.Format("dddd"))
		if len(allRootOptions.SkipDays) > 0 && strings.Index(allSkipDaysLower,skipDayLower) >= 0 {
			continue
		}
		included = append(included, req)
	}
	return included
}

func IsNear(csv goment.Goment, reqDate []goment.Goment) (bool, int) {
	for r, req := range reqDate {
		a, b := GetPaddingRange(req)
/*		fmt.Println("  a:", a.Format(dateOutputFmt))
		fmt.Println("  b:", b.Format(dateOutputFmt))
		fmt.Println("csv:", csv.Format(dateOutputFmt))*/
		if csv.IsBetween(&a, &b) {
			//fmt.Printf("%s is near %s\n", csv.Format(dateOutputFmt), req.Format(dateOutputFmt))
			return true, r
		}
	}
	return false, -1
}

func findMissingDates(csvDate, reqDate []goment.Goment) []goment.Goment {
	for _, csv := range csvDate {
		found, r := IsNear(csv, reqDate)
		if found {
			reqDate = RemoveSliceItem(reqDate,r)
		}
	}

/*	fmt.Println()
	DisplayTable(csvDate,"csvDate", false, -1)
	DisplayTable(reqDate,"reqDate", false, -1)*/

	filterDate := ExcludeDates(reqDate)
	//DisplayTable(filterDate,"filterDate", false, -1)

	return filterDate
}

func getCsvDates(allRecords [][]string) ([]goment.Goment, map[string][]string) {
	// build csvDates
	var csvDates []goment.Goment
	allRows := make(map [string][]string)
	for i, d := range allRecords {
		if !allRootOptions.HasNoHeader && i == 0 {
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
	debugLevel := allRootOptions.Debug

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
	if !allRootOptions.HasNoHeader {
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
	durationInSeconds := GetDurationInSeconds(allRootOptions.Amount, allRootOptions.Unit)
	current, _ := goment.New(first)
	for {
		if current.IsAfter(last) {
			break
		}
		if debugLevel > 99998 {
			fmt.Println("current:", current.Format(dateOutputFmt))
		}
		requiredDates = append(requiredDates, *current)
		current.Add(durationInSeconds, "seconds")
	}
	if current.IsBetween(&last) {
		requiredDates = append(requiredDates, *current)
	}

	return csvDates, requiredDates, csvStyleDate
}
