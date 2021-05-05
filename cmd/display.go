package cmd

import (
	"fmt"
	"github.com/nleeper/goment"
	"github.com/olekukonko/tablewriter"
	"os"
)

func waitForInput() {
	fmt.Println()
	var pause string
	fmt.Scanln(&pause)
}

func DisplayTable(g []goment.Goment, desc string, pause bool, highlight int) {
	if len(g) == 0 {
		if pause {
			waitForInput()
		}
		return
	}
	var display [][]string
	for i, gom := range g {
		star := ""
		if i == highlight {
			star = "  * "
		}
		row := []string{ fmt.Sprintf("%s%d", star, i), gom.Format(dateOutputFmt)}
		display = append(display,row)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Index", fmt.Sprintf("Date - %s", desc)})
	table.SetAutoWrapText(false)
	table.AppendBulk(display)
	fmt.Println()
	table.Render()
	if pause {
		waitForInput()
	}
}
