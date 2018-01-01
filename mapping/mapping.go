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

package mapping

import (
	"fmt"
	"github.com/czcorpus/ictools/common"
	"log"
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
func NewPosRange(parsedVals []string) PosRange {
	return PosRange{common.Str2Int(parsedVals[0]), common.Str2Int(parsedVals[len(parsedVals)-1])}
}

// ----------------------------------------------

// Mapping represents a mapping between
// two structures from aligned corpora.
// These mappings are in general M:N
// (which is why we use PosRange internally here)
type Mapping struct {
	From PosRange
	To   PosRange
}

func (m Mapping) String() string {
	return fmt.Sprintf("%s\t%s", m.From, m.To)
}

func NewMapping(from1 int, from2 int, to1 int, to2 int) Mapping {
	return Mapping{
		PosRange{from1, from2},
		PosRange{to1, to2},
	}
}

// ----------------------------------------------

// SortableMapping implements sort.Interface
// for either [a, b] + [a, -1] or [-1, a] items
// (i.e. you cannot combine the two together
// as it is undefined how to compare [a, -1]
// and [-1, b]).
type SortableMapping []Mapping

func (sm SortableMapping) Len() int {
	return len(sm)
}

func (sm SortableMapping) Swap(i, j int) {
	sm[i], sm[j] = sm[j], sm[i]
}

func (sm SortableMapping) Less(i, j int) bool {
	if sm[i].From.First > -1 && sm[j].From.First > -1 {
		return sm[i].From.First < sm[j].From.First

	} else if sm[i].From.First == -1 || sm[j].From.First == -1 {
		return sm[i].To.First < sm[j].To.First

	} else if sm[i].To.First == -1 && sm[j].To.First == -1 {
		return sm[i].From.First < sm[j].From.First
	}
	panic("unknow type combination")
}

// ----------------------------------------------

// Iterator is used to merge two sorted
// mappings together.
type Iterator struct {
	mapping []Mapping
	currIdx int
	final   Mapping
}

// NewIterator creates a new Iterator instance
func NewIterator(data []Mapping, final Mapping) Iterator {
	return Iterator{
		mapping: data,
		currIdx: 0,
		final:   final,
	}
}

// Next returns the next mapping item. In
// case the last item is reached, Next returns
// a special 'final' item which is manually
// set when instantiating the iterator.
func (m *Iterator) Next() Mapping {
	defer func() {
		m.currIdx++
	}()
	if m.currIdx < len(m.mapping) {
		return m.mapping[m.currIdx]
	}
	return m.final
}

// ----------------------------------------------

func max(v1 int, v2 int) int {
	if v1 > v2 {
		return v1
	}
	return v2
}

func calcFinalItem(mainMapping []Mapping, backEmptyMapping []Mapping) Mapping {
	ifinal := -1
	if len(mainMapping) > 0 {
		ifinal = max(mainMapping[len(mainMapping)-1].From.Last, mainMapping[len(mainMapping)-1].To.Last)
	}
	if len(backEmptyMapping) > 0 {
		ifinal = max(ifinal, backEmptyMapping[len(backEmptyMapping)-1].To.Last)
	}

	return Mapping{
		PosRange{ifinal, ifinal},
		PosRange{ifinal, ifinal},
	}
}

// MergeMappings merges two sorted mappings, one containing items
// [a, b], [a, -1] (where a, b > -1) and one
// containing items [-1, a] (where a > -1) into a single
// sorted one. The sorting rule is the following:
// the [-1, a1] mapping gets priority over [a2, b], [a3, -1]
// if a1 < b (a1 < -1 cannot happen).
// The function does not create a new slice for the merged
// items. It's up to a function user to provide a function
// specifying what to do with each item.
func MergeMappings(mainMapping []Mapping, mapFromEmpty []Mapping, onItem func(item Mapping)) {
	finalMapping := calcFinalItem(mainMapping, mapFromEmpty)
	iterL2L3 := NewIterator(mainMapping, finalMapping)
	iterL3 := NewIterator(mapFromEmpty, finalMapping)

	vL2L3 := iterL2L3.Next()
	vL3 := iterL3.Next()
	var curr Mapping
	for vL2L3 != finalMapping || vL3 != finalMapping {
		if vL3.To.LessThan(vL2L3.To) {
			curr = vL3
			vL3 = iterL3.Next()

		} else {
			curr = vL2L3
			vL2L3 = iterL2L3.Next()
		}
		onItem(curr)
	}
	log.Print("...done.")
}
