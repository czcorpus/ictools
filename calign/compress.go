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

func mkMapping(beg int, end int, rightEmpty bool) mapping.Mapping {
	if beg == end {
		if rightEmpty {
			return mapping.NewMapping(beg, beg, -1, -1)
		}
		return mapping.NewMapping(-1, -1, beg, beg)
	}
	if rightEmpty {
		return mapping.NewMapping(beg, end, -1, -1)
	}
	return mapping.NewMapping(-1, -1, beg, end)
}

func compressStep(item *mapping.Mapping, lastItem *mapping.Mapping, onItem func(item mapping.Mapping)) {
	if item.To.First == -1 {
		if lastItem.From.First == -2 {
			lastItem.From.First = item.From.First
			lastItem.From.Last = item.From.Last

		} else {
			lastItem.From.Last = item.From.Last
		}
		return

	} else if lastItem.From.First != -2 {
		onItem(mkMapping(lastItem.From.First, lastItem.From.Last, true))
		lastItem.From.First = -2
	}

	if item.From.First == -1 {
		if lastItem.To.First == -2 {
			lastItem.To.First = item.To.First
			lastItem.To.Last = item.To.Last

		} else {
			lastItem.To.Last = item.To.Last
		}
		return

	} else if lastItem.To.First != -2 {
		onItem(mkMapping(lastItem.To.First, lastItem.To.Last, false))
		lastItem.To.First = -2
	}
	onItem(*item)
}

// CompressFromChan reduces subsequent lines with -1 in one of the columns
// to a single line with proper range (e.g. "-1   am,an" where 'am' is the
// beginning of the first line in the series and 'an' is the end of the last
// line in the series.
func CompressFromChan(ch chan []mapping.Mapping, onItem func(item mapping.Mapping)) {
	lastItem := mapping.NewMapping(-2, -2, -2, -2) // -2 is an empty value placeholder

	for buff := range ch {
		for _, item := range buff {
			compressStep(&item, &lastItem, onItem)
		}
	}

	if lastItem.From.First != -2 {
		onItem(mkMapping(lastItem.From.First, lastItem.From.Last, true))

	} else if lastItem.To.First != -2 {
		onItem(mkMapping(lastItem.To.First, lastItem.To.Last, false))
	}
}

// CompressFromFile runs in the same way as CompressFromChan except that
// the data source is a file in this case.
func CompressFromFile(file *os.File, onItem func(item mapping.Mapping)) {
	fr := bufio.NewScanner(file)
	lastItem := mapping.NewMapping(-2, -2, -2, -2) // -2 is an empty value placeholder

	for i := 0; fr.Scan(); i++ {
		item, err := mapping.NewMappingFromString(fr.Text())
		if err == nil {
			compressStep(&item, &lastItem, onItem)

		} else {
			log.Printf("[WARNING] Failed to process line %d: %s", i, err)
		}
	}

	if lastItem.From.First != -2 {
		onItem(mkMapping(lastItem.From.First, lastItem.From.Last, true))

	} else if lastItem.To.First != -2 {
		onItem(mkMapping(lastItem.To.First, lastItem.To.Last, false))
	}
}
