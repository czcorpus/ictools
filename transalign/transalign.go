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

// Run implements an algorith for finding a mapping
// between L1 and L1 based on two "half mappings"
// L1 -> LP and L2 -> LP.
func Run(pivotMapping1 *PivotMapping, pivotMapping2 *PivotMapping) {

	mapL1L2 := make([]mapping.Mapping, 0, pivotMapping1.Size()) // TODO size estimation
	// we have to keep one of [-1, x], [x, -1] mapping separate
	// because these two cannot be sorted together in a traditional way
	mapEmptyL2 := make([]mapping.Mapping, 0, pivotMapping1.Size()/10) // 10 is just an estimate
	log.Print("INFO: Computing new alignment...")

	log.Print("INFO: Done")
	log.Print("INFO: Generating output...")

	//l2Start := 0
	//l2End := 0
	//pivot2Pos := 0
	//var pivot2Row mapping.PosRange

	for _, l1Row := range pivotMapping1.ranges {
		fmt.Println(l1Row)
	}

	done := make(chan bool, 2)
	go func() {
		sort.Sort(mapping.SortableMapping(mapL1L2))
		done <- true
	}()
	go func() {
		sort.Sort(mapping.SortableMapping(mapEmptyL2))
		done <- true
	}()
	<-done
	<-done

	mapping.MergeMappings(mapL1L2, mapEmptyL2, func(item mapping.Mapping) {
		fmt.Println(item)
	})
}
