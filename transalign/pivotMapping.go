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
	"path/filepath"
	"strings"

	"github.com/czcorpus/ictools/common"
	"github.com/czcorpus/ictools/mapping"
)

const (
	// A magical constant to estimate number of lines
	// of a PivotMapping based on source file size.
	fileToCapacityRatio = 14
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

	// list of lines mapping language ranges to pivot ranges (pivot language).
	ranges []*mapping.PosRange

	// maps ranges from other language to pivot
	pivots []*mapping.PosRange

	// gaps identifes all the rows which represent gaps between texts etc.
	// These rows cannot be used to expand translation range.
	gaps map[int]bool

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
	log.Printf("INFO: pivot mapping size estimation for %s: %d",
		filepath.Base(file.Name()), initialCap)
	return &PivotMapping{
		file:       file,
		reader:     bufio.NewScanner(file),
		ranges:     make([]*mapping.PosRange, 0, initialCap),
		pivots:     make([]*mapping.PosRange, 0, initialCap),
		itemsEstim: initialCap,
		gaps:       make(map[int]bool),
	}, nil
}

// Size returns number of mapping definitions
// (= number of lines in a respective source file)
func (hm *PivotMapping) Size() int {
	return len(hm.ranges)
}

// HasGapAtRow tests whether the mapping [-1 -> X]
// or [X -> -1] (typically only the former variant
// should be present as we expect pivot to be superset
// of all the other language stuff) is due to a
// missing text/package etc.
// Because there is an important distinction here.
// In non-gap cases the transalign algorithm is allowed
// to extend the alignment across such an empty
// alignment if pivot's range includes one or more
// languages in the "left" language.
func (hm *PivotMapping) HasGapAtRow(idx int) bool {
	return hm.gaps[idx]
}

// Load loads the respective data from a predefined file.
func (hm *PivotMapping) Load() error {

	log.Printf("INFO: Loading %s ...", hm.file.Name())
	var i int
	for hm.reader.Scan() {
		elms := strings.Split(hm.reader.Text(), "\t")
		if elms[0] == mapping.ErrorMark {
			return fmt.Errorf("Refusing to continue due to the 'ERROR' mark in the source file")
		}
		// the mapping in the file is (SOME_LANG -> PIVOT_LANG)
		pivot := strings.Split(elms[1], ",")
		l2 := strings.Split(elms[0], ",")
		pivotPair, err1 := mapping.NewPosRange(pivot)
		if err1 != nil {
			return fmt.Errorf("ERROR: Failed to parse pivot on line %d: %s", i, err1)
		}
		l2Pair, err2 := mapping.NewPosRange(l2)
		if err2 != nil {
			return fmt.Errorf("ERROR: Failed to parse other lang on line %d: %s", i, err2)
		}

		hm.ranges = append(hm.ranges, &l2Pair)
		hm.pivots = append(hm.pivots, &pivotPair)
		i = len(hm.ranges) - 1
		hm.gaps[i] = len(elms) == 3
	}
	log.Printf("INFO: ...Done (%d items).", len(hm.ranges))
	return nil
}
