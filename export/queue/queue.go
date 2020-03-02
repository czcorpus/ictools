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
	"log"
	"strings"

	"github.com/czcorpus/ictools/mapping"
)

// Element groups items for the Deque.
// It wraps both original numeric mapping
// and also the first string variant of
// a respective group ID. We say 'first'
// because the mapping.Mapping may in
// general represent multiple texts.
//
type Element struct {
	GroupID string
	Mapping *mapping.Mapping
	next    *Element
}

// Deque is a double ended queue as needed
// by Ictools export algorithm.
type Deque struct {
	first *Element
	last  *Element
	size  int
}

// New creates a new Deque instance
func New() *Deque {
	return new(Deque)
}

// Size returns a size of the queue.
func (q *Deque) Size() int {
	return q.size
}

// PushBack adds a new item to the end of the Deque.
// The complexity is O(1).
func (q *Deque) PushBack(groupID string, mp *mapping.Mapping) {
	if groupID == "" {
		log.Fatalf("FATAL: entering empty group for mapping %v", mp)
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

// PopBack removes an item from the back of the
// Deque. This operation's complexity is O(n).
func (q *Deque) PopBack() (*Element, error) {
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

// PushFront adds a new item to the front of the Deque.
// The complexity is O(1).
func (q *Deque) PushFront(groupID string, mp *mapping.Mapping) {
	n := &Element{GroupID: groupID, Mapping: mp, next: q.first}
	q.first = n
	if q.size == 0 {
		q.last = n
	}
	q.size++
}

// PopFront removes an item from the front of the Deque.
// The complexity is O(1).
func (q *Deque) PopFront() (*Element, error) {
	if q.first != nil {
		v := q.first
		q.first = q.first.next
		q.size--
		return v, nil
	}
	return nil, fmt.Errorf("Empty queue")
}

// BackGroup returns a group identifier of the
// item at the back of the Deque. In case of
// an empty Deque, an empty string is returned.
// The complexity is O(1).
func (q *Deque) BackGroup() string {
	if q.last != nil {
		return q.last.GroupID
	}
	return ""
}

// FrontGroup returns a group identifier of the
// item at the front of the Deque. In case of
// an empty Deque, an empty string is returned.
// The complexity is O(1).
func (q *Deque) FrontGroup() string {
	if q.first != nil {
		return q.first.GroupID
	}
	return ""
}

// ForEach applies a provided function on all
// the items starting from the front item.
func (q *Deque) ForEach(fn func(groupID string, m *mapping.Mapping)) {
	curr := q.first
	for ; curr != nil; curr = curr.next {
		fn(curr.GroupID, curr.Mapping)
	}
}

func (q *Deque) String() string {
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
