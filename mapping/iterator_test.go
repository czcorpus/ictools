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

// Package mapping provides data types and functions used
// to manipulate numeric mapping between two aligned structures

package mapping

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIteratorFactory(t *testing.T) {
	mList := make([]Mapping, 4)
	mList[0] = Mapping{PosRange{1, 2}, PosRange{-1, -1}, false}
	mList[1] = Mapping{PosRange{3, 4}, PosRange{-1, -1}, false}
	mList[2] = Mapping{PosRange{5, 5}, PosRange{-1, -1}, false}
	mList[3] = Mapping{PosRange{6, 6}, PosRange{0, 0}, false}

	itr := NewIterator(mList)

	assert.True(t, itr.Unfinished())
}

func TestIteratorFactoryEmptyData(t *testing.T) {
	mList := make([]Mapping, 0)
	itr := NewIterator(mList)
	assert.True(t, !itr.Unfinished())
}

func TestIteratorApplyFinishes(t *testing.T) {
	mList := make([]Mapping, 1)
	mList[0] = Mapping{PosRange{1, 2}, PosRange{-1, -1}, false}
	itr := NewIterator(mList)

	itr.Next()
	assert.True(t, itr.Unfinished())

	itr.Apply(func(v Mapping) {
	})

	assert.False(t, itr.Unfinished())
}

func TestIteratorCanApplyWithoutNext(t *testing.T) {
	mList := make([]Mapping, 1)
	mList[0] = Mapping{PosRange{1, 2}, PosRange{-1, -1}, false}
	itr := NewIterator(mList)
	itr.Apply(func(v Mapping) {
		assert.Equal(t, mList[0], v)
	})
}

func TestIteratorHasPriorityOver(t *testing.T) {
	mList := make([]Mapping, 4)
	mList[0] = Mapping{PosRange{1, 2}, PosRange{-1, -1}, false}
	mList[1] = Mapping{PosRange{3, 4}, PosRange{-1, -1}, false}
	mList[2] = Mapping{PosRange{5, 5}, PosRange{-1, -1}, false}
	mList[3] = Mapping{PosRange{6, 6}, PosRange{0, 0}, false}
	itr1 := NewIterator(mList)
	itr2 := NewIterator(mList)

	itr1.Next()
	itr2.Next()
	itr2.Next()
	assert.False(t, itr1.HasPriorityOver(&itr2)) // "to" values are in play

	itr2.Next()
	itr2.Next()
	assert.True(t, itr1.HasPriorityOver(&itr2)) // finally we're at (0, 0) against (-1, -1)
}

func TestMergeMappings(t *testing.T) {
	mList1 := make([]Mapping, 4)
	mList1[0] = Mapping{PosRange{1, 2}, PosRange{-1, -1}, false}
	mList1[1] = Mapping{PosRange{3, 3}, PosRange{1, 1}, false}
	mList1[2] = Mapping{PosRange{4, 5}, PosRange{4, 4}, false}
	mList1[3] = Mapping{PosRange{6, 6}, PosRange{6, 7}, false}

	mList2 := make([]Mapping, 2)
	mList2[0] = Mapping{PosRange{-1, -1}, PosRange{2, 3}, false}
	mList2[1] = Mapping{PosRange{-1, -1}, PosRange{5, 5}, false}

	i := 0
	ans := make([]Mapping, 6)
	MergeMappings(mList1, mList2, func(item Mapping) {
		ans[i] = item
		i++
	})

	assert.Equal(t, mList1[0], ans[0])
	assert.Equal(t, mList1[1], ans[1])
	assert.Equal(t, mList2[0], ans[2])
	assert.Equal(t, mList1[2], ans[3])
	assert.Equal(t, mList2[1], ans[4])
	assert.Equal(t, mList1[3], ans[5])
}

func TestMergeMappingsAlternatingItems(t *testing.T) {
	mList1 := make([]Mapping, 2)
	mList1[0] = Mapping{PosRange{1, 1}, PosRange{1, 1}, false}
	mList1[1] = Mapping{PosRange{2, 2}, PosRange{4, 4}, false}

	mList2 := make([]Mapping, 2)
	mList2[0] = Mapping{PosRange{-1, -1}, PosRange{2, 3}, false}
	mList2[1] = Mapping{PosRange{-1, -1}, PosRange{5, 5}, false}

	i := 0
	ans := make([]Mapping, 4)
	MergeMappings(mList1, mList2, func(item Mapping) {
		ans[i] = item
		i++
	})

	assert.Equal(t, mList1[0], ans[0])
	assert.Equal(t, mList2[0], ans[1])
	assert.Equal(t, mList1[1], ans[2])
	assert.Equal(t, mList2[1], ans[3])
}

func TestMergeMappingsWaitingColumn(t *testing.T) {
	mList1 := make([]Mapping, 2)
	mList1[0] = Mapping{PosRange{1, 1}, PosRange{3, 3}, false}
	mList1[1] = Mapping{PosRange{2, 2}, PosRange{4, 4}, false}

	mList2 := make([]Mapping, 2)
	mList2[0] = Mapping{PosRange{-1, -1}, PosRange{1, 1}, false}
	mList2[1] = Mapping{PosRange{-1, -1}, PosRange{2, 2}, false}

	i := 0
	ans := make([]Mapping, 4)
	MergeMappings(mList1, mList2, func(item Mapping) {
		ans[i] = item
		fmt.Println(item)
		i++
	})

	assert.Equal(t, mList2[0], ans[0])
	assert.Equal(t, mList2[1], ans[1])
	assert.Equal(t, mList1[0], ans[2])
	assert.Equal(t, mList1[1], ans[3])
}

func TestMergeMappingsEmptySources(t *testing.T) {
	mList1 := make([]Mapping, 0)
	mList2 := make([]Mapping, 0)
	i := 0
	MergeMappings(mList1, mList2, func(item Mapping) {
		i++
	})
	assert.Equal(t, 0, i)
}

func TestStringMethod(t *testing.T) {
	m := Mapping{From: PosRange{1, 2}, To: PosRange{3, 4}, IsGap: true}
	assert.Equal(t, "1,2\t3,4\tg", m.String())

	m2 := Mapping{From: PosRange{1, 2}, To: PosRange{3, 4}, IsGap: false}
	assert.Equal(t, "1,2\t3,4", m2.String())

	m3 := Mapping{From: PosRange{1, 1}, To: PosRange{3, 3}, IsGap: false}
	assert.Equal(t, "1\t3", m3.String())
}

func TestNewGapMapping(t *testing.T) {
	m := NewGapMapping(1, 2, 3, 4)
	assert.Equal(t, 1, m.From.First)
	assert.Equal(t, 2, m.From.Last)
	assert.Equal(t, 3, m.To.First)
	assert.Equal(t, 4, m.To.Last)
}

func TestIsEmpty(t *testing.T) {
	m := Mapping{From: PosRange{-1, -1}, To: PosRange{-1, -1}, IsGap: false}
	assert.True(t, m.IsEmpty())
	m2 := Mapping{From: PosRange{-1, -1}, To: PosRange{-1, -1}, IsGap: false}
	assert.True(t, m2.IsEmpty())

	m3 := Mapping{From: PosRange{2, -1}, To: PosRange{-1, -1}, IsGap: false}
	assert.False(t, m3.IsEmpty())
}
