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

// Package transalign provides functions to generate
// mapping L1 -> L2 from L1 -> P and L2 -> P
package transalign

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/czcorpus/ictools/common"
	"github.com/czcorpus/ictools/mapping"
)

const (
	// a magical constant to estimate the size of
	// PivotMapping internal data slice
	fileToCapacityRatio = 9
)

// PosRangeMap maps data rows to PosRange values
// in other language
type PosRangeMap map[int]mapping.PosRange

// PivotMapping represents a list of mappings
// from a pivot language to a non-pivot one.
// The mapping is in expanded form which means
// that each data row with range r1, r2 (r2 > r1)
// is stored as r1, r1+1, ..., r2-1, r2 (but each
// line still knows the original range it belongs to)
type PivotMapping struct {
	// source file
	file *os.File

	// source file reader
	reader *bufio.Scanner

	// maps ranges from other language to pivot
	langToPivot map[int]*mapping.PosRange

	// gaps identifes all the rows which represent gaps between texts etc.
	// These rows cannot be used to expand translation range.
	gaps map[int]bool

	// list of lines mapping language ranges to pivot ranges (pivot language).
	ranges []*mapping.PosRange

	// estimation of items number for efficient memory pre-allocation
	itemsEstim int
}

// NewPivotMapping creates a new instance of PivotMapping
// and opens a file scanner for it. No data is loaded in
// this function (see PivotRange.Load()).
func NewPivotMapping(file *os.File) (*PivotMapping, error) {
	fSize, err := common.FileSize(file.Name())
	if err != nil {
		return nil, err
	}
	initialCap := fSize / fileToCapacityRatio
	return &PivotMapping{
		file:        file,
		reader:      bufio.NewScanner(file),
		langToPivot: make(map[int]*mapping.PosRange),
		itemsEstim:  initialCap,
		gaps:        make(map[int]bool),
		ranges:      make([]*mapping.PosRange, 0, initialCap),
	}, nil
}

// LangToPivot translates a data line of pivot lang. into a PosRange
// within the other lang.
func (hm *PivotMapping) LangToPivot(idx int) (*mapping.PosRange, bool) {
	ans, ok := hm.langToPivot[idx]
	return ans, ok
}

// Name returns a name of a source data file
func (hm *PivotMapping) Name() string {
	return hm.file.Name()
}

func (hm *PivotMapping) Size() int {
	return len(hm.ranges)
}

func (hm *PivotMapping) HasGapAtRow(idx int) bool {
	return hm.gaps[idx]
}

func (hm *PivotMapping) addMapping(langPair *mapping.PosRange) int {
	hm.ranges = append(hm.ranges, langPair)
	return len(hm.ranges) - 1
}

func (hm *PivotMapping) slicePivot(rightLimit int) {
	if rightLimit <= len(hm.ranges) {
		hm.ranges = hm.ranges[:rightLimit]

	} else {
		panic(fmt.Sprintf("Failed to slice ranges (pivot cap: %d, len: %d, idx: %d)",
			cap(hm.ranges), len(hm.ranges), rightLimit))
	}
}

// Load loads the respective data from a predefined file.
func (hm *PivotMapping) Load() {

	log.Printf("INFO: Loading %s ...", hm.file.Name())
	var i int
	for hm.reader.Scan() {
		elms := strings.Split(hm.reader.Text(), "\t")
		// the mapping in the file is (L1/L2 -> pivot)
		pivot := strings.Split(elms[1], ",")
		l2 := strings.Split(elms[0], ",")
		pivotPair, err1 := mapping.NewPosRange(pivot)
		if err1 != nil {
			log.Printf("ERROR: Failed to parse pivot on line %d: %s", i, err1)
		}
		l2Pair, err2 := mapping.NewPosRange(l2)
		if err2 != nil {
			log.Printf("ERROR: Failed to parse other lang on line %d: %s", i, err2)
		}
		i = hm.addMapping(&l2Pair)
		hm.langToPivot[i] = &pivotPair
		hm.gaps[i] = len(elms) == 3
	}
	log.Printf("INFO: Done (%d items).", len(hm.ranges))
}
