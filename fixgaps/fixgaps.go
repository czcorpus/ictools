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

// Package fixgaps provides functions for filling missing mapping
// lines to an extracted alignment file.
package fixgaps

import (
	"bufio"
	"os"
	"strings"

	"github.com/czcorpus/ictools/mapping"
)

// FromFile inserts [-1, a] or [a, -1] between identifiers
// A1 and A2 where A2 > A1+1 (but also with respect to two possible
// positions in a column).
// Data are read from file 'file'. If startFromZero is true then
// the list starts from zero else from the first found item.
// The function does not print anything to stdout.
func FromFile(file *os.File, startFromZero bool, onItem func(item mapping.Mapping)) {
	fr := bufio.NewScanner(file)
	lastL1 := -1
	lastL2 := -1
	for fr.Scan() {
		line := fr.Text()
		items := strings.Split(line, "\t")
		l1t := strings.Split(items[0], ",")
		l2t := strings.Split(items[1], ",")
		r1 := mapping.NewPosRange(l1t)
		r2 := mapping.NewPosRange(l2t)
		if !startFromZero && lastL1 == -1 && lastL2 == -1 {
			lastL1 = r1.First
			lastL2 = r2.First
		}
		for r1.First > lastL1+1 {
			lastL1++
			onItem(mapping.NewMapping(lastL1, lastL1, -1, -1))
		}
		for r2.First > lastL2+1 {
			lastL2++
			onItem(mapping.NewMapping(-1, -1, lastL2, lastL2))
		}
		if r1.Last != -1 {
			lastL1 = r1.Last
		}
		if r2.Last != -1 {
			lastL2 = r2.Last
		}
		onItem(mapping.Mapping{r1, r2})
	}
}

// FromChan is the same as FromFile except from the source
// of data. In this case, a channel is used.
func FromChan(ch chan []mapping.Mapping, startFromZero bool, onItem func(item mapping.Mapping)) {
	lastL1 := -1
	lastL2 := -1
	for buff := range ch {
		for _, item := range buff {
			if !startFromZero && lastL1 == -1 && lastL2 == -1 {
				lastL1 = item.From.First
				lastL2 = item.To.First
			}
			for item.From.First > lastL1+1 {
				lastL1++
				onItem(mapping.NewMapping(lastL1, lastL1, -1, -1))
			}
			for item.To.First > lastL2+1 {
				lastL2++
				onItem(mapping.NewMapping(-1, -1, lastL2, lastL2))
			}
			if item.From.Last != -1 {
				lastL1 = item.From.Last
			}
			if item.To.Last != -1 {
				lastL2 = item.To.Last
			}
			onItem(item)
		}
	}
}
