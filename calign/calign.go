// Copyright 2012 Milos Jakubicek
// Copyright 2017 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2017 Charles University, Faculty of Arts,
//                Institute of the Czech National Corpus
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

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
	"regexp"
	"sort"
	"strings"

	"github.com/czcorpus/ictools/attrib"
	"github.com/czcorpus/ictools/common"
	"github.com/czcorpus/ictools/mapping"
)

var (
	attrRegexp1 = regexp.MustCompile(".*xtargets=\"([^\"]+?)\".*")
	attrRegexp2 = regexp.MustCompile(".*xtargets='([^']+?)'.*")
)

// Processor represents an object used
// to process an alignment XML input file.
type Processor struct {
	attr1 attrib.GoPosAttr
	attr2 attrib.GoPosAttr
}

// NewProcessor creates a new instance of Processor
func NewProcessor(attr1 attrib.GoPosAttr, attr2 attrib.GoPosAttr) *Processor {
	return &Processor{
		attr1: attr1,
		attr2: attr2,
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

// processLine parses a single line of XML input file
// it parses lines of the form:
// <link type='1-1' xtargets='pl:_ACQUIS:jrc21959A1006_01:28:1;cs:_ACQUIS:jrc21959A1006_01:28:1' status='auto'/>
// any other xml element is ignored
func (p *Processor) processLine(line string, lineNum int) (mapping.Mapping, error) {
	x := attrRegexp2.FindStringSubmatch(line)
	if len(x) > 0 {
		aligned := strings.Split(x[1], ";")
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
		return mapping.Mapping{l1, l2}, nil
	}
	return mapping.Mapping{}, fmt.Errorf("Ignoring line: %d", lineNum)
}

// ProcessFile reads an input XML file containing mappings between
// structures (typically <s> for a sentence) of two languages and
// transforms them into a numeric representation based on internal
// identifiers used by Manatee.
// The function does not print anything to stdout.
func (p *Processor) ProcessFile(file *os.File, onItem func(item mapping.Mapping)) {
	reader := bufio.NewScanner(file)
	initialCap := common.FileSize(file.Name()) / 80. // TODO - estimation
	items := make([]mapping.Mapping, 0, initialCap)
	fromUndefItems := make([]mapping.Mapping, 0, initialCap/10)
	for i := 0; reader.Scan(); i++ {
		mp, err := p.processLine(reader.Text(), i)
		if err == nil {
			if mp.From.First > -1 {
				items = append(items, mp)

			} else {
				fromUndefItems = append(fromUndefItems, mp)
			}

		} else {
			log.Print(err)
		}
	}

	done := make(chan bool, 2)

	go func() {
		sort.Sort(mapping.SortableMapping(items))
		done <- true
	}()
	go func() {
		sort.Sort(mapping.SortableMapping(fromUndefItems))
		done <- true
	}()
	<-done
	<-done
	mapping.MergeMappings(items, fromUndefItems, onItem)
}
