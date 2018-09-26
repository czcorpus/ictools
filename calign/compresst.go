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
	"fmt"

	"github.com/czcorpus/ictools/mapping"
)

func compress(data []mapping.Mapping, onItem func(mapping.Mapping)) {
	lastItem := mapping.NewGapMapping(-2, -2, -2, -2) // -2 is an empty value placeholder

	for _, item := range data {
		compressStep(&item, &lastItem, false, onItem)
	}

	if lastItem.From.First != -2 {
		onItem(mkMapping(lastItem.From.First, lastItem.From.Last, true))

	} else if lastItem.To.First != -2 {
		onItem(mkMapping(lastItem.To.First, lastItem.To.Last, false))
	}
}

func TestCompressGaps() {
	data := make([]mapping.Mapping, 0, 5)
	data = append(data, mapping.Mapping{
		From: mapping.PosRange{0, 0},
		To:   mapping.PosRange{-1, -1},
	})
	data = append(data, mapping.Mapping{
		From: mapping.PosRange{1, 2},
		To:   mapping.PosRange{-1, -1},
	})
	data = append(data, mapping.Mapping{
		From: mapping.PosRange{3, 3},
		To:   mapping.PosRange{-1, -1},
	})
	data = append(data, mapping.Mapping{
		From: mapping.PosRange{-1, -1},
		To:   mapping.PosRange{0, 0},
	})
	data = append(data, mapping.Mapping{
		From: mapping.PosRange{-1, -1},
		To:   mapping.PosRange{1, 2},
	})
	data = append(data, mapping.Mapping{
		From: mapping.PosRange{4, 4},
		To:   mapping.PosRange{-1, -1},
	})
	data = append(data, mapping.Mapping{
		From: mapping.PosRange{5, 6},
		To:   mapping.PosRange{3, 3},
	})

	compress(data, func(item mapping.Mapping) {
		fmt.Println("ITEM ", item)
	})

}
