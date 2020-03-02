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
	"fmt"
	"strings"

	"github.com/czcorpus/ictools/mapping"
)

type Element struct {
	GroupID string
	Mapping *mapping.Mapping
	next    *Element
}

type Queue struct {
	first *Element
	last  *Element
	size  int
}

func New() *Queue {
	return new(Queue)
}

func (q *Queue) Size() int {
	return q.size
}

func (q *Queue) AddLast(groupID string, mp *mapping.Mapping) {
	if groupID == "" {
		panic("")
	}
	n := &Element{GroupID: groupID, Mapping: mp}
	if q.first == nil {
		q.first = n
		q.last = n

	} else {
		q.last.next = n
		q.last = n
	}
	q.size++
}

func (q *Queue) AddFirst(groupID string, mp *mapping.Mapping) {
	n := &Element{GroupID: groupID, Mapping: mp, next: q.first}
	q.first = n
	if q.size == 0 {
		q.last = n
	}
	q.size++
}

func (q *Queue) LastGroup() string {
	if q.last != nil {
		return q.last.GroupID
	}
	return ""
}

func (q *Queue) FirstGroup() string {
	if q.first != nil {
		return q.first.GroupID
	}
	return ""
}

func (q *Queue) RemoveFirst() (*Element, error) {
	if q.first != nil {
		v := q.first
		q.first = q.first.next
		q.size--
		return v, nil
	}
	return nil, fmt.Errorf("Empty queue")
}

func (q *Queue) RemoveLast() (*Element, error) {
	if q.first == nil {
		return nil, fmt.Errorf("Empty queue")
	}
	var prev *Element
	var curr *Element
	for curr = q.first; curr.next != nil; curr = curr.next {
		prev = curr
	}
	if prev != nil {
		prev.next = nil
		q.last = prev

	} else {
		q.first = nil
		q.last = nil
	}
	q.size--
	return curr, nil
}

func (q *Queue) ForEach(fn func(groupID string, m *mapping.Mapping)) {
	curr := q.first
	for ; curr != nil; curr = curr.next {
		fn(curr.GroupID, curr.Mapping)
	}
}

func (q *Queue) String() string {
	if q.size == 0 {
		return "Queue []"
	}
	tmp := make([]string, 0, 10)
	curr := q.first
	i := 0
	for i = 0; i < 4 && curr != nil; i++ {
		tmp = append(tmp, curr.GroupID)
		curr = curr.next
	}
	i = q.size - 2
	if i < 10 {
		tmp = append(tmp, "...")
		i++
		for ; i < 10 && curr != nil; i++ {
			tmp = append(tmp, curr.GroupID)
			curr = curr.next
		}
	}
	return fmt.Sprintf("Queue (len=%d) [%s]", q.size, strings.Join(tmp, ", "))
}
