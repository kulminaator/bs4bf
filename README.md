# bs4bf - binary search for big files

---
## Intro
If you have a big file which has records in a predictable order and want to search for a pattern in it, 
then this tool is for you. An example would be a syslog file or a postgresql csv log file.

The tool does a binary search to find the row starting with "search_start" 
and matches lines starting from there until it meets the "search_end".

If the files that you need to search are hundreds of GBs, this tool will be very useful 
and way faster than a regular grep. You can also combine it with grep by matching everything with "" instead 
of the search pattern.

## Building the tool

`go build`


---
## Usage:

```
bs4bf file search_start search_end pattern
```

example:
```
./bs4bf /var/log/syslog "Oct 12 07:05:02" "Oct 12 09:05:02" "ERROR"
```

This should find lines after 07:05 until 09:05 and find the ERROR statements from the range, ignoring the rest of the file.

## Comments

Word of warning. I'm not a professional go programmer. I just did this as an exercise and to rehearse the language.


