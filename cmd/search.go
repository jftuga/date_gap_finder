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
	"sort"
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

func findMissingDates(csvDate, requiredDates []goment.Goment) []goment.Goment {
	//debugLevel := allRootOptions.Debug
	//fmt.Println()
	//fmt.Println("=================================================================")
	//fmt.Println()
	//maxTimeDiff := GetDuration(allRootOptions.Amount, allRootOptions.Unit)
	//padTime, _ := time.ParseDuration("0s")
	//fmt.Println("maxTimeDiff:", maxTimeDiff)
	//fmt.Println("    padTime:", padTime)
	//fmt.Println()

	var seenDates []goment.Goment
	//csvDate := csvDates
	c := len(csvDate) - 1
	r := 0
	reqDate := requiredDates

	for {
		for {
			fmt.Println()
			fmt.Println("========================================================")
			fmt.Println("c, len(csvDate), r, len(reqDate)", c, len(csvDate), r, len(reqDate))
			fmt.Printf("about to compare csv:%s\n", csvDate[c].Format(dateOutputFmt))
			fmt.Printf("                 req:%s\n", reqDate[r].Format(dateOutputFmt))
			//show(reqDate, csvDate, seenDates, c, r)
			if csvDate[c].IsSameOrBefore(&reqDate[r]) {
				fmt.Printf("added to seenDate: %s   c:%d", reqDate[r].Format(dateOutputFmt), c)
				seenDates = append(seenDates, reqDate[r])
				reqDate = RemoveSliceItem(reqDate, r)
				csvDate = BeheadSlice(csvDate, c)
				if len(reqDate) == 0 {
					break
				}
				c = len(csvDate)
			}
			c -= 1
			if c == -1 {
				break
			}
		} // for inner
		r += 1
		c = len(csvDate) - 1
		if c == -1 {
			break
		}
		if r == len(reqDate) {
			break
		}
		fmt.Println("22 c, len(csvDate), r, len(reqDate)", c, len(csvDate), r, len(reqDate))
	}


	fmt.Println()
	fmt.Println("vvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvvv")
	DisplayTable(seenDates,"seenDates", false, -1)
	DisplayTable(csvDate,"csvDate", false, -1)
	DisplayTable(reqDate,"reqDate", false, -1)

	/*
	for {
		fmt.Println("c, len(csvDate), r, len(reqDate)", c, len(csvDates), r, len(reqDate))
		show(reqDate,csvDate,missingDates)
		if len(reqDate) == 0 {
			break
		}
		if csvDate[c].IsSameOrBefore(&reqDate[r]) {
			reqDate = RemoveSliceItem(reqDate,r)
			csvDate = RemoveSliceItem(csvDate,c)
			c = -1
			r = -1
		} else {
			missingDates = append(missingDates, reqDate[r])
			fmt.Println("zzzzzzzzzzzzzzzzzzzzzzzzzz")
			//reqDate = RemoveSliceItem(reqDate,r)
			csvDate = RemoveSliceItem(csvDate,c)
			c = -1
			r = -1
		}
		c += 1
		r += 1
	}
	*/

	/*
	for {
		show(reqDate,csvDate)
		fmt.Printf("cmp c:%s\n", csvDate[c].Format(dateOutputFmt))
		fmt.Printf("    r:%s\n", reqDate[r].Format(dateOutputFmt))
		fmt.Println()
		if csvDate[c].IsAfter(&reqDate[r]) {
			fmt.Println("BOOM!")
			missingDates = append(missingDates, reqDate[r])
			reqDate = RemoveSliceItem(reqDate,r)
			fmt.Println("xxxxxxxxxxxxxx:", len(reqDate))
			r = -1
		}
		c += 1
		r += 1
		fmt.Println("r:", r, len(reqDate), " c:", c, len(csvDates))
		if r == len(reqDate) {
			break
		}
	}
	*/

	/*
	fmt.Println("reqDates")
	fmt.Println("============")
	for _, m := range reqDate {
		fmt.Println(m.Format(dateOutputFmt))
	}
	fmt.Println()

	fmt.Println("missingDates")
	fmt.Println("============")
	for _, m := range missingDates {
		fmt.Println(m.Format(dateOutputFmt))
	}

	fmt.Println("xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx")

	 */
	return reqDate
}

func findMissingDates2(csvDates, requiredDates []goment.Goment) []goment.Goment {
	debugLevel := allRootOptions.Debug

	maxTimeDiff := GetDuration(allRootOptions.Amount, allRootOptions.Unit)
	padTime, _ := time.ParseDuration("0s")
	var err error
	if len(allRootOptions.Padding) > 0 {
		padTime, err = time.ParseDuration(allRootOptions.Padding)
		if err != nil {
			log.Fatalf("Error #29680: unable to convert to time.Duration: '%s, %s'\n", allRootOptions.Padding, err)
		}
	}
	maxTimeDiff = time.Duration(maxTimeDiff + padTime)
	if debugLevel > 98 {
		fmt.Println("maxTimeDiff:", maxTimeDiff)
		fmt.Println("===========================")
	}
	seenDates := make(map [time.Time]bool)
	for _, reqDate := range requiredDates {
		for i, csvDate := range csvDates {
			if debugLevel > 98 {
				fmt.Println()
				fmt.Println("csvDate:", csvDate.Format(dateOutputFmt))
				fmt.Println("reqDate:", reqDate.Format(dateOutputFmt))

			}
			if csvDate.IsSameOrBefore(&reqDate) {
				key := reqDate.ToTime()
				if csvDate.Format(dateOutputFmt) == "03/12/2021 6:40:01 PM Friday" {
					fmt.Println("dbg")
				}
				// compare the time duration difference
				diff := GetTimeDifference(csvDate,reqDate)
				if debugLevel > 98 {
					fmt.Println("diff:", diff, "  maxTimeDiff:", maxTimeDiff, "   diff < maxTimeDiff:", diff.Seconds() < maxTimeDiff.Seconds())
					//fmt.Println()
				}
				if diff.Seconds() <= maxTimeDiff.Seconds() {
					fmt.Println("seenDate:", key)
					seenDates[key] = true
					csvDates = RemoveSliceItem(csvDates,i)
					break // FIXME
				}
			}
		}
	}

	if debugLevel > 98 {
		fmt.Println()
		fmt.Println("seenDates")
		fmt.Println("============")
		var sortedSeen []string
		for key := range seenDates {
			sortedSeen = append(sortedSeen,key.String())
		}
		sort.Strings(sortedSeen)
		for _, seen := range sortedSeen {
			fmt.Println(seen)
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
		toCheck := reqDate.ToTime()  //FIXME
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
				//g, _ := goment.Unix(toCheck) //FIXME
				g, _ := goment.New(toCheck)
				fmt.Printf("missing date: %s\n", g.Format(dateOutputFmt))
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
