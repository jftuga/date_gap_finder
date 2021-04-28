# date_gap_finder

**WORK IN PROGRESS**

## Description
Searches for missing dates with in CSV files and optionally insert new CSV entries for those missing dates

## Motivation
I have many automated tasks that will add data to CSV files (one per task) each day.  However, if for some reason that 
task does not run, there is a possibility that I will not be notified.  Each of these CSV tasks are monitored via a 
Grafana dashboard using the [grafana-csv-datasource](https://github.com/marcusolsson/grafana-csv-datasource) data source
plugin.  This allows me to quickly view the status of all automated tasks.  When a task has missed a day, then this 
program will insert CSV data in such a way that I will be notified that a task was not run.

## Usage
```
searches for missing dates with in CSV files and optionally insert CSV entries for those missing dates

Usage:
  date_gap_finder [command]

Available Commands:
  help        Help about any command
  insert      insert missing CSV entries
  replace     A brief description of your command
  search      search CSV files for missing dates

Flags:
  -a, --amount int         a maximum, numeric duration (default -1)
  -c, --column int         CSV column number (starts at zero)
  -D, --debug int          enable verbose debugging, set to 999 or 9999
  -d, --delimiter string   CSV delimiter (default ",")
  -H, --header             if CSV file has header line (default true)
  -h, --help               help for date_gap_finder
  -p, --period string      period of time, such as: days, hours, minutes
  -S, --skipDays string    skip comma-delimited set of fully spelled out days
  -s, --skipWeekends       allow gaps on weekends when set
  -v, --version            version for date_gap_finder

Use "date_gap_finder [command] --help" for more information about a command.

```

## Example: Search For Missing Dates

This example file is called `e.csv`, which *should* get updated once per day.

| Date | Errors | Warnings
|------|--------|----------
| 2021-04-15 06:55:01 | 0 | 23
| 2021-04-15 08:30:26 | 1 | 22
| 2021-04-16 06:55:01 | 0 | 23
| 2021-04-19 06:55:01 | 2 | 21

```
# I usually allow for a couple of extra minutes in case a process runs a little longer than usual. Therefore, 1442 instead of 1440.
$ date_gap_finder search -a 1442 -p minutes e.csv
2021-04-17 06:59:01
2021-04-18 07:01:01

# For the search verb, the program's exit code will be equal to the number of missed dates
# For Powershell, use $LASTEXITCODE
$ echo $?
2

# Skip weekends with -s
$ date_gap_finder search -a 1442 -p -s minutes e.csv
(no output - as all missed dates occurs on either a Saturday or Sunday)

# This also skips weekends, but you could also include other days of the week
$ date_gap_finder search -a 1442 -p -S Saturday,Sunday minutes e.csv
(no output - as all missed dates occurs on either a Saturday or Sunday)

```

## Example: Insert Records For Missing Dates

```

# Columns numbers start at zero.
# Insert -1 at column 1 and 1 at column 2

PS C:\> .\date_gap_finder.exe insert -a 1442 -p minutes -r 1,-1 -r 2,0 .\e.csv
Date,Errors,Warnings
2021-04-15 06:55:01,0,23
2021-04-15 08:30:26,0,23
2021-04-16 06:55:01,0,23
2021-04-17 06:59:01,-1,0
2021-04-18 07:01:01,-1,0
2021-04-19 06:55:01,0,23

# Use -R to set all missing columns and also skip Sundays
PS C:\> .\date_gap_finder.exe insert -a 1442 -p minutes -R 999 -S Sunday .\e.csv
Date,Errors,Warnings
2021-04-15 06:55:01,0,23
2021-04-15 08:30:26,0,23
2021-04-16 06:55:01,0,23
2021-04-17 06:59:01,999,999
2021-04-19 06:55:01,0,23

```

___


## License
* [MIT License](LICENSE)

## Acknowledgments
* https://github.com/spf13/cobra
* https://github.com/nleeper/goment
* https://github.com/araddon/dateparse
* https://github.com/pivotal/go-ape
