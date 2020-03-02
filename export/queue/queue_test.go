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

func TestQueue(t *testing.T) {
	q := &Queue{}
	q.AddLast("one", nil)
	q.AddLast("two", nil)
	q.AddLast("three", nil)
	assert.Equal(t, 3, q.Size())
	ans, err := q.RemoveFirst()
	assert.Equal(t, "one", ans.GroupID)
	assert.Nil(t, err)
	assert.Equal(t, 2, q.Size())
	ans, err = q.RemoveFirst()
	assert.Nil(t, err)
	assert.Equal(t, "two", ans.GroupID)
	assert.Equal(t, 1, q.Size())
	ans, err = q.RemoveFirst()
	assert.Nil(t, err)
	assert.Equal(t, "three", ans.GroupID)
	assert.Equal(t, 0, q.Size())
}

func TestQueueRemoveFromEmpty(t *testing.T) {
	q := &Queue{}
	_, err := q.RemoveFirst()
	assert.Error(t, err)
	assert.Equal(t, 0, q.Size())
}

func TestQueueGetFirstGroup(t *testing.T) {
	q := &Queue{}
	q.AddLast("one", nil)
	q.AddLast("two", nil)
	assert.Equal(t, "one", q.FirstGroup())
}

func TestQueueFirstGroupEmpty(t *testing.T) {
	q := &Queue{}
	assert.Equal(t, "", q.FirstGroup())
	assert.Equal(t, 0, q.Size())
}

func TestQueueGetLastGroup(t *testing.T) {
	q := &Queue{}
	q.AddLast("one", nil)
	q.AddLast("two", nil)
	assert.Equal(t, "two", q.LastGroup())
}

func TestQueueLastGroupEmpty(t *testing.T) {
	q := &Queue{}
	assert.Equal(t, "", q.LastGroup())
}

func TestQueueRemoveLast(t *testing.T) {
	q := &Queue{}
	q.AddLast("one", nil)
	q.AddLast("two", nil)
	q.AddLast("three", nil)
	v, err := q.RemoveLast()
	assert.Equal(t, "three", v.GroupID)
	assert.Equal(t, "two", q.last.GroupID)
	assert.Equal(t, 2, q.Size())
	assert.Nil(t, err)
}

func TestQueueRemoveLastOnEmpty(t *testing.T) {
	q := &Queue{}
	v, err := q.RemoveLast()
	assert.Nil(t, v)
	assert.Error(t, err)
}

func TestQueueRemoveLastOnSizeOne(t *testing.T) {
	q := &Queue{}
	q.AddLast("one", nil)
	v, err := q.RemoveLast()
	assert.Equal(t, "one", v.GroupID)
	assert.Nil(t, q.last)
	assert.Nil(t, q.first)
	assert.Nil(t, err)
	assert.Equal(t, 0, q.Size())
}

func TestQueueAddFirst(t *testing.T) {
	q := &Queue{}
	q.AddLast("one", nil)
	q.AddLast("two", nil)
	q.AddLast("three", nil)
	q.AddFirst("zero", nil)
	assert.Equal(t, q.first.GroupID, "zero")
	assert.Equal(t, q.last.GroupID, "three")
	assert.Equal(t, 4, q.Size())
}

func TestQueueAddFirstOnSizeOne(t *testing.T) {
	q := &Queue{}
	q.AddLast("one", nil)
	q.AddFirst("zero", nil)
	assert.Equal(t, q.first.GroupID, "zero")
	assert.Equal(t, q.last.GroupID, "one")
	assert.Equal(t, 2, q.Size())
}

func TestQueueAddFirstOnEmpty(t *testing.T) {
	q := &Queue{}
	q.AddFirst("zero", nil)
	assert.Equal(t, q.first.GroupID, "zero")
	assert.Equal(t, q.last.GroupID, "zero")
	assert.Equal(t, 1, q.Size())
}
