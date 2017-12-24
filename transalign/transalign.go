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
	"fmt"
	"log"
	"sort"
	"time"
)

type mapping struct {
	from PosRange
	to   PosRange
}

func (m mapping) String() string {
	return fmt.Sprintf("%s\t%s", m.from, m.to)
}

// ---------------------------

type sortableMapping []mapping

func (sm sortableMapping) Len() int {
	return len(sm)
}

func (sm sortableMapping) Swap(i, j int) {
	sm[i], sm[j] = sm[j], sm[i]
}

func (sm sortableMapping) Less(i, j int) bool {
	if sm[i].from.first > -1 && sm[j].from.first > -1 {
		return sm[i].from.first > sm[j].from.first

	} else if sm[i].from.first == -1 || sm[j].from.first == -1 {
		return sm[i].to.first > sm[j].to.first

	} else if sm[i].to.first == -1 && sm[j].to.first == -1 {
		return sm[i].from.first > sm[j].from.first
	}
	panic("unknow type combination")
}

// ---------------------------

func max(v1 int, v2 int) int {
	if v1 > v2 {
		return v1
	}
	return v2
}

func calcFinalItem(mapL2L3 []mapping, mapEmptyL3 []mapping) mapping {
	ifinal := -1
	if len(mapL2L3) > 0 {
		ifinal = max(mapL2L3[0].from.last, mapL2L3[0].to.last)
	}
	if len(mapEmptyL3) > 0 {
		ifinal = max(ifinal, mapEmptyL3[0].to.last) // TODO !!! overit
	}

	return mapping{
		PosRange{ifinal, ifinal},
		PosRange{ifinal, ifinal},
	}
}

func Run(halfMapping1 *HalfMapping, halfMapping2 *HalfMapping) {
	next := 0
	mapL2L3 := make([]mapping, 0, halfMapping1.RangesSize()) // TODO size estimation
	// we have to keep one of [-1, x], [x, -1] mapping separate
	// because these two cannot be sorted together in a traditional way
	mapEmptyL3 := make([]mapping, 0, halfMapping1.RangesSize()/10) // 10 is just an estimate
	log.Print("Computing new alignment:")
	t1 := time.Now().UnixNano()
	for i, rng := range halfMapping1.ranges {
		if i == 1000 {
			t1 = time.Now().UnixNano() - t1
			log.Printf("estimated proc. time: %01.2f seconds.", float64(t1)*1e-9*1e-3*float64(halfMapping1.RangesSize()))
		}
		if i < next {
			continue
		}
		changed := true
		for changed {
			changed = false
			if halfMapping2.HasRange(rng.first) {
				rng, changed = ensureContains(rng, halfMapping2.GetRange(rng.first))
			}
			if halfMapping2.HasRange(rng.last) {
				var lChanged bool
				rng, lChanged = ensureContains(rng, halfMapping2.GetRange(rng.last))
				changed = changed || lChanged
			}
			if changed {
				halfMapping1.SetRange(i, rng)
				changed = false
				if halfMapping1.HasRange(rng.first) {
					rng, changed = ensureContains(rng, halfMapping1.GetRange(rng.first))
				}
				if halfMapping1.HasRange(rng.last) {
					var lChanged bool
					rng, lChanged = ensureContains(rng, halfMapping1.GetRange(rng.last))
					changed = changed || lChanged
				}
			}
		}
		next = rng.last + 1
		l2 := PosRange{
			skipEmpty(rng.first, rng.last+1, halfMapping1.PivotMap()),
			skipEmpty(rng.last, rng.first-1, halfMapping1.PivotMap()),
		}
		l3 := PosRange{
			skipEmpty(rng.first, rng.last+1, halfMapping2.PivotMap()),
			skipEmpty(rng.last, rng.first-1, halfMapping2.PivotMap()),
		}
		if l2.first == -1 && l3.first == -1 { // nothing to export (-1 to -1)
			continue

		} else if l2.first != -1 && l3.first == -1 {
			mapEmptyL3 = append(mapEmptyL3, mapping{l2, l3})

		} else {
			mapL2L3 = append(mapL2L3, mapping{l2, l3})
		}
	}

	log.Print("Sorting intermediate data...")
	sort.Sort(sortableMapping(mapL2L3))
	sort.Sort(sortableMapping(mapEmptyL3))
	log.Print("...done.")
	log.Print("Generating output...")

	finalMapping := calcFinalItem(mapL2L3, mapEmptyL3)
	iterL2L3 := newMapIterator(mapL2L3, finalMapping)
	iterL3 := newMapIterator(mapEmptyL3, finalMapping)

	vL2L3 := iterL2L3.Next()
	vL3 := iterL3.Next()
	var curr mapping
	for vL2L3 != finalMapping || vL3 != finalMapping {
		curr = vL2L3
		if vL3.from.LessThan(vL2L3.from) {
			curr = vL3
			vL3 = iterL3.Next()

		} else {
			vL2L3 = iterL2L3.Next()
		}
		fmt.Println(curr)
	}
	log.Print("...done.")
}
