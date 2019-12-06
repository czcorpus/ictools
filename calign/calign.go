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
	"path/filepath"
	"strings"

	"github.com/czcorpus/ictools/mapping"
)

const (
	quoteStyleSingle = 1
	quoteStyleDouble = 2
)

// AttribMapper is a general type allowing transformation
// from a string value of a structural attribute
// (e.g. "cs:Adams-Holisticka_det_k:0:8:1") to a numeric
// form representing structure position in Manatee index.
type AttribMapper interface {
	Str2ID(value string) int
	ID2Str(ident int) string
}

// Processor represents an object used
// to process an alignment XML input file.
type Processor struct {
	attr1        AttribMapper
	attr2        AttribMapper
	valPrefix    string
	valSuffix    string
	valOffset    int
	lastPos      int
	lastPivotPos int
}

// NewProcessor creates a new instance of Processor
func NewProcessor(attr1 AttribMapper, attr2 AttribMapper, quoteStyle int) *Processor {
	valPrefix := "xtargets='"
	valSuffix := "'"
	if quoteStyle == quoteStyleDouble {
		valPrefix = "xtargets=\""
		valSuffix = "\""
	}

	return &Processor{
		attr1:        attr1,
		attr2:        attr2,
		valPrefix:    valPrefix,
		valSuffix:    valSuffix,
		valOffset:    len(valPrefix),
		lastPos:      0,
		lastPivotPos: 0,
	}
}

// processColElm parses a left or right item of a mapping line
func (p *Processor) processColElm(value string, attr AttribMapper, lineNum int) (mapping.PosRange, error) {
	if value == "" {
		return mapping.PosRange{-1, -1}, nil
	}
	elms := strings.Split(value, " ")
	beg, end := elms[0], elms[len(elms)-1]
	if beg == end {
		b := attr.Str2ID(beg)
		if b == -1 {
			return mapping.PosRange{}, fmt.Errorf("skipping invalid position [ %s ] on line %d", beg, lineNum+1)
		}
		return mapping.PosRange{b, b}, nil
	}
	b := attr.Str2ID(beg)
	e := attr.Str2ID(end)

	if b == -1 && e == -1 {
		return mapping.PosRange{}, fmt.Errorf("skipping invalid position range [ %s, %s ] on line %d", beg, end, lineNum+1)

	} else if b == -1 {
		log.Printf("ERROR: invalid left side of position range [ %s ] on line %d, using right side", beg, lineNum+1)
		return mapping.PosRange{e, e}, nil

	} else if e == -1 {
		log.Printf("ERROR: invalid right side of position range [ %s ] on line %d, using left side", end, lineNum+1)
		return mapping.PosRange{b, b}, nil
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
		endIdx := strings.Index(src[startIdx+p.valOffset:], p.valSuffix)
		if endIdx > -1 {
			return src[startIdx+p.valOffset : startIdx+p.valOffset+endIdx]
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
			return mapping.Mapping{}, fmt.Errorf("skipping invalid mapping on line %d", lineNum+1)
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
		if l2.Last > -1 {
			p.lastPivotPos = l2.Last
		}
		return mapping.Mapping{l1, l2, false}, nil
	}
	return mapping.Mapping{}, NewIgnorableError("skipping non-alignment line %d", lineNum)
}

// ProcessFile reads an input XML file containing mappings between
// structures (typically <s> for a sentence) of two languages and
// transforms them into a numeric representation based on internal
// identifiers used by Manatee.
// The function does not print anything to stdout.
func (p *Processor) ProcessFile(file *os.File, bufferSize int, onItem func(item mapping.Mapping, i int)) error {
	reader := bufio.NewScanner(file)
	reader.Buffer(make([]byte, bufio.MaxScanTokenSize), bufferSize)
	var i int
	count := 0
	for i = 0; reader.Scan(); i++ {
		if i%1000000 == 0 {
			log.Printf("INFO: Read %dm lines", i/1000000)
		}
		mp, err := p.processLine(reader.Text(), i)
		if err == nil {
			onItem(mp, count)
			count++

		} else {
			switch err.(type) {
			case IgnorableError:
				log.Print("INFO: ", err)
			default:
				log.Printf("ERROR: %s (file: %s)", err, filepath.Base(file.Name()))
			}
		}
	}
	err := reader.Err()
	if err != nil {
		return NewFileImportError(err, i)
	}
	return nil
}
