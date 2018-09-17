// Copyright 2012 Milos Jakubicek
// Copyright 2017 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2017 Charles University, Faculty of Arts,
//                Institute of the Czech National Corpus
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// Package calign implements functions used to transform
// an alignment XML file to two-column numeric format
// Please note that Manatee library must be installed
// on the system (py, pl, java etc. wrappers are not needed).
package calign

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/czcorpus/ictools/attrib"
	"github.com/czcorpus/ictools/mapping"
)

// Processor represents an object used
// to process an alignment XML input file.
type Processor struct {
	attr1           attrib.GoPosAttr
	attr2           attrib.GoPosAttr
	valPrefix       string
	valSuffix       string
	lastPos         int
	lastPivotPos    int
	pivotStructSize int
}

// NewProcessor creates a new instance of Processor
func NewProcessor(attr1 attrib.GoPosAttr, attr2 attrib.GoPosAttr, quoteStyle int) *Processor {
	valPrefix := "xtargets='"
	valSuffix := "'"
	if quoteStyle == 2 {
		valPrefix = "xtargets=\""
		valSuffix = "\""
	}

	return &Processor{
		attr1:           attr1,
		attr2:           attr2,
		valPrefix:       valPrefix,
		valSuffix:       valSuffix,
		lastPos:         0,
		lastPivotPos:    0,
		pivotStructSize: attr2.Size(),
	}
}

// processColElm parses a left or right item of a mapping line
func (p *Processor) processColElm(value string, attr attrib.GoPosAttr, lineNum int) (mapping.PosRange, error) {
	if value == "" {
		return mapping.PosRange{-1, -1}, nil
	}
	elms := strings.Split(value, " ")
	beg, end := elms[0], elms[len(elms)-1]
	if beg == end {
		b := attr.Str2ID(beg)
		if b == -1 {
			return mapping.PosRange{}, fmt.Errorf("skipping invalid beg/end ('%s') on line %d", beg, lineNum+1)
		}
		return mapping.PosRange{b, b}, nil
	}
	b := attr.Str2ID(beg)
	e := attr.Str2ID(end)

	if b == -1 || e == -1 {
		if b == -1 && e == -1 {
			return mapping.PosRange{}, fmt.Errorf("skipping invalid beg, end ('%s','%s') on line %d", beg, end, lineNum+1)

		} else if b == -1 {
			log.Printf("invalid beg ('%s') on line %d, using end", beg, lineNum+1)
			return mapping.PosRange{e, e}, nil

		} else {
			log.Printf("invalid end ('%s') on line %d, using beg", end, lineNum+1)
			return mapping.PosRange{b, b}, nil
		}
	}
	return mapping.PosRange{b, e}, nil
}

// parseLine accepts lines of the form:
// <link type='1-1' xtargets='pl:_ACQUIS:jrc21959A1006_01:28:1;cs:_ACQUIS:jrc21959A1006_01:28:1' status='auto'/>
// other lines are ignored (i.e. an empty string is returned).
// Devel note: we try to avoid regexp here as it is quite slow compared with
// Python's or Perl's regexp engines (tested).

func (p *Processor) parseLine(src string) string {
	startIdx := strings.Index(src, p.valPrefix)
	if startIdx > -1 {
		endIdx := strings.Index(src[startIdx+10:], p.valSuffix)
		if endIdx > -1 {
			return src[startIdx+10 : startIdx+10+endIdx]
		}
	}
	return ""
}

// processLine parses a single line of XML input file
// any other xml element is ignored
func (p *Processor) processLine(line string, lineNum int) (mapping.Mapping, error) {
	srch := p.parseLine(line)
	if len(srch) > 0 {
		aligned := strings.Split(srch, ";")
		if len(aligned) > 2 {
			return mapping.Mapping{}, fmt.Errorf("Skipping invalid mapping on line %d", lineNum+1)
		}
		l1, err1 := p.processColElm(aligned[0], p.attr1, lineNum)
		if err1 != nil {
			return mapping.Mapping{}, err1
		}
		l2, err2 := p.processColElm(aligned[1], p.attr2, lineNum)
		if err2 != nil {
			return mapping.Mapping{}, err2
		}
		p.lastPos = l1.Last
		p.lastPivotPos = l2.Last
		return mapping.Mapping{l1, l2}, nil
	}
	return mapping.Mapping{}, fmt.Errorf("Ignoring line: %d", lineNum)
}

// ProcessFile reads an input XML file containing mappings between
// structures (typically <s> for a sentence) of two languages and
// transforms them into a numeric representation based on internal
// identifiers used by Manatee.
// The function does not print anything to stdout.
func (p *Processor) ProcessFile(file *os.File, bufferSize int, onItem func(item mapping.Mapping)) {
	reader := bufio.NewScanner(file)
	reader.Buffer(make([]byte, bufio.MaxScanTokenSize), bufferSize)
	for i := 0; reader.Scan(); i++ {
		if i%1000000 == 0 {
			log.Printf("Read %dm lines", i/1000000)
		}
		mp, err := p.processLine(reader.Text(), i)
		if err == nil {
			onItem(mp)

		} else {
			log.Print(err)
		}
	}
	if p.lastPos < p.pivotStructSize-1 {
		onItem(mapping.Mapping{
			From: mapping.NewEmptyPosRange(),
			To: mapping.PosRange{
				First: p.lastPivotPos + 1,
				Last:  p.pivotStructSize - 1,
			},
		})
	}
}
