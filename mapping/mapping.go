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

// Package mapping provides data types and functions used
// to manipulate numeric mapping between two aligned structures
package mapping

import (
	"fmt"
	"strings"

	"github.com/czcorpus/ictools/common"
)

const (
	// ErrorMark is written into a numeric LANG-PIVOT_LANG
	// mapping file (created during import) in case of an error.
	// This prevents transalign task from running when such
	// a file is loaded.
	ErrorMark = "ERROR"

	errorPositionValue = -2
)

// PosRange defines a range of (Manatee) structure
// positions. The most typical range encountered
// in data is of size 1. In such case First == Last.
type PosRange struct {
	First int
	Last  int
}

// LessThan defines an ordering for range
// items. First position is most significant
// and Last position is least significant.
func (pr PosRange) LessThan(pr2 PosRange) bool {
	return pr.First < pr2.First || pr.First == pr2.First && pr.Last < pr2.Last
}

// String converts the range into a
// format required by other applications.
func (pr PosRange) String() string {
	if pr.First == pr.Last {
		return fmt.Sprintf("%d", pr.First)
	}
	return fmt.Sprintf("%d,%d", pr.First, pr.Last)
}

// NewPosRange creates a new PosRange from
// a list of string-encoded integers.
// An empty string is treated as -1 (i.e. 'undefined')
func NewPosRange(parsedVals []string) (PosRange, error) {
	v1 := common.Str2Int(parsedVals[0])
	v2 := common.Str2Int(parsedVals[len(parsedVals)-1])
	if (v1 == -1 || v2 == -1) && v1 != v2 {
		return PosRange{}, fmt.Errorf("Cannot use -1 with a different value: %d, %d", v1, v2)
	}
	if v1 > v2 {
		return PosRange{}, fmt.Errorf("First value must be smaller than the last one. Got: %d, %d", v1, v2)
	}
	return PosRange{v1, v2}, nil
}

// NewEmptyPosRange creates a range "-1,-1"
func NewEmptyPosRange() PosRange {
	return PosRange{First: -1, Last: -1}
}

// ----------------------------------------------

// Mapping represents a mapping between two structures from aligned corpora.
// These mappings are in general M:N (which is why we use PosRange internally
// here).
// Besides positions 0,...,N the code here uses also -1 for undefined mapping
// and -2 for an error record (but the error record is exported into a special
// string when creating the outout).
type Mapping struct {
	From  PosRange
	To    PosRange
	IsGap bool
}

func (m Mapping) String() string {
	if m.IsGap {
		return fmt.Sprintf("%s\t%s\tg", m.From, m.To)

	} else if m.From.First == errorPositionValue {
		return fmt.Sprintf(ErrorMark)
	}
	return fmt.Sprintf("%s\t%s", m.From, m.To)
}

func (m *Mapping) IsEmpty() bool {
	return m.From.First == -1 && m.To.First == -1
}

// NewMapping creates a new instance of Mapping.
// The arguments can be understood as follows:
// from1,from2[TAB]to1,to2
func NewMapping(from1 int, from2 int, to1 int, to2 int) Mapping {
	return Mapping{
		PosRange{from1, from2},
		PosRange{to1, to2},
		false,
	}
}

// NewGapMapping creates a new instance of Mapping
// with isGap set to true.
// The arguments can be understood as follows:
// from1,from2[TAB]to1,to2
func NewGapMapping(from1 int, from2 int, to1 int, to2 int) Mapping {
	return Mapping{
		PosRange{from1, from2},
		PosRange{to1, to2},
		true,
	}
}

// NewMappingFromString creates a new Mapping instance
// from a two-column numeric source code line used as
// an intermediate format.
func NewMappingFromString(src string) (Mapping, error) {
	items := strings.Split(src, "\t")
	if len(items) < 2 {
		return Mapping{}, fmt.Errorf("No TAB separated data found")
	}
	l1t := strings.Split(items[0], ",")
	l2t := strings.Split(items[1], ",")
	r1, err1 := NewPosRange(l1t)
	if err1 != nil {
		return Mapping{}, err1
	}
	r2, err2 := NewPosRange(l2t)
	if err2 != nil {
		return Mapping{}, err2
	}
	return Mapping{r1, r2, len(items) == 3}, nil
}

// NewErrorMapping creates a mapping with all the
// values set to the value of constant errorPositionValue which
// is then transformed into an error string mark when exported
// to string.
func NewErrorMapping() Mapping {
	return Mapping{
		From: PosRange{errorPositionValue, errorPositionValue},
		To:   PosRange{errorPositionValue, errorPositionValue},
	}
}

// ----------------------------------------------

// SortableMapping implements sort.Interface
// for either [a, b] + [a, -1] (type A)
// or [-1, b] (type B) items
// (i.e. you cannot combine the two together
// as it is undefined how to compare [a, -1]
// and [-1, b]).
// Please note that this sorting is not able
// to process files correctly in case they
// contain different value mixing than the one
// described above (i.e. it may finish without
// an error but the result won't be a properly
// sorted mapping list)
type SortableMapping []Mapping

func (sm SortableMapping) Len() int {
	return len(sm)
}

func (sm SortableMapping) Swap(i, j int) {
	sm[i], sm[j] = sm[j], sm[i]
}

// Less compares items from either A iterator or B iterator
// TODO simplify this - it's not necessary to implement it this way
func (sm SortableMapping) Less(i, j int) bool {
	if sm[i].From.First > -1 && sm[j].From.First > -1 {
		return sm[i].From.First < sm[j].From.First

	} else if sm[i].From.First == -1 && sm[j].From.First == -1 {
		return sm[i].To.First < sm[j].To.First
	}
	panic(fmt.Sprintf("Unknow mapping type combination: [%s] [%s]", sm[i], sm[j]))
}

// ----------------------------------------------

// MergeMappings merges two sorted mappings, one containing items
// [a, b], [a, -1] (where a, b > -1) and one
// containing items [-1, a] (where a > -1) into a single
// sorted one. The sorting rule is the following:
// the [-1, a1] mapping gets priority over [a2, b], [a3, -1]
// if a1 < b (a1 < -1 cannot happen).
// The function does not create a new slice for the merged
// items. It's up to a function user to provide a function
// specifying what to do with each item.
//
// The onItem argument is a function applied on each encountered
// item (no mattter how data are chunked). The 'pos' argument
// of the function provides last position for each language column.
// This can be used to test (and resolve) possible gaps between
// items (e.g. Manatee *mkalign* does not like them).
func MergeMappings(mainMapping []Mapping, mapFromEmpty []Mapping, onItem func(item Mapping)) {
	iterL2L3 := NewIterator(mainMapping)
	iterL3 := NewIterator(mapFromEmpty)

	for iterL2L3.Unfinished() || iterL3.Unfinished() {
		if iterL3.HasPriorityOver(&iterL2L3) || !iterL2L3.Unfinished() {
			iterL3.Apply(onItem)
			iterL3.Next()

		} else {
			iterL2L3.Apply(onItem)
			iterL2L3.Next()
		}
	}
}
