# date_gap_finder

**WORK IN PROGRESS**

## Description
Searches for missing dates with in CSV files and optionally insert new CSV entries for those missing dates

## Motivation
I have many automated tasks that will add data to a CSV file each day.  However, if for some reason that task does not
run, there is a possibility that I will not be notified.  Each of these CSV tasks are monitored via a Grafana dashboard
using the [grafana-csv-datasource](https://github.com/marcusolsson/grafana-csv-datasource) data source plugin.  This
allows me to quickly the status of all tasks.  When a task has missed a day, this program can insert data letting me know
that it was not run.

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
  -a, --amount int        a maximum, numeric duration (default -1)
  -c, --column int        CSV column number (starts at zero)
  -D, --debug int         enable verbose debugging, set to 999 or 9999
  -H, --header            if CSV file has header line (default true)
  -h, --help              help for date_gap_finder
  -p, --period string     period of time, such as: days, hours, minutes
  -S, --skipDays string   skip comma-delimited set of fully spelled out days
  -s, --skipWeekends      allow gaps on weekends when set
  -v, --version           version for date_gap_finder

Use "date_gap_finder [command] --help" for more information about a command.

```

## Examples

## License
* [MIT License](LICENSE)

## Acknowledgments
* https://github.com/spf13/cobra
* https://github.com/nleeper/goment
* https://github.com/araddon/dateparse
* https://github.com/pivotal/go-ape
