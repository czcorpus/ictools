// Copyright 2020 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2020 Charles University, Faculty of Arts,
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

package queue

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestDeque(t *testing.T) {
	q := &Deque{}
	q.PushBack("one", nil)
	q.PushBack("two", nil)
	q.PushBack("three", nil)
	assert.Equal(t, 3, q.Size())
	ans, err := q.PopFront()
	assert.Equal(t, "one", ans.GroupID)
	assert.Nil(t, err)
	assert.Equal(t, 2, q.Size())
	ans, err = q.PopFront()
	assert.Nil(t, err)
	assert.Equal(t, "two", ans.GroupID)
	assert.Equal(t, 1, q.Size())
	ans, err = q.PopFront()
	assert.Nil(t, err)
	assert.Equal(t, "three", ans.GroupID)
	assert.Equal(t, 0, q.Size())
}

func TestDequeRemoveFromEmpty(t *testing.T) {
	q := &Deque{}
	_, err := q.PopFront()
	assert.Error(t, err)
	assert.Equal(t, 0, q.Size())
}

func TestDequeGetFirstGroup(t *testing.T) {
	q := &Deque{}
	q.PushBack("one", nil)
	q.PushBack("two", nil)
	assert.Equal(t, "one", q.FrontGroup())
}

func TestDequeFirstGroupEmpty(t *testing.T) {
	q := &Deque{}
	assert.Equal(t, "", q.FrontGroup())
	assert.Equal(t, 0, q.Size())
}

func TestDequeGetLastGroup(t *testing.T) {
	q := &Deque{}
	q.PushBack("one", nil)
	q.PushBack("two", nil)
	assert.Equal(t, "two", q.BackGroup())
}

func TestDequeLastGroupEmpty(t *testing.T) {
	q := &Deque{}
	assert.Equal(t, "", q.BackGroup())
}

func TestDequeRemoveLast(t *testing.T) {
	q := &Deque{}
	q.PushBack("one", nil)
	q.PushBack("two", nil)
	q.PushBack("three", nil)
	v, err := q.PopBack()
	assert.Equal(t, "three", v.GroupID)
	assert.Equal(t, "two", q.last.GroupID)
	assert.Equal(t, 2, q.Size())
	assert.Nil(t, err)
}

func TestDequeRemoveLastOnEmpty(t *testing.T) {
	q := &Deque{}
	v, err := q.PopBack()
	assert.Nil(t, v)
	assert.Error(t, err)
}

func TestDequeRemoveLastOnSizeOne(t *testing.T) {
	q := &Deque{}
	q.PushBack("one", nil)
	v, err := q.PopBack()
	assert.Equal(t, "one", v.GroupID)
	assert.Nil(t, q.last)
	assert.Nil(t, q.first)
	assert.Nil(t, err)
	assert.Equal(t, 0, q.Size())
}

func TestDequeAddFirst(t *testing.T) {
	q := &Deque{}
	q.PushBack("one", nil)
	q.PushBack("two", nil)
	q.PushBack("three", nil)
	q.PushFront("zero", nil)
	assert.Equal(t, q.first.GroupID, "zero")
	assert.Equal(t, q.last.GroupID, "three")
	assert.Equal(t, 4, q.Size())
}

func TestDequeAddFirstOnSizeOne(t *testing.T) {
	q := &Deque{}
	q.PushBack("one", nil)
	q.PushFront("zero", nil)
	assert.Equal(t, q.first.GroupID, "zero")
	assert.Equal(t, q.last.GroupID, "one")
	assert.Equal(t, 2, q.Size())
}

func TestDequeAddFirstOnEmpty(t *testing.T) {
	q := &Deque{}
	q.PushFront("zero", nil)
	assert.Equal(t, q.first.GroupID, "zero")
	assert.Equal(t, q.last.GroupID, "zero")
	assert.Equal(t, 1, q.Size())
}
