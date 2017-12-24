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

package transalign

import (
	"bufio"
	"fmt"
	"github.com/czcorpus/ictools/common"
	"log"
	"os"
	"strings"
)

const (
	fileToCapacityRatio = 9
)

type PosRange struct {
	first int
	last  int
}

func (pr PosRange) LessThan(pr2 PosRange) bool {
	return pr.first < pr2.first || pr.first == pr2.first && pr.last < pr2.last
}

func (pr PosRange) IsEmpty() bool {
	return pr.first == -1 && pr.last == -1
}

func (pr PosRange) String() string {
	if pr.first == pr.last {
		return fmt.Sprintf("%d", pr.first)
	}
	return fmt.Sprintf("%d,%d", pr.first, pr.last)
}

func newPosRange(parsedVals []string) PosRange {
	return PosRange{common.Str2Int(parsedVals[0]), common.Str2Int(parsedVals[len(parsedVals)-1])}
}

// ensureContains checks and extends limit (if necessary)
// of range r1 to include r2. A possibly extended version
// of r1 is returned along with status whether the bounds
// have been changed
func ensureContains(r1 PosRange, r2 PosRange) (PosRange, bool) {
	changed := false
	if r2.first < r1.first {
		r1.first = r2.first
		changed = true
	}
	if r2.last > r1.last {
		r1.last = r2.last
		changed = true
	}
	return r1, changed
}

type PosRangeMap map[int]PosRange

func skipEmpty(idx int, final int, mapd PosRangeMap) int {
	var step int

	if idx < final {
		step = 1

	} else {
		step = -1
	}
	val := -1

	for idx != final && val == -1 {
		tmp, ok := mapd[idx]
		if !ok {
			continue
		}
		if idx < final {
			val = tmp.first

		} else {
			val = tmp.last
		}
		idx += step
	}

	return val
}

type HalfMapping struct {
	// source file
	file *os.File

	// source file reader
	reader *bufio.Scanner

	// maps ranges from L2 to pivot ranges
	mapToPivot map[int]PosRange

	// maps indices to position ranges
	ranges []PosRange

	// estimation of items number for efficient memory pre-allocation
	itemsEstim int
}

func NewHalfMapping(file *os.File) *HalfMapping {
	initialCap := common.FileSize(file.Name()) / fileToCapacityRatio
	return &HalfMapping{
		file:       file,
		reader:     bufio.NewScanner(file),
		mapToPivot: make(map[int]PosRange),
		itemsEstim: initialCap,
		ranges:     make([]PosRange, initialCap),
	}
}

func (hm *HalfMapping) RangesSize() int {
	return len(hm.ranges)
}

func (hm *HalfMapping) GetRange(idx int) PosRange {
	return hm.ranges[idx]
}

func (hm *HalfMapping) SetRange(idx int, v PosRange) {
	hm.ranges[idx] = v
}

func (hm *HalfMapping) HasRange(idx int) bool {
	return idx < len(hm.ranges) // TODO maybe we will need something more sophisticated
}

func (hm *HalfMapping) PivotMap() PosRangeMap {
	return hm.mapToPivot
}

func (hm *HalfMapping) Load() {
	var part int
	log.Printf("Loading %s ...", hm.file.Name())
	for hm.reader.Scan() {
		elms := strings.Split(hm.reader.Text(), "\t")
		// the mapping is (L2 -> L1/pivot)
		l1 := strings.Split(elms[1], ",")
		if l1[0] == "-1" {
			continue
		}
		l2 := strings.Split(elms[0], ",")
		l1Pair := newPosRange(l1)
		l2Pair := newPosRange(l2)
		for part = l1Pair.first; part <= l1Pair.last; part++ {
			hm.mapToPivot[part] = l2Pair
			if part >= len(hm.ranges) {
				hm.ranges = append(hm.ranges, make([]PosRange, hm.itemsEstim/2)...)
			}
			hm.ranges[part] = l1Pair
		}
	}
	hm.ranges = hm.ranges[:part+1]
	log.Printf("...done.")
}
