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
	"os"
	"path/filepath"
	"testing"

	"github.com/czcorpus/ictools/mapping"
	"github.com/stretchr/testify/assert"
)

func TestCompressStepLeftEmptyMerge(t *testing.T) {
	currRanges := mapping.NewMapping(-2, -2, 2, 3)
	newItem := mapping.NewMapping(-1, -1, 4, 5)
	i := 0
	compressStep(&newItem, &currRanges, false, func(item mapping.Mapping) {
		i++
	})
	assert.Equal(t, 0, i)
	assert.Equal(t, 2, currRanges.To.First)
	assert.Equal(t, 5, currRanges.To.Last)
}

// TestCompressStepLeftEmptyNoMerge adds an item we should ignore
// because it does not extend the current left range 1..4. On the other
// hand we create a new range 13..17 as the current right range has not been
// initialized yet (-2...-2).
func TestCompressStepLeftEmptyNoMerge(t *testing.T) {
	currRanges := mapping.NewMapping(1, 4, -2, -2)
	newItem := mapping.NewMapping(-1, -1, 13, 17)
	var closedItem *mapping.Mapping
	i := 0
	compressStep(&newItem, &currRanges, false, func(item mapping.Mapping) {
		closedItem = &item
		i++
	})
	assert.Equal(t, 1, i)
	assert.Equal(t, 1, closedItem.From.First)
	assert.Equal(t, 4, closedItem.From.Last)
	assert.Equal(t, -2, currRanges.From.First)
	assert.Equal(t, -2, currRanges.From.Last)
	assert.Equal(t, 13, currRanges.To.First)
	assert.Equal(t, 17, currRanges.To.Last)
}

func TestCompressStepRightEmptyMerge(t *testing.T) {
	currRanges := mapping.NewMapping(1, 4, -2, -2)
	newItem := mapping.NewMapping(4, 7, -1, -1)
	i := 0
	compressStep(&newItem, &currRanges, false, func(item mapping.Mapping) {
		i++
	})
	assert.Equal(t, 0, i)
	assert.Equal(t, 1, currRanges.From.First)
	assert.Equal(t, 7, currRanges.From.Last)
}

// TestCompressStepRightEmptyNoMerge adds an item we should ignore
// because it does not extend the current right range 5..9. On the other
// hand we create a new range 21..23 as the current left range has not been
// initialized yet (-2...-2).
func TestCompressStepRightEmptyNoMerge(t *testing.T) {
	currRanges := mapping.NewMapping(-2, -2, 5, 9)
	newItem := mapping.NewMapping(21, 23, -1, -1)
	i := 0
	compressStep(&newItem, &currRanges, false, func(item mapping.Mapping) {
		i++
	})
	assert.Equal(t, 0, i)
	assert.Equal(t, 5, currRanges.To.First)
	assert.Equal(t, 9, currRanges.To.Last) // not reset (but it should se respective todo)
	assert.Equal(t, 21, currRanges.From.First)
	assert.Equal(t, 23, currRanges.From.Last)
}

func TestCompressFromFile(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f, err := os.Open(filepath.Join(cwd, "..", "testdata", "to-compress.txt"))
	if err != nil {
		panic(err)
	}

	ans := make([]*mapping.Mapping, 0, 10)
	CompressFromFile(f, true, func(item mapping.Mapping) {
		ans = append(ans, &item)
	})
	validate := []mapping.Mapping{
		mapping.NewMapping(0, 0, 0, 0),
		mapping.NewMapping(-1, -1, 1, 1),
		mapping.NewMapping(1, 1, -1, -1),
		mapping.NewMapping(2, 2, 2, 2),
		mapping.NewMapping(3, 4, 3, 4),
		mapping.NewGapMapping(4, 6, -1, -1),
		mapping.NewGapMapping(-1, -1, 5, 7),
		mapping.NewMapping(7, 7, -1, -1),
	}
	for i, v := range ans {
		log.Print(v)
		assert.Equal(t, &validate[i], v)
	}
	assert.Equal(t, 8, len(ans))

}
