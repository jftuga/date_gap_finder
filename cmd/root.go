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
	"github.com/spf13/cobra"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

type rootOptions struct {
	Column int
	HasHeader bool
	Amount int
	Period string
	SkipWeekends bool
	SkipDays string
	Debug int
}

var allRootOptions rootOptions

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "date_gap_finder",
	Short: "searches for missing dates with in CSV files and optionally insert CSV entries for those missing dates",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var dateOutputFmt string = "L LTS dddd"
var debugLevel int

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.date_gap_finder.yaml)")
	rootCmd.PersistentFlags().IntVarP(&allRootOptions.Column, "column", "c", 0, "CSV column number (starts at zero)")
	rootCmd.PersistentFlags().BoolVarP(&allRootOptions.HasHeader, "header", "H", true, "if CSV file has header line")
	rootCmd.PersistentFlags().IntVarP(&allRootOptions.Amount, "amount", "a", -1, "a maximum, numeric duration")
	rootCmd.PersistentFlags().StringVarP(&allRootOptions.Period, "period", "p", "", "period of time, such as: days, hours, minutes")
	rootCmd.PersistentFlags().BoolVarP(&allRootOptions.SkipWeekends, "skipWeekends", "s", false, "allow gaps on weekends when set")
	rootCmd.PersistentFlags().StringVarP(&allRootOptions.SkipDays, "skipDays", "S", "", "skip comma-delimited set of fully spelled out days")
	rootCmd.PersistentFlags().IntVarP(&allRootOptions.Debug, "debug", "D", 0, "enable verbose debugging, set to 999 or 9999")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".date_gap_finder" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".date_gap_finder")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
