// Copyright 2012 Milos Jakubicek
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
	"fmt"
	"log"
	"sort"

	"github.com/czcorpus/ictools/mapping"
)

// ---------------------------

// skipEmpty searches for a valid non-empty (!= -1) non-pivot language
//  range corresponding matching to provided pivot range [idx, final]
// (or [final, idx] with step -1).
//
// The function returns both non-pivot lang. range and matching pivot range
// (which may have different 'final' value).
// The latter value is used to detect empty intersections.
func skipEmpty(idx int, final int, hMapping *PivotMapping) (int, int) {
	var step int

	if idx < final {
		step = 1

	} else {
		step = -1
	}
	val := -1

	for idx != final && val == -1 {
		tmp, ok := hMapping.PivotToLang(idx)
		if ok {
			if idx < final {
				val = tmp.First

			} else {
				val = tmp.Last
			}
		}
		idx += step
	}
	return val, idx - step
}

// expandToAlign checks and extends limit (if necessary)
// of range r1 to include r2. Status whether the bounds
// have been changed is returned.
func expandToAlign(r1 *mapping.PosRange, r2 *mapping.PosRange) bool {
	changed := false

	if r2.First < r1.First {
		r1.First = r2.First
		changed = true
	}
	if r2.Last > r1.Last {
		r1.Last = r2.Last
		changed = true
	}
	return changed
}

// Run implements an algorith for finding a mapping
// between L2 and L3 based on two "half mappings"
// L1 -> L2 and L1 -> L3.
func Run(pivotMapping1 *PivotMapping, pivotMapping2 *PivotMapping) {
	next := 0
	mapL2L3 := make([]mapping.Mapping, 0, pivotMapping1.PivotSize()) // TODO size estimation
	// we have to keep one of [-1, x], [x, -1] mapping separate
	// because these two cannot be sorted together in a traditional way
	mapEmptyL3 := make([]mapping.Mapping, 0, pivotMapping1.PivotSize()/10) // 10 is just an estimate
	log.Print("INFO: Computing new alignment...")

	var i int
	var extRng *mapping.PosRange

	for ix, rng := range pivotMapping1.pivot {
		i = pivotMapping1.deindex(ix)

		if i < next || rng == nil {
			continue
		}
		extRng = &mapping.PosRange{
			First: rng.First,
			Last:  rng.Last,
		}

		changed := true
		for changed {
			changed = false
			if pivotMapping2.HasPivotRange(extRng.First) {
				changed = expandToAlign(extRng, pivotMapping2.GetPivotRange(extRng.First))
			}
			if pivotMapping2.HasPivotRange(extRng.Last) {
				lChanged := expandToAlign(extRng, pivotMapping2.GetPivotRange(extRng.Last))
				changed = changed || lChanged
			}

			if changed {
				changed = false
				if pivotMapping1.HasPivotRange(extRng.First) {
					changed = expandToAlign(extRng, pivotMapping1.GetPivotRange(extRng.First))
				}
				if pivotMapping1.HasPivotRange(extRng.Last) {
					lChanged := expandToAlign(extRng, pivotMapping1.GetPivotRange(extRng.Last))
					changed = changed || lChanged
				}
			}
		}

		next = extRng.Last + 1
		l2f, i2f := skipEmpty(extRng.First, extRng.Last+1, pivotMapping1)
		l2l, i2l := skipEmpty(extRng.Last, extRng.First-1, pivotMapping1)
		l3f, i3f := skipEmpty(extRng.First, extRng.Last+1, pivotMapping2)
		l3l, i3l := skipEmpty(extRng.Last, extRng.First-1, pivotMapping2)
		l2 := mapping.PosRange{First: l2f, Last: l2l}
		l3 := mapping.PosRange{First: l3f, Last: l3l}

		if l2.First == -1 && l3.First == -1 { // nothing to export (-1 to -1)
			continue

		} else {

			// empty intersection (expansion got too wide through empty mappings)
			if (i3f > i2l || i3l < i2f) && l3f != -1 && l2f != -1 {
				l3.First = -1
				l3.Last = -1
			}

			if l2.First == -1 && l3.First != -1 {
				mapEmptyL3 = append(mapEmptyL3, mapping.Mapping{From: l2, To: l3})

			} else {
				mapL2L3 = append(mapL2L3, mapping.Mapping{From: l2, To: l3})
			}
		}
	}
	log.Print("INFO: Done")
	log.Print("INFO: Generating output...")

	done := make(chan bool, 2)
	go func() {
		sort.Sort(mapping.SortableMapping(mapL2L3))
		done <- true
	}()
	go func() {
		sort.Sort(mapping.SortableMapping(mapEmptyL3))
		done <- true
	}()
	<-done
	<-done

	mapping.MergeMappings(mapL2L3, mapEmptyL3, func(item mapping.Mapping, pos *mapping.ProcPosition) {
		if pos.Left == -1 && item.From.First > 0 {
			fmt.Println(mapping.Mapping{
				From: mapping.PosRange{
					First: 0,
					Last:  item.From.First - 1,
				},
				To: mapping.NewEmptyPosRange(),
			})
			pos.Left = item.From.First - 1

		} else if pos.Right == -1 && item.To.First > 0 {
			fmt.Println(mapping.Mapping{
				From: mapping.NewEmptyPosRange(),
				To: mapping.PosRange{
					First: 0,
					Last:  item.To.First - 1,
				},
			})
			pos.Right = item.To.First - 1

		}
		if item.From.First > pos.Left+1 {
			fmt.Println(mapping.Mapping{
				From: mapping.PosRange{
					First: pos.Left + 1,
					Last:  item.From.First - 1,
				},
				To: mapping.NewEmptyPosRange(),
			})
		}
		if item.To.First > pos.Right+1 {
			fmt.Println(mapping.Mapping{
				From: mapping.NewEmptyPosRange(),
				To: mapping.PosRange{
					First: pos.Right + 1,
					Last:  item.To.First - 1,
				},
			})
		}
		fmt.Println(item)
	})
}
