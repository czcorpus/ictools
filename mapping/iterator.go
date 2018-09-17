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

package mapping

type ProcPosition struct {
	Left  int
	Right int
}

// Iterator is used when merging two sorted mappings together.
// It provides a way how to apply a function to each item rather
// than exposing the item.

type Iterator struct {
	mapping  []Mapping
	currIdx  int
	pos      *ProcPosition
	finished bool
}

// NewIterator creates a new Iterator instance
func NewIterator(data []Mapping, pos *ProcPosition) Iterator {
	finished := false
	if len(data) == 0 {
		finished = true
	}
	return Iterator{
		mapping:  data,
		currIdx:  0,
		pos:      pos,
		finished: finished,
	}
}

// Apply calls a provided function with the current
// item as its argument. After the method is called,
// a possible "finished" state is .
func (m *Iterator) Apply(onItem func(item Mapping, pos *ProcPosition)) {
	onItem(m.mapping[m.currIdx], m.pos)
	if m.mapping[m.currIdx].From.First != -1 {
		m.pos.Left = m.mapping[m.currIdx].From.Last
	}
	if m.mapping[m.currIdx].To.First != -1 {
		m.pos.Right = m.mapping[m.currIdx].To.Last
	}
	if m.currIdx == len(m.mapping)-1 {
		m.finished = true
	}
}

// HasPriorityOver compares latest items of two iterators
// and returns true if the item from the first one is
// less then (see how LessThan is defined on PosRange)
// the second one.
func (m *Iterator) HasPriorityOver(m2 *Iterator) bool {
	return !m.finished && m.mapping[m.currIdx].To.LessThan(m2.mapping[m2.currIdx].To)
}

// Next moves an internal index to the next item.
// In case the index reached the end, nothing is
// done.
func (m *Iterator) Next() {
	if m.currIdx < len(m.mapping)-1 {
		m.currIdx++
	}
}

// Unfinished tells whether Next() will
// provide another item.
func (m *Iterator) Unfinished() bool {
	return !m.finished
}
