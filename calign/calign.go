// Copyright 2012 Milos Jakubicek
// Copyright 2017 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2017 Charles University, Faculty of Arts,
//                Institute of the Czech National Corpus
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

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

type Processor struct {
	attr1 attrib.GoPosAttr
	attr2 attrib.GoPosAttr
}

func NewProcessor(attr1 attrib.GoPosAttr, attr2 attrib.GoPosAttr) *Processor {
	return &Processor{
		attr1: attr1,
		attr2: attr2,
	}
}

func (p *Processor) trLine(aligned string, attr attrib.GoPosAttr, lineNr int) (mapping.PosRange, error) {
	if aligned == "" {
		return mapping.PosRange{-1, -1}, nil
	}
	elms := strings.Split(aligned, " ")
	beg, end := elms[0], elms[len(elms)-1]
	if beg == end {
		b := attr.Str2ID(beg)
		if b == -1 {
			return mapping.PosRange{}, fmt.Errorf("skipping invalid beg/end ('%s') on line %d", beg, lineNr+1)
		}
		return mapping.PosRange{b, b}, nil
	}
	b := attr.Str2ID(beg)
	e := attr.Str2ID(end)

	if b == -1 || e == -1 {

		if b == -1 && e == -1 {
			return mapping.PosRange{}, fmt.Errorf("skipping invalid beg, end ('%s','%s') on line %d", beg, end, lineNr+1)

		} else if b == -1 {
			log.Printf("invalid beg ('%s') on line %d, using end", beg, lineNr+1)
			return mapping.PosRange{e, e}, nil

		} else {
			log.Printf("invalid end ('%s') on line %d, using beg", end, lineNr+1)
			return mapping.PosRange{b, b}, nil
		}
	}
	return mapping.PosRange{b, e}, nil
}

func (p *Processor) processLine(line string, lineNum int) (mapping.Mapping, error) {
	x := attrRegexp2.FindStringSubmatch(line)
	if len(x) > 0 {
		aligned := strings.Split(x[1], ";")
		if len(aligned) > 2 {
			return mapping.Mapping{}, fmt.Errorf("Skipping invalid mapping on line %d", lineNum+1)
		}
		l1, err1 := p.trLine(aligned[0], p.attr1, lineNum)
		if err1 != nil {
			return mapping.Mapping{}, err1
		}
		l2, err2 := p.trLine(aligned[1], p.attr2, lineNum)
		if err2 != nil {
			return mapping.Mapping{}, err2
		}
		return mapping.Mapping{l1, l2}, nil
	}
	return mapping.Mapping{}, fmt.Errorf("Ignoring line: %d", lineNum)
}

func (p *Processor) ProcessFile(file *os.File) {
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

	sort.Sort(mapping.SortableMapping(items))
	sort.Sort(mapping.SortableMapping(fromUndefItems))

	mapping.MergeMappings(items, fromUndefItems, func(item mapping.Mapping) {
		fmt.Printf("%s\t%s\n", item.From, item.To)
	})
}
