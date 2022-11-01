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
	"log"
	"sort"

	"github.com/czcorpus/ictools/mapping"
)

// fetchRow sets new language range and pivot range for provided langPos, pivotPos
// arguments using PivotMapping data on line langIdx.
// It returns true in case there was a range information available at the langIdx index.
// Otherwise (if langIdx > length of pm.ranges data), false is returned which
// means that the caller reached the end of data.
func fetchRow(
	langIdx int,
	langPos *mapping.PosRange,
	pivotPos *mapping.PosRange,
	pm *PivotMapping,
) bool {
	if langIdx >= len(pm.ranges) {
		return false
	}
	langPos.First = pm.ranges[langIdx].First
	langPos.Last = pm.ranges[langIdx].Last
	pivotPos.First = pm.pivots[langIdx].First
	pivotPos.Last = pm.pivots[langIdx].Last
	return true
}

// appendRow extends provided langPos, pivotPos using data loaded from line langIdx
func appendRow(langIdx int, langPos *mapping.PosRange, pivotPos *mapping.PosRange, pm *PivotMapping) {
	if langIdx >= len(pm.ranges) {
		return
	}
	if langPos.First == -1 {
		langPos.First = pm.ranges[langIdx].First
	}

	if pm.ranges[langIdx].Last != -1 {
		langPos.Last = pm.ranges[langIdx].Last
	}
	pivotPos.Last = pm.pivots[langIdx].Last
}

// addMapping is a simple wrapper around 'append' for the mapping
// pivot mapping slices which deliberately ignores -1 --> -1 mappings
// the algorithm sometimes produces.
func addMapping(list []mapping.Mapping, v mapping.Mapping) []mapping.Mapping {
	if v.From.First != -1 || v.To.First != -1 {
		return append(list, v)
	}
	return list
}

// Run implements an algorith for finding a mapping
// between L1 and L1 based on two "half mappings"
// L1 -> LP and L2 -> LP.
func Run(pivotMapping1 *PivotMapping, pivotMapping2 *PivotMapping, onItem func(mapping.Mapping)) {
	log.Print("INFO: Computing new alignment...")

	l1Idx := 0                  // current line in L1 source
	l1Pos := mapping.PosRange{} // current L1 range
	p1Pos := mapping.PosRange{} // current P1 range (pivot for L1)
	l2Idx := 0                  // current line in L2 source
	l2Pos := mapping.PosRange{} // current L2 range
	p2Pos := mapping.PosRange{} // current P2 range (pivot for L2)

	// We have to create two separate lists for the mappings as
	// one of the [-1, x], [x, -1] mappings must be kept separate
	// to be able to sort them. Final merging/sorting is done via
	// mapping.Iterator.
	mapL1L2 := make([]mapping.Mapping, 0, pivotMapping1.Size())      // TODO size estimation
	mapNoneL2 := make([]mapping.Mapping, 0, pivotMapping1.Size()/10) // 10 is just an estimate

	l1FetchOK := fetchRow(l1Idx, &l1Pos, &p1Pos, pivotMapping1)
	l2FetchOK := fetchRow(l2Idx, &l2Pos, &p2Pos, pivotMapping2)

	//for l1Idx < pivotMapping1.Size() || l2Idx < pivotMapping2.Size() {
	for l1FetchOK && l2FetchOK {
		if p1Pos.First < p2Pos.First { // must align beginning of pivots
			if p1Pos.Last == -1 {
				mapL1L2 = addMapping(mapL1L2, mapping.Mapping{
					From: l1Pos,
					To:   mapping.NewEmptyPosRange(),
				})
			}
			l1Idx++
			l1FetchOK = fetchRow(l1Idx, &l1Pos, &p1Pos, pivotMapping1)

		} else if p1Pos.First > p2Pos.First { // must align beginning of pivots
			if p2Pos.Last == -1 {
				mapNoneL2 = addMapping(mapNoneL2, mapping.Mapping{
					From: mapping.NewEmptyPosRange(),
					To:   l2Pos,
				})
			}
			l2Idx++
			l2FetchOK = fetchRow(l2Idx, &l2Pos, &p2Pos, pivotMapping2)

		} else { // pivots start at the same position; now try to align end positions
			if p1Pos.Last > p2Pos.Last {
				if pivotMapping1.HasGapAtRow(l1Idx) { // we cannot extend alignment across a gap
					mapNoneL2 = addMapping(mapNoneL2, mapping.Mapping{
						From: mapping.NewEmptyPosRange(),
						To:   l2Pos,
					})
					l2Idx++
					l2FetchOK = fetchRow(l2Idx, &l2Pos, &p2Pos, pivotMapping2)
					// a correction to keep pivots aligned (a spec. situation)
					// but we're losing compression here (TODO improve)
					p1Pos.First = p2Pos.First

				} else {
					l2Idx++
					appendRow(l2Idx, &l2Pos, &p2Pos, pivotMapping2)
				}

			} else if p2Pos.Last > p1Pos.Last {
				if pivotMapping2.HasGapAtRow(l2Idx) {
					mapL1L2 = addMapping(mapL1L2, mapping.Mapping{
						From: l1Pos,
						To:   mapping.NewEmptyPosRange(),
					})
					l1Idx++
					l1FetchOK = fetchRow(l1Idx, &l1Pos, &p1Pos, pivotMapping1)
					p2Pos.First = p1Pos.First // a correction to keep pivots aligned (a spec. situation)

				} else {
					l1Idx++
					appendRow(l1Idx, &l1Pos, &p1Pos, pivotMapping1)
				}

			} else if p1Pos.Last == -1 && p2Pos.Last == -1 {
				mapL1L2 = addMapping(mapL1L2, mapping.Mapping{
					From: l1Pos,
					To:   mapping.NewEmptyPosRange(),
				})
				mapNoneL2 = addMapping(mapNoneL2, mapping.Mapping{
					From: mapping.NewEmptyPosRange(),
					To:   l2Pos,
				})
				l1Idx++
				l1FetchOK = fetchRow(l1Idx, &l1Pos, &p1Pos, pivotMapping1)
				l2Idx++
				l2FetchOK = fetchRow(l2Idx, &l2Pos, &p2Pos, pivotMapping2)

			} else {
				if l1Pos.First != -1 {
					mapL1L2 = addMapping(mapL1L2, mapping.Mapping{
						From: l1Pos,
						To:   l2Pos,
					})

				} else {
					mapNoneL2 = addMapping(mapNoneL2, mapping.Mapping{
						From: l1Pos,
						To:   l2Pos,
					})
				}
				l1Idx++
				l1FetchOK = fetchRow(l1Idx, &l1Pos, &p1Pos, pivotMapping1)
				l2Idx++
				l2FetchOK = fetchRow(l2Idx, &l2Pos, &p2Pos, pivotMapping2)

			}
		}
	}

	log.Print("INFO: Sorting L1->L2/None and None->L2 lists...")
	done := make(chan bool, 2)
	go func() {
		sort.Sort(mapping.SortableMapping(mapL1L2))
		done <- true
	}()
	go func() {
		sort.Sort(mapping.SortableMapping(mapNoneL2))
		done <- true
	}()
	<-done
	<-done

	log.Print("INFO: Compressing and generating output...")

	mapping.MergeMappings(mapL1L2, mapNoneL2, onItem)

}
