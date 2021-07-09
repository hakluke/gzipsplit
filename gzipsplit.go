package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"fmt"
	"log"
	"os"
)

type F struct {
	f  *os.File
	gf *gzip.Writer
	fw *bufio.Writer
}

func CreateGZ(s string) (f F) {

	fi, err := os.OpenFile(s, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Printf("Error in Create\n")
		panic(err)
	}
	gf := gzip.NewWriter(fi)
	fw := bufio.NewWriter(gf)
	f = F{fi, gf, fw}
	return
}

func WriteGZ(f F, s string) {
	(f.fw).WriteString(s)
}

func CloseGZ(f F) {
	f.fw.Flush()
	// Close the gzip first.
	f.gf.Close()
	f.f.Close()
}

func main() {
	buffer := flag.Int("b", 10000, "How many lines to write to each gzip file.")
	filePrefix := flag.String("f", "split", "filename prefix. Files will be in the format $prefix$number.gz")
	flag.Parse()

	s := bufio.NewScanner(os.Stdin)
	fileCounter := 1
	counter := 0
	lines := ""
	for s.Scan() {
		lines = lines + "\n" + s.Text()
		counter++
		if counter >= *buffer {
			f := CreateGZ(fmt.Sprintf("%s%d.gz", *filePrefix, fileCounter))
			fileCounter++
			WriteGZ(f, lines)
			lines = ""
			counter = 0
			CloseGZ(f)
		}
	}
	// write the final file
	f := CreateGZ(fmt.Sprintf("%s%d.gz", *filePrefix, fileCounter))
	WriteGZ(f, lines)
	CloseGZ(f)
}
