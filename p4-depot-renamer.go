package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	depotFrom                = flag.String("depot", "", "Name of the depot to rename")
	depotTo                  = flag.String("rename-to", "", "New name of the depot")
	checkpointFileName       = flag.String("cp", "", "Checkpoint file to process")
	checkpointOutputFileName = flag.String("o", "checkpoint.renamed", "Checkpoint file to write out")
)

func endOfFile(c chan<- JournalLine) {
	c <- JournalLine{EndOfFile: true}
}

func readCheckpoint(c chan<- JournalLine) {
	defer endOfFile(c)
	f, err := os.OpenFile(*checkpointFileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Fatalf("open file error: %v", err)
		return
	}
	defer f.Close()

	sc := bufio.NewReader(f)
	for {
		if jl, _ := ScanJournalLine(sc); !jl.EndOfFile {
			c <- jl
		} else {
			c <- jl
			break
		}

	}
}

func linePrinter(c <-chan JournalLine) {
	for line := range c {
		if line.EndOfFile {
			break
		}
		fmt.Println(line)
	}
}

func main() {
	flag.Parse()
	if *depotFrom == "" {
		log.Fatalf("-depot not specified")
		return
	}
	if *depotTo == "" {
		log.Fatalf("-rename-to not specified")
		return
	}
	if *checkpointFileName == "" {
		log.Fatalf("-cp not specified")
		return
	}
	if *checkpointOutputFileName == "" {
		log.Fatalf("-o not specified")
		return
	}

	readerOut := make(chan JournalLine)
	transOut := make(chan JournalLine)
	go readCheckpoint(readerOut)
	go transformer(readerOut, transOut, *depotFrom, *depotTo)
	linePrinter(transOut)

}
