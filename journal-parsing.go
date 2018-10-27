package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"
)

var (
	eof = byte(0)
)

//JournalElem individual journal elements
type JournalElem struct {
	EndOfLine    bool
	Encapsulated bool
	Data         string
}

//JournalLine stores a journal line for processing
type JournalLine struct {
	Idx       int
	Operator  string
	Version   string
	Table     string
	Raw       string
	Parsed    bool
	EndOfFile bool
	RowElems  []JournalElem
}

type lineScanner struct {
	Idx    int
	Reader *bufio.Reader
	Raw    bytes.Buffer
}

func (je *JournalElem) applyTransform(from string, to string) {
	re := regexp.MustCompile(from)
	je.Data = re.ReplaceAllString(je.Data, to)
}

func (ls *lineScanner) read() byte {
	r, err := ls.Reader.ReadByte()
	if err != nil {
		//probably an eof
		r = '\n'
	}
	ls.Raw.WriteByte(r)
	return r
}

func (ls *lineScanner) unread() {
	ls.Raw.Truncate(ls.Raw.Len() - 1)
	ls.Reader.UnreadByte()
}

func isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t'
}

func isEndLine(ch byte) bool {
	return ch == '\n'
}

func (ls *lineScanner) scanWhitespace() {
	for {
		if ch := ls.read(); ch == '\n' {
			ls.unread()
			break
		} else if !isWhitespace(ch) {
			ls.unread()
			break
		}
	}
}

func isDelimiter(ch byte) bool {
	return ch == '@'
}

func (ls *lineScanner) scanString() (token JournalElem) {
	var buf bytes.Buffer
	sanity := ls.read() //consume opening delimiter
	if !isDelimiter(sanity) {
		log.Fatal("scanString started without '@' delimiter")
	}
	for {
		if ch := ls.read(); ch == eof {
			break
		} else if isDelimiter(ch) {
			//possible end point, check next character in case this is "@@"
			nextCh := ls.read()
			if isDelimiter(nextCh) {
				//This is a "@@", write is as "@" in the data and keep going
				buf.WriteByte(nextCh)
			} else {
				//reached end delimiter
				ls.unread()
				break
			}
		} else {
			buf.WriteByte(ch)
		}
	}
	return JournalElem{Data: buf.String(), Encapsulated: true}
}

func (ls *lineScanner) scanNakedElem() (token JournalElem) {
	var buf bytes.Buffer
	for {
		if ch := ls.read(); ch == eof {
			break
		} else if isWhitespace(ch) || isEndLine(ch) {
			//reached end of element
			ls.unread()
			break
		} else {
			buf.WriteByte(ch)
		}
	}
	return JournalElem{Data: buf.String(), Encapsulated: false}
}

func (ls *lineScanner) scan() (token JournalElem) {
	ch := ls.read()

	//Ditch whitespace separators
	if isWhitespace(ch) {
		ls.unread()
		ls.scanWhitespace()
		ch = ls.read()
	}

	if isDelimiter(ch) {
		ls.unread()
		return ls.scanString()
	} else if isEndLine(ch) {
		return JournalElem{EndOfLine: true}
	} else {
		ls.unread()
		return ls.scanNakedElem()
	}
}

func (ls *lineScanner) scanToEndLine() {
	for {
		if ch := ls.read(); isEndLine(ch) {
			break
		}
	}
}

func trimEndLine(s string) string {
	return strings.TrimSuffix(s, "\n")
}

//ScanJournalLine parses a journal line into individual elements
func ScanJournalLine(r *bufio.Reader) (JournalLine, error) {
	//eg: @pv@ 1 @db.view@ @rpetti_client1@ 1 0 @//rpetti_client1/iperf-branches/...@ @//depot/iperf-branches/...@
	ls := lineScanner{Reader: r}
	var jl JournalLine

	//only support pv, dv, rv, vv
	je := ls.scan()
	if je.EndOfLine {
		return JournalLine{EndOfFile: true}, nil
	}
	if je.Data != "pv" &&
		je.Data != "dv" &&
		je.Data != "rv" &&
		je.Data != "vv" {
		//Unsupported, don't parse
		jl.Parsed = false
		ls.scanToEndLine()
		jl.Raw = trimEndLine(ls.Raw.String())
		return jl, nil
	}
	jl.Operator = je.Data
	je = ls.scan()
	jl.Version = je.Data
	je = ls.scan()
	jl.Table = je.Data

	for {
		je = ls.scan()
		if !je.EndOfLine {
			jl.RowElems = append(jl.RowElems, je)
		} else {
			break
		}
	}

	jl.Parsed = true
	return jl, nil
}

func (je JournalElem) String() string {
	if je.Encapsulated {
		var buf bytes.Buffer
		buf.WriteString("@")
		buf.WriteString(strings.Replace(je.Data, "@", "@@", -1))
		buf.WriteString("@")
		return buf.String()
	} else {
		return je.Data
	}
}

//Print JournalLine
func (jl JournalLine) String() string {
	if !jl.Parsed {
		return jl.Raw
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("@%s@ ", jl.Operator))
	buf.WriteString(fmt.Sprintf("%s ", jl.Version))
	buf.WriteString(fmt.Sprintf("@%s@ ", jl.Table))
	for _, je := range jl.RowElems {
		buf.WriteString(fmt.Sprintf("%s ", je))
	}
	return trimEndLine(buf.String())
}
