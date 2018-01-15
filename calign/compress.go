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

// CompressFromChan reduces subsequent lines with -1 in one of the columns
// to a single line with proper range (e.g. "-1   am,an" where 'am' is the
// beginning of the first line in the series and 'am' is the end of the last
// line in the series.
func CompressFromChan(ch chan []mapping.Mapping, onItem func(item mapping.Mapping)) {
	fromFirst := -2 // -2 is an empty value placeholder
	fromLast := -2
	toFirst := -2
	toLast := -2

	for buff := range ch {
		for _, item := range buff {

			if item.To.First == -1 {
				if fromFirst == -2 {
					fromFirst = item.From.First
					fromLast = item.From.Last

				} else {
					fromLast = item.From.Last
				}
				continue

			} else if fromFirst != -2 {
				onItem(mkMapping(fromFirst, fromLast, true))
				fromFirst = -2
			}

			if item.From.First == -1 {
				if toFirst == -2 {
					toFirst = item.To.First
					toLast = item.To.Last

				} else {
					toLast = item.To.Last
				}
				continue

			} else if toFirst != -2 {
				onItem(mkMapping(toFirst, toLast, false))
				toFirst = -2
			}
			onItem(item)
		}
	}

	if fromFirst != -2 {
		onItem(mkMapping(fromFirst, fromLast, true))

	} else if toFirst != -2 {
		onItem(mkMapping(toFirst, toLast, false))
	}
}
