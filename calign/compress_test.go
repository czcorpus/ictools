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
	"log"
	"testing"

	"github.com/czcorpus/ictools/mapping"
	"github.com/stretchr/testify/assert"
)

func compress(data []mapping.Mapping, onItem func(mapping.Mapping)) {
	lastItem := mapping.NewMapping(-2, -2, -2, -2) // -2 is an empty value placeholder

	for _, item := range data {
		compressStep(&item, &lastItem, false, onItem)
	}

	if lastItem.From.First != -2 {
		onItem(mkMapping(lastItem.From.First, lastItem.From.Last, true))

	} else if lastItem.To.First != -2 {
		onItem(mkMapping(lastItem.To.First, lastItem.To.Last, false))
	}
}

func TestCompress(t *testing.T) {
	data := make([]mapping.Mapping, 0, 5)
	data = append(data, mapping.NewMapping(0, 0, -1, -1))
	data = append(data, mapping.NewMapping(1, 2, -1, -1))
	data = append(data, mapping.NewMapping(3, 3, -1, -1))
	data = append(data, mapping.NewMapping(-1, -1, 0, 0))
	data = append(data, mapping.NewMapping(-1, -1, 1, 2))
	data = append(data, mapping.NewMapping(4, 4, -1, -1))
	data = append(data, mapping.NewMapping(5, 6, 3, 3))

	valData := []mapping.Mapping{
		mapping.NewGapMapping(0, 3, -1, -1),
		mapping.NewGapMapping(4, 4, -1, -1),
		mapping.NewGapMapping(-1, -1, 0, 2),
		mapping.NewMapping(5, 6, 3, 3),
	}

	i := 0
	compress(data, func(item mapping.Mapping) {
		log.Print(item)
		assert.Equal(t, valData[i], item)
		i++
	})

}
