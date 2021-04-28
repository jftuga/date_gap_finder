# date_gap_finder

## Description
`date_gap_finder` searches for missing dates with in CSV files and optionally inserts new CSV entries for those missing dates.

Binaries for Windows, Mac, and Linux are available on the [releases](https://github.com/jftuga/date_gap_finder/releases) page.

## Motivation
I have many automated tasks that will append data to individual CSV files once per day.  However, if for some reason that 
task does not start, there is a possibility that I will not be notified of any errors.  Each of these tasks are monitored via a 
`Grafana Dashboard` using the [grafana-csv-datasource](https://github.com/marcusolsson/grafana-csv-datasource) data source
plugin.  This allows me to quickly view the status of all automated tasks.  When a task has missed a day, `date_gap_finder` will
insert CSV data in such a way that I will be notified that a task was not run.

This image displays the before and after of using `date_gap_finder`.  The **Photo Import** job runs every weekday.  Notice that `Monday, April 26` is missing from the *Before* image on the left.  This could be easily missed since it occurs on a Monday *(and the job isn't run on the weekends)*.  The *After* image on the right shows what can be displayed when using `date_gap_finder` to insert missing CSV data.  It is now much easier to detect a date gap.

![Grafana Before and After](dgf_before_after.png)

## Usage
```
date_gap_finder searches for missing dates with in CSV files and optionally inserts CSV entries for those missing dates.

Usage:
  date_gap_finder [command]

Available Commands:
  help        Help about any command
  insert      insert missing CSV entries
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


Flags for "insert":
  -R, --allRecords string    insert data to all columns of a missing row
  -O, --overwrite            overwrite existing CSV file; original file saved as .bak
  -r, --record stringArray   insert record with missing data

```

___

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
# Insert -1 at column 1 and 0 at column 2

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
