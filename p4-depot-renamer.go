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

func lineWriter(c <-chan JournalLine) {
	f, err := os.Create(*checkpointOutputFileName)
	if err != nil {
		log.Fatalf("error opening output file: %v", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for line := range c {
		if line.EndOfFile {
			break
		}
		w.WriteString(fmt.Sprintln(line))
	}
	w.Flush()
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
	if _, err := os.Stat(*checkpointOutputFileName); !os.IsNotExist(err) {
		log.Fatalf("%s already exists, cannot overwrite", *checkpointOutputFileName)
		return
	}

	log.Printf("will read and transform %s", *checkpointFileName)
	log.Printf("renaming depot \"%s\" to \"%s\"...", *depotFrom, *depotTo)
	readerOut := make(chan JournalLine)
	transOut := make(chan JournalLine)
	go readCheckpoint(readerOut)
	go transformer(readerOut, transOut, *depotFrom, *depotTo)
	lineWriter(transOut)

	log.Printf("transforms have been applied and saved to %s, you can now use it to restore your database", *checkpointOutputFileName)
}
