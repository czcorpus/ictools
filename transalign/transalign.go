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

func fetchRow(langIdx int, langPos *mapping.PosRange, pivotPos *mapping.PosRange, pm *PivotMapping) bool {
	if langIdx >= len(pm.ranges) {
		return true // TODO !!!
	}
	langPos.First = pm.ranges[langIdx].First
	langPos.Last = pm.ranges[langIdx].Last
	pivot, ok := pm.LangToPivot(langIdx)
	if !ok {
		// TODO
	}
	pivotPos.First = pivot.First
	pivotPos.Last = pivot.Last
	return pm.HasGapAtRow(langIdx)
}

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
	pivot, ok := pm.LangToPivot(langIdx)
	if !ok {
		// TODO
	}
	pivotPos.Last = pivot.Last
}

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

	l1Idx := 0
	l1Pos := mapping.PosRange{}
	p1Pos := mapping.PosRange{}
	l2Idx := 0
	l2Pos := mapping.PosRange{}
	p2Pos := mapping.PosRange{}

	mapL1L2 := make([]mapping.Mapping, 0, pivotMapping1.Size()) // TODO size estimation
	// we have to keep one of [-1, x], [x, -1] mapping separate
	// because these two cannot be sorted together in a traditional way
	mapNoneL2 := make([]mapping.Mapping, 0, pivotMapping1.Size()/10) // 10 is just an estimate

	fetchRow(l1Idx, &l1Pos, &p1Pos, pivotMapping1)
	fetchRow(l2Idx, &l2Pos, &p2Pos, pivotMapping2)

	for l1Idx < len(pivotMapping1.ranges) && l2Idx < len(pivotMapping2.ranges) {
		if p1Pos.First < p2Pos.First { // must align beginning of pivots
			if p1Pos.Last == -1 {
				mapL1L2 = addMapping(mapL1L2, mapping.Mapping{
					From: l1Pos,
					To:   mapping.NewEmptyPosRange(),
				})
			}
			l1Idx++
			fetchRow(l1Idx, &l1Pos, &p1Pos, pivotMapping1)

		} else if p1Pos.First > p2Pos.First { // must align beginning of pivots
			if p2Pos.Last == -1 {
				mapNoneL2 = addMapping(mapNoneL2, mapping.Mapping{
					From: mapping.NewEmptyPosRange(),
					To:   l2Pos,
				})
			}
			l2Idx++
			fetchRow(l2Idx, &l2Pos, &p2Pos, pivotMapping2)

		} else {
			if p1Pos.Last > p2Pos.Last {
				if pivotMapping1.HasGapAtRow(l1Idx) { // we cannot extend alignment across a gap
					mapNoneL2 = addMapping(mapNoneL2, mapping.Mapping{
						From: mapping.NewEmptyPosRange(),
						To:   l2Pos,
					})
					l2Idx++
					fetchRow(l2Idx, &l2Pos, &p2Pos, pivotMapping2)
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
					fetchRow(l1Idx, &l1Pos, &p1Pos, pivotMapping1)
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
				l2Idx++
				fetchRow(l1Idx, &l1Pos, &p1Pos, pivotMapping1)
				fetchRow(l2Idx, &l2Pos, &p2Pos, pivotMapping2)

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
				l2Idx++
				fetchRow(l1Idx, &l1Pos, &p1Pos, pivotMapping1)
				fetchRow(l2Idx, &l2Pos, &p2Pos, pivotMapping2)

			}
		}
	}

	log.Print("INFO: Sorting L1-L2 and None->L2 lists...")
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
