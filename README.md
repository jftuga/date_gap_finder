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
  -h, --help               help for date_gap_finder
  -n, --noheader           set if CSV file does not have header line
  -p, --padding string     add time to range before considering a gap between two dates (default "1s")
  -S, --skipDays string    skip comma-delimited set of fully spelled out days
  -s, --skipWeekends       allow gaps on weekends when set
  -t, --tab                use tab character as CSV delimiter
  -u, --unit string        unit of time, such as: days, hours, minutes
  -v, --version            version for date_gap_finder

Use "date_gap_finder [command] --help" for more information about a command.


Flags for "insert":
  -R, --allRecords string    insert data to all columns of a missing row
  -O, --overwrite            overwrite existing CSV file; original file saved as .bak
  -r, --record stringArray   insert record with missing data; format: col#,value

```

___

## Example 1

This example file is called `e.csv`, which *should* get updated once per day. On the `15th`, it got updated twice. There are no entries for the `17th` and `18th`.

| Date | Errors | Warnings
|------|--------|----------
| 2021-04-15 06:55:01 | 0 | 23
| 2021-04-15 08:30:26 | 1 | 22
| 2021-04-16 06:55:01 | 0 | 23
| 2021-04-19 06:55:01 | 2 | 21

__

## search for date gaps
* Search for gaps occurring more than `1 day` apart

```
$ date_gap_finder search -a 1 -u days e.csv
2021-04-17 06:55:01
2021-04-18 06:55:01
```

## search for date gaps, but exclude weekends
* Use -s to skip weekends
```
PS C:\> .\date_gap_finder.exe search -a 1 -u days -s .\e.csv
(no results, because the 17th and 18th occur on a Saturday and Sunday)
```

## search for date gaps, but exclude on certain days
* Use -S to create a comma-separated list of days to skip
```
$ date_gap_finder search -a 1 -u days -S Sunday,Saturday e.csv
(no results, because the 17th and 18th occur on a Saturday and Sunday)

$ date_gap_finder search -a 1 -u days -S Sunday e.csv
2021-04-17 06:55:01
```

## insert rows for date gaps
* Adds column data for missing dates
* **CSV columns start 0**
* * `-r 1,1` - in column `1`, insert a value of `-1`
* * `-r 2,999` - in column `2`, insert a value of `999`
* * Note that multiple `-r` switches can be used
```
PS C:\> .\date_gap_finder.exe insert -a 1 -u days -r 1,-1 -r 2,999 .\e.csv
Date,Errors,Warnings
2021-04-15 06:55:01,0,23
2021-04-15 08:30:26,0,23
2021-04-16 06:55:01,0,23
2021-04-17 06:55:01,-1,999
2021-04-18 06:55:01,-1,999
2021-04-19 06:55:01,0,23
```

## insert rows for date gaps with the same value
* Similar to above, but use one `-R` switch instead of two `-r` switches
* * `-R 567` - excluding the `non-date` column, insert values of `567`
* * Note that `-R` and `-r` can be simultaneously used, with `-r` having precedence over `-R`
```
$ date_gap_finder insert -a 1 -u days -R 567 e.csv
Date,Errors,Warnings
2021-04-15 06:55:01,0,23
2021-04-15 08:30:26,0,23
2021-04-16 06:55:01,0,23
2021-04-17 06:55:01,567,567
2021-04-18 06:55:01,567,567
2021-04-19 06:55:01,0,23
```

## insert rows for date gaps and overwrite original file
* Use `-O` to overwrite the original file
* * a *date-versioned* backup file will be created containing the original file
* Use `-R` to add `0` to all missing columns of all missing rows
* Note that a backup file similar to this has been created: `e--20210506.162151.bak`
* * The naming convention is `filename--YYYYMMDD.HHMMSS.bak`
```
PS C:> .\date_gap_finder.exe insert -a 1 -u days -R 0 -O .\e.csv
PS C:> cat e.csv
Date,Errors,Warnings
2021-04-15 06:55:01,0,23
2021-04-15 08:30:26,0,23
2021-04-16 06:55:01,0,23
2021-04-17 06:55:01,0,0
2021-04-18 06:55:01,0,0
2021-04-19 06:55:01,0,23
```
___

## Example 2

This example file is called `f.csv`, which *should* get updated once per day. However, each entry is off by a few seconds.  You can use the `-p` switch to correct for this.  It will pad time before and after the correct time. It is missing `5` days: `12-14` and `17-18`.

| Date | Total | 
|------|-------|
| 2021-03-10 18:40:01 | 317
| 2021-03-11 18:40:01 | 249
| 2021-03-15 18:40:04 | 287
| 2021-03-16 18:40:03 | 320
| 2021-03-19 18:40:06 | 102

__

## search for date gaps and use time padding
* If a `6 second` time padding is not used, then `8` rows will be returned.  This is most likely not an accurate result since the times are only off by a few seconds. By using `-p 6s` an accurate list of missing rows is returned.
* Padding can end in `s` for seconds, `m` for minutes, or `h` for hours. `Days` are not supported.
* Note that `24 hours` is used instead of `1 day`.
* **It is better to use time padding (-p) vs. using a longer time amount (-a)**
```
$ date_gap_finder search -a 24 -u hours -p 6s f.csv
2021-03-12 18:40:01
2021-03-13 18:40:01
2021-03-14 18:40:01
2021-03-17 18:40:01
2021-03-18 18:40:01
```
___

## Example 3

This example file is called `g.csv`.  The `date` column is in `1` instead of the normal column `0`.  This file is also delimited by the `tab` character instead of the `comma` character. It has a date gap consisting of `04-15` through `04-18`.

| Processed | Date |
|-----------|------|
| 5125 | 2021-04-12
| 5197 | 2021-04-13
| 5206 | 2021-04-14
| 5222 | 2021-04-19
__

## search for gaps when date is not the first column and use a different column delimiter
* Use `-d "\t"` to define the column delimiter
* * A short cut for the `tab` character is to simply use `-t` instead of `-d "\t"`
* Use `-c` to denote column `1` instead of the default of column `0`
* Note that `1440 minutes` is equal to `1 day`
```
PS C:\> .\date_gap_finder.exe search -a 1440 -u minutes -d "\t" -c 1 .\g.csv
2021-04-15
2021-04-16
2021-04-17
2021-04-18
```

## insert data for date-gapped missing rows
* Similar switches to the above example, except that it is using the `insert` verb instead of `search`
* `-r 0,9999` - in column `0`, insert a value of `9990`
```
$ date_gap_finder insert -a 1440 -u minutes -t -c 1 -r 0,9999 g.csv
Processed	Date
5125	2021-04-12
5197	2021-04-13
5206	2021-04-14
9999	2021-04-15
9999	2021-04-16
9999	2021-04-17
9999	2021-04-18
5222	2021-04-19
```
___


## License
* [MIT License](LICENSE)

## Acknowledgments
* https://github.com/spf13/cobra
* https://github.com/nleeper/goment
* https://github.com/araddon/dateparse
* https://github.com/pivotal/go-ape
