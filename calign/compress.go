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

package calign

import (
	"bufio"
	"log"
	"os"

	"github.com/czcorpus/ictools/mapping"
)

func mkLeftToEmpty(beg int, end int, isGap bool) mapping.Mapping {
	if isGap {
		return mapping.NewGapMapping(beg, end, -1, -1)
	}
	return mapping.NewMapping(beg, end, -1, -1)
}

func mkEmptyToRight(beg int, end int, isGap bool) mapping.Mapping {
	if isGap {
		return mapping.NewGapMapping(-1, -1, beg, end)
	}
	return mapping.NewMapping(-1, -1, beg, end)
}

// compressStep decides whether 'item' should be either added to
// one of currently expanded ranges (currRanges) or directly printed.
// Please note that 'currRanges' is of a little misused type here as it
// stores no concrete mapping line we want eventually store but rather currently
// reached non-empty ranges for left and right sizes.
func compressStep(item *mapping.Mapping, currRanges *mapping.Mapping, gapsOnly bool, onItem func(item mapping.Mapping)) {
	if item.To.First == -1 && (gapsOnly && item.IsGap || !gapsOnly) {
		if currRanges.From.First == -2 {
			currRanges.From.First = item.From.First
			currRanges.From.Last = item.From.Last
			currRanges.IsGap = item.IsGap

		} else {
			currRanges.From.Last = item.From.Last
		}
		return

	} else if currRanges.From.First != -2 {
		onItem(mkLeftToEmpty(currRanges.From.First, currRanges.From.Last, currRanges.IsGap))
		currRanges.From.First = -2
		currRanges.From.Last = -2
		currRanges.IsGap = false
		// TODO also reset From.Last
	}
	if item.From.First == -1 && (gapsOnly && item.IsGap || !gapsOnly) {
		if currRanges.To.First == -2 {
			currRanges.To.First = item.To.First
			currRanges.To.Last = item.To.Last
			currRanges.IsGap = item.IsGap

		} else {
			currRanges.To.Last = item.To.Last
		}
		return

	} else if currRanges.To.First != -2 {
		onItem(mkEmptyToRight(currRanges.To.First, currRanges.To.Last, currRanges.IsGap))
		currRanges.To.First = -2
		currRanges.To.Last = -2
		currRanges.IsGap = false
	}
	onItem(*item)
}

// CompressFromChan reduces subsequent lines with -1 in one of the columns
// to a single line with proper range (e.g. "-1   am,an" where 'am' is the
// beginning of the first line in the series and 'an' is the end of the last
// line in the series.
func CompressFromChan(ch chan []mapping.Mapping, gapsOnly bool, onItem func(mapping.Mapping)) {
	currRanges := mapping.NewMapping(-2, -2, -2, -2) // -2 is an empty value placeholder

	for buff := range ch {
		for _, item := range buff {
			compressStep(&item, &currRanges, gapsOnly, onItem)
		}
	}

	if currRanges.From.First != -2 {
		onItem(mkLeftToEmpty(currRanges.From.First, currRanges.From.Last, currRanges.IsGap))
	}
	if currRanges.To.First != -2 {
		onItem(mkEmptyToRight(currRanges.To.First, currRanges.To.Last, currRanges.IsGap))
	}
}

// CompressFromFile runs in the same way as CompressFromChan except that
// the data source is a file in this case.
func CompressFromFile(file *os.File, gapsOnly bool, onItem func(item mapping.Mapping)) {
	fr := bufio.NewScanner(file)
	currRanges := mapping.NewMapping(-2, -2, -2, -2) // -2 is an empty value placeholder

	for i := 0; fr.Scan(); i++ {
		item, err := mapping.NewMappingFromString(fr.Text())
		if err == nil {
			compressStep(&item, &currRanges, gapsOnly, onItem)

		} else {
			log.Printf("ERROR: Failed to process line %d: %s", i, err)
		}
	}

	if currRanges.From.First != -2 {
		onItem(mkLeftToEmpty(currRanges.From.First, currRanges.From.Last, currRanges.IsGap))
	}
	if currRanges.To.First != -2 {
		onItem(mkEmptyToRight(currRanges.To.First, currRanges.To.Last, currRanges.IsGap))
	}
}
