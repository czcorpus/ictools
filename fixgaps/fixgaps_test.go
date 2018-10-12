// Copyright 2018 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2018 Charles University, Faculty of Arts,
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
	"testing"

	"github.com/czcorpus/ictools/mapping"
	"github.com/stretchr/testify/assert"
)

func TestFromChan(t *testing.T) {
	// func FromChan(ch chan []mapping.Mapping, startFromZero bool, struct1Size int, struct2Size int, onItem func(item mapping.Mapp

	ch := make(chan []mapping.Mapping, 1)

	ch <- []mapping.Mapping{
		mapping.NewMapping(1, 1, 0, 2),
		mapping.NewMapping(2, 2, 3, 3),
		mapping.NewMapping(4, 4, 5, 5),
	}
	close(ch)

	ans := make([]mapping.Mapping, 0, 10)
	FromChan(ch, true, 10, 20, func(item mapping.Mapping) {
		ans = append(ans, item)
	})

	assert.Equal(t, mapping.Mapping{From: mapping.PosRange{0, 0}, To: mapping.PosRange{-1, -1}, IsGap: true}, ans[0])
	assert.Equal(t, mapping.Mapping{From: mapping.PosRange{1, 1}, To: mapping.PosRange{0, 2}, IsGap: false}, ans[1])
	assert.Equal(t, mapping.Mapping{From: mapping.PosRange{2, 2}, To: mapping.PosRange{3, 3}, IsGap: false}, ans[2])
	assert.Equal(t, mapping.Mapping{From: mapping.PosRange{3, 3}, To: mapping.PosRange{-1, -1}, IsGap: true}, ans[3])
	assert.Equal(t, mapping.Mapping{From: mapping.PosRange{-1, -1}, To: mapping.PosRange{4, 4}, IsGap: true}, ans[4])
	assert.Equal(t, mapping.Mapping{From: mapping.PosRange{4, 4}, To: mapping.PosRange{5, 5}, IsGap: false}, ans[5])
	assert.Equal(t, mapping.Mapping{From: mapping.PosRange{5, 9}, To: mapping.PosRange{-1, -1}, IsGap: true}, ans[6])
	assert.Equal(t, mapping.Mapping{From: mapping.PosRange{-1, -1}, To: mapping.PosRange{6, 19}, IsGap: true}, ans[7])
	assert.Equal(t, 8, len(ans))
}
