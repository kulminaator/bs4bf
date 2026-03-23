package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

var debugging = false

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	args := os.Args

	if len(args) < 4 {
		fatal("Usage: bs4bf filename range_start range_end pattern\n")
	}

	filename := args[1]
	start := args[2]
	end := args[3]
	pattern := args[4]

	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)

	debug("Started with %s %s %s %s\n", filename, start, end, pattern)

	if err != nil {
		fatal("Unable to open file" + err.Error())
	}

	findInFile(file, start, end, pattern)

	err2 := file.Close()

	if err2 != nil {
		fatal("Unable to close file" + err.Error())
	}
}

func findInFile(file *os.File, searchStart string, searchEnd string, pattern string) {
	fileSize := getFileSize(file)
	searchPosition := binarySearchFilePosition(file, fileSize, searchStart, searchEnd)
	scanLinesInRange(file, searchPosition, searchStart, searchEnd, pattern)
}

func getFileSize(file *os.File) int64 {
	fileStat, err := file.Stat()
	if err != nil {
		fatal("Unable to get file stats %s", err.Error())
	}
	return fileStat.Size()
}

func binarySearchFilePosition(file *os.File, fileSize int64, searchStart string, searchEnd string) int64 {
	seekSize := fileSize/2 + 1
	seekPosition := seekSize
	prefixLength := maxInt(len(searchStart), len(searchEnd))

	for seekSize > 1 {
		debug("Searching at %d with hop %d \n", seekPosition, seekSize)

		seekToLineStart(file, seekPosition)

		linePrefix := readLinePrefix(file, prefixLength)
		seekSize = seekSize / 2

		if linePrefix < searchStart {
			debug("Seeking forward from> %s \n", truncateString(linePrefix, 25))
			seekPosition = seekPosition + seekSize
		} else if linePrefix > searchEnd {
			debug("Seeking backwards from> %s \n", truncateString(linePrefix, 25))
			seekPosition = seekPosition - seekSize
		} else {
			// We found a line within our target range, we can stop here
			debug("Found position in range: %s \n", truncateString(linePrefix, 25))
			break
		}
	}

	return seekPosition
}

func maxInt(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func readLinePrefix(file *os.File, prefixLength int) string {
	buf := make([]byte, prefixLength)
	_, err := file.Read(buf)
	if err != nil {
		fatal("failed to read file %s: %s\n", file.Name(), err.Error())
	}
	return string(buf)
}

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[0:maxLen]
	}
	return s
}

func scanLinesInRange(file *os.File, startPosition int64, searchStart string, searchEnd string, pattern string) {
	seekToLineStart(file, startPosition)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		debug("Checking %s for %s", line, pattern)

		if line < searchStart {
			continue
		}
		if line > searchEnd {
			debug("Line is bigger than search end")
			break
		}
		if strings.Index(line, pattern) >= 0 {
			fmt.Println(line)
		}
	}
}

func seekToLineStart(file *os.File, currentPosition int64) {
	start := currentPosition
	var smallBuffer = make([]byte, 1)
	for start > 0 {
		_, seekErr := file.Seek(start-1, io.SeekStart)
		if seekErr != nil {
			fatal("Unable to seek to offset of file %s\n", file.Name())
		}
		_, readErr := file.Read(smallBuffer)
		if readErr != nil {
			fatal("failed to read file %s: %s\n", file.Name(), readErr.Error())
		}
		if smallBuffer[0] == '\n' {
			break
		}
		start--
	}
	_, seekErr := file.Seek(start, io.SeekStart)
	if seekErr != nil {
		fatal("Unable to seek to offset of file %s\n", file.Name())
	}
}

func fatal(format string, a ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, a...)
	os.Exit(1)
}

func debug(format string, a ...any) {
	if debugging {
		fmt.Printf("DEBUG:"+format, a...)
	}
}
