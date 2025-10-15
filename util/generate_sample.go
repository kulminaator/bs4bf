package main

import (
	"fmt"
	"io"
	"os"
)

// generates a big new sample file for testing. don't use outside development
func main() {
	fmt.Println("Generating sample file, this may take a while")
	file, err := os.OpenFile("../sample_files/bigfile1", os.O_WRONLY|os.O_CREATE, 0666)
	check(err)

	// 16 million lines, roughly 1.7gb
	for i := 0; i < 16*1000*1000; i++ {
		_, err := fmt.Fprintf(file,
			"Offset %020d content looking like a%020db and some other unrleated x%dy things\n",
			i, i, i%44)
		check(err)
	}

	_ = file.Close()
}

func check(err error) {
	if err != nil {
		_, _ = io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}
}
