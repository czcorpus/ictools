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

	// maps ranges from pivot to the other language
	pivotToLang map[int]*mapping.PosRange

	// pivot maps indices (= data rows) to position ranges (pivot language).
	// PosRange with To > From is cloned across all lines it describes.
	// e.g. PosRange{3, 5} will exist in three instances on lines 3, 4 and 5.
	pivot []*mapping.PosRange

	// estimation of items number for efficient memory pre-allocation
	itemsEstim int

	minIdx int
}

// NewPivotMapping creates a new instance of PivotMapping
// and opens a file scanner for it. No data is loaded in
// this function (see PivotRange.Load()).
func NewPivotMapping(file *os.File) *PivotMapping {
	initialCap := common.FileSize(file.Name()) / fileToCapacityRatio
	return &PivotMapping{
		file:        file,
		reader:      bufio.NewScanner(file),
		pivotToLang: make(map[int]*mapping.PosRange),
		itemsEstim:  initialCap,
		pivot:       make([]*mapping.PosRange, initialCap),
		minIdx:      0,
	}
}

// PivotSize returns number of data lines in pivot language
func (hm *PivotMapping) PivotSize() int {
	return len(hm.pivot)
}

// GetPivotRange returns a range located on a specified data line
func (hm *PivotMapping) GetPivotRange(idx int) *mapping.PosRange {
	return hm.pivot[hm.index(idx)]
}

// SetPivotRange sets a range for a specified line
func (hm *PivotMapping) SetPivotRange(idx int, v *mapping.PosRange) {
	hm.pivot[hm.index(idx)] = v
}

// HasPivotRange tests whether there is a data line
// containing a PosRange.
func (hm *PivotMapping) HasPivotRange(idx int) bool {
	i := hm.index(idx)
	return idx >= hm.minIdx && i < len(hm.pivot) && hm.pivot[i] != nil
}

// PivotToLang translates a data line of pivot lang. into a PosRange
// within the other lang.
func (hm *PivotMapping) PivotToLang(idx int) (*mapping.PosRange, bool) {
	ans, ok := hm.pivotToLang[idx]
	return ans, ok
}

// Name returns a name of a source data file
func (hm *PivotMapping) Name() string {
	return hm.file.Name()
}

func (hm *PivotMapping) index(i int) int {
	return i - hm.minIdx
}

func (hm *PivotMapping) deindex(i int) int {
	return i + hm.minIdx
}

// Load loads the respective data from a predefined file.
func (hm *PivotMapping) Load() {
	var part int
	log.Printf("Loading %s ...", hm.file.Name())
	i := 0
	for hm.reader.Scan() {
		elms := strings.Split(hm.reader.Text(), "\t")
		// the mapping in the file is (L2 -> L1/pivot)
		pivot := strings.Split(elms[1], ",")
		if pivot[0] == "-1" {
			continue
		}
		l2 := strings.Split(elms[0], ",")
		pivotPair, err1 := mapping.NewPosRange(pivot)
		if err1 != nil {
			log.Printf("[WARNING] Failed to parse pivot on line %d: %s", i, err1)
		}
		l2Pair, err2 := mapping.NewPosRange(l2)
		if err2 != nil {
			log.Printf("[WARNING] Failed to parse other lang on line %d: %s", i, err2)
		}
		if i == 0 {
			hm.minIdx = pivotPair.First
		}
		for part = pivotPair.First; part <= pivotPair.Last; part++ {
			hm.pivotToLang[part] = &l2Pair
			if hm.index(part) >= len(hm.pivot) {
				hm.pivot = append(hm.pivot, make([]*mapping.PosRange, hm.index(part)-len(hm.pivot)+1)...)
			}
			hm.pivot[hm.index(part)] = &pivotPair
		}
		i++
	}
	hm.pivot = hm.pivot[:hm.index(part+1)]
	log.Printf("...done.")
}
