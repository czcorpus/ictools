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
	"github.com/czcorpus/ictools/common"
	"github.com/czcorpus/ictools/mapping"
	"log"
	"os"
	"strings"
)

const (
	fileToCapacityRatio = 9
)

type PosRangeMap map[int]mapping.PosRange

// PivotMapping represents a list of mappings
// from a pivot language to a non-pivot one.
type PivotMapping struct {
	// source file
	file *os.File

	// source file reader
	reader *bufio.Scanner

	// maps ranges from pivot to the other language
	pivotToLang map[int]mapping.PosRange

	// maps indices to position ranges (pivot language)
	pivot []mapping.PosRange

	// estimation of items number for efficient memory pre-allocation
	itemsEstim int
}

// NewPivotMapping creates a new instance of PivotMapping
// and opens a file scanner for it. No data is loaded in
// this function (see PivotRange.Load()).
func NewPivotMapping(file *os.File) *PivotMapping {
	initialCap := common.FileSize(file.Name()) / fileToCapacityRatio
	return &PivotMapping{
		file:        file,
		reader:      bufio.NewScanner(file),
		pivotToLang: make(map[int]mapping.PosRange),
		itemsEstim:  initialCap,
		pivot:       make([]mapping.PosRange, initialCap),
	}
}

func (hm *PivotMapping) PivotSize() int {
	return len(hm.pivot)
}

func (hm *PivotMapping) GetPivotRange(idx int) mapping.PosRange {
	return hm.pivot[idx]
}

func (hm *PivotMapping) SetPivotRange(idx int, v mapping.PosRange) {
	hm.pivot[idx] = v
}

func (hm *PivotMapping) HasPivotRange(idx int) bool {
	return idx < len(hm.pivot)
}

func (hm *PivotMapping) PivotToLang(idx int) (mapping.PosRange, bool) {
	ans, ok := hm.pivotToLang[idx]
	return ans, ok
}

// Load loads the respective data from a predefined file.
func (hm *PivotMapping) Load() {
	var part int
	log.Printf("Loading %s ...", hm.file.Name())
	for hm.reader.Scan() {
		elms := strings.Split(hm.reader.Text(), "\t")
		// the mapping in the file is (L2 -> L1/pivot)
		l1 := strings.Split(elms[1], ",")
		if l1[0] == "-1" {
			continue
		}
		l2 := strings.Split(elms[0], ",")
		l1Pair := mapping.NewPosRange(l1)
		l2Pair := mapping.NewPosRange(l2)
		for part = l1Pair.First; part <= l1Pair.Last; part++ {
			hm.pivotToLang[part] = l2Pair
			if part >= len(hm.pivot) {
				hm.pivot = append(hm.pivot, make([]mapping.PosRange, hm.itemsEstim/2)...)
			}
			hm.pivot[part] = l1Pair
		}
	}
	hm.pivot = hm.pivot[:part+1]
	log.Printf("...done.")
}
