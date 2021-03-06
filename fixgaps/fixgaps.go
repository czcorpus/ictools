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
	"fmt"
	"log"
	"os"

	"github.com/czcorpus/ictools/mapping"
)

type FixGapsError struct {
	Item  mapping.Mapping
	Left  int
	Pivot int
}

func (f *FixGapsError) Error() string {
	return fmt.Sprintf("alignment [%s] overlaps an already covered range (LEFT position: %d, PIVOT position: %d)", f.Item, f.Left, f.Pivot)
}

func NewFixGapsError(item mapping.Mapping, left int, pivot int) *FixGapsError {
	return &FixGapsError{Item: item, Left: left, Pivot: pivot}
}

// FromFile inserts [-1, a] or [a, -1] between identifiers
// A1 and A2 where A2 > A1+1 (but also with respect to two possible
// positions in a column).
// Data are read from 'file'. If startFromZero is true then
// the list is always build so it starts from position 0.
// Otherwise, the list starts from the first found item.
// The function does not print anything to stdout.
func FromFile(file *os.File, startFromZero bool, struct1Size int, struct2Size int, onItem func(item mapping.Mapping)) {
	fr := bufio.NewScanner(file)
	lastL1 := -1
	lastL2 := -1
	for i := 0; fr.Scan(); i++ {
		item, err := mapping.NewMappingFromString(fr.Text())
		if err != nil {
			fmt.Printf("[WARNING] Failed to process line %d: %s", i, err)
			continue
		}
		if !startFromZero && lastL1 == -1 && lastL2 == -1 {
			lastL1 = item.From.First
			lastL2 = item.From.First
		}
		for item.From.First > lastL1+1 {
			lastL1++
			onItem(mapping.NewGapMapping(lastL1, lastL1, -1, -1))
		}
		for item.To.First > lastL2+1 {
			lastL2++
			onItem(mapping.NewGapMapping(-1, -1, lastL2, lastL2))
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

// FromChan is the same as FromFile except from the source
// of data. In this case, a channel is used.
func FromChan(ch chan []mapping.Mapping, startFromZero bool, struct1Size int, struct2Size int,
	onItem func(item mapping.Mapping, err *FixGapsError)) {
	lastL1 := -1
	lastL2 := -1
	for buff := range ch {
		for _, item := range buff {
			var err *FixGapsError
			if item.From.First != -1 && item.From.First <= lastL1 ||
				item.To.First != -1 && item.To.First <= lastL2 {
				err = NewFixGapsError(item, lastL1, lastL2)
				onItem(mapping.Mapping{}, err)
			}
			if !startFromZero && lastL1 == -1 && lastL2 == -1 {
				lastL1 = item.From.First
				lastL2 = item.To.First
			}
			for item.From.First > lastL1+1 {
				lastL1++
				onItem(mapping.NewGapMapping(lastL1, lastL1, -1, -1), nil)
			}
			for item.To.First > lastL2+1 {
				lastL2++
				onItem(mapping.NewGapMapping(-1, -1, lastL2, lastL2), nil)
			}
			if item.From.Last != -1 {
				lastL1 = item.From.Last
			}
			if item.To.Last != -1 {
				lastL2 = item.To.Last
			}
			onItem(item, nil)
		}
	}

	if lastL1 < struct1Size-1 {
		log.Printf("WARNING: Filled in missing end %d,%d in the LEFT language. Please make sure this is correct.", lastL1+1, struct1Size-1)
		onItem(mapping.Mapping{
			From: mapping.PosRange{
				First: lastL1 + 1,
				Last:  struct1Size - 1,
			},
			To:    mapping.NewEmptyPosRange(),
			IsGap: true,
		}, nil)
	}

	if lastL2 < struct2Size-1 {
		log.Printf("WARNING: Filled in missing end %d,%d in the PIVOT language. Please make sure this is correct.", lastL2+1, struct2Size-1)
		onItem(mapping.Mapping{
			From: mapping.NewEmptyPosRange(),
			To: mapping.PosRange{
				First: lastL2 + 1,
				Last:  struct2Size - 1,
			},
			IsGap: true,
		}, nil)
	}
}
