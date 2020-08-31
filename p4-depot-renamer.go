package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	depotFrom                = flag.String("depot", "", "Name of the depot to rename")
	depotTo                  = flag.String("rename-to", "", "New name of the depot")
	checkpointFileName       = flag.String("cp", "", "Checkpoint file to process")
	checkpointOutputFileName = flag.String("o", "checkpoint.renamed", "Checkpoint file to write out")
	batchArgumentsFileName   = flag.String("batch", "", "Path to the batch arguments - use instead of depot/rename-to")
	forceOverwrite           = flag.Bool("f", false, "Overwrites the target file if exists")
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

	var batchArguments []BatchArgument
	var err error

	if *batchArgumentsFileName != "" {
		if *depotFrom != "" || *depotTo != "" {
			log.Fatalf("please do not specify -depot or -rename-to if -batch is specified")
			return
		}

		batchArguments, err = BatchArguments(*batchArgumentsFileName)
		if err != nil {
			log.Fatalf("error reading batch arguments: %v", err)
			return
		}
	}

	if *depotFrom == "" && len(batchArguments) == 0 {
		log.Fatalf("-depot not specified")
		return
	}
	if *depotTo == "" && len(batchArguments) == 0 {
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
	if strings.EqualFold(*checkpointFileName, *checkpointOutputFileName) {
		log.Fatalf("files specified by -cp and -o cannot be the same")
		return
	}

	if _, err := os.Stat(*checkpointOutputFileName); !os.IsNotExist(err) {
		if !(*forceOverwrite) {
			log.Fatalf("%s already exists, cannot overwrite", *checkpointOutputFileName)
			return
		}
	}

	if len(batchArguments) == 0 {
		batchArguments = append(batchArguments, BatchArgument{
			PathFrom:           *depotFrom,
			PathTo:             *depotTo,
			IncludedTransforms: []string{},
			ExcludedTransforms: []string{},
		})
	}

	log.Printf("will read and transform %s", *checkpointFileName)
	readerOut := make(chan JournalLine)
	transOut := make(chan JournalLine)
	go readCheckpoint(readerOut)
	go transformer(readerOut, transOut, batchArguments)
	lineWriter(transOut)

	log.Printf("transforms have been applied and saved to %s, you can now use it to restore your database", *checkpointOutputFileName)
}
