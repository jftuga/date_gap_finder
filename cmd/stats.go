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
	"github.com/hako/durafmt"
	"github.com/jftuga/date_gap_finder/fileOps"
	"github.com/spf13/cobra"
	"log"
	"strings"
	"time"
)

type dateCount struct {
	Date string
	count int
}

// statsCmd represents the stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show CSV file statistics",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		showAllStats(args)
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}

func showAllStats(args []string) {
	for _, fname := range args {
		statsData := statsOneFile(fname)
		DisplayStats(statsData)
	}
}

func buildRow(col1 string, col2 interface{}) []string {
	switch v := col2.(type) {
	case string:
		return []string {col1, v}
	case int:
		return []string{col1, fmt.Sprintf("%d", v)}
	}
	return []string {col1, "unknown data type"}
}

func statsOneFile(fname string) [][]string {
	debug := allRootOptions.Debug
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
		log.Fatalf("Error #45221: Unable to close CSV file: '%s'; %s\n", fname, err)
		return nil
	}

	row := 0
	numOfRecords := len(allRecords)
	if !allRootOptions.HasNoHeader {
		row = 1
		numOfRecords -= 1
	}
	layout := allRecords[row][allRootOptions.Column]
	numOfColumns := len(allRecords[row])

	var tableData [][]string
	tableData = append(tableData, buildRow("file",fname))
	tableData = append(tableData, buildRow("records",numOfRecords))
	tableData = append(tableData, buildRow("columns",numOfColumns))
	tableData = append(tableData, buildRow("date/time layout",layout))

	if numOfRecords <= 1 {
		return nil
	}

	// build a map containing a duration(string) and the frequency of that duration(int)
	// also get the average time between entries and total time
	var sum float64
	var count float64
	var a, b time.Time
	freq := make(map [string]int)
	for i := 2 ; i <= numOfRecords; i++ {
		a, err = dateparse.ParseAny(allRecords[i-1][allRootOptions.Column])
		b, err = dateparse.ParseAny(allRecords[i][allRootOptions.Column])

		if err != nil {
			log.Fatalf("Error #29099: %s\n", err)
		}
		diff := b.Sub(a)
		freq[diff.String()] += 1
		sum += diff.Seconds()
		count += 1
	}
	tmp := sum / count * 1000000000
	avg := time.Duration(tmp)
	dp := durafmt.Parse(avg)
	tableData = append(tableData, buildRow("average time between entries", LimitToSeconds(*dp)))

	first, _ := dateparse.ParseAny(allRecords[1][allRootOptions.Column])
	last, _ := dateparse.ParseAny(allRecords[numOfRecords][allRootOptions.Column])
	total := last.Sub(first)
	duration := durafmt.Parse(total)
	tableData = append(tableData, buildRow("total duration", LimitToSeconds(*duration)))

	// get the most frequently occurring duration
	var mode string
	modeCount := 0
	for key, val := range freq {
		if val > modeCount {
			mode = key
			modeCount = val
		}
	}

	modeDuration, _ := time.ParseDuration(mode)
	fullMode := fmt.Sprintf("%s (%s)", LimitToSeconds(*durafmt.Parse(modeDuration)), mode)
	tableData = append(tableData, buildRow("mode", fullMode))
	tableData = append(tableData, buildRow("mode count", modeCount))
	return tableData
}

func LimitToSeconds(d durafmt.Durafmt) string {
	var try string
	for n := 6; n > 0; n-- {
		try = d.LimitFirstN(n).String()
		if strings.Contains(try, "milli") || strings.Contains(try, "milli") {
			continue
		} else {
			break
		}
	}
	return try
}