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

package gpool

import (
	"github.com/czcorpus/ictools/mapping"
)

// TextGroupPool is a pool of gradually built text groups
type TextGroupPool struct {
	data             map[string]*TextGroup
	lastGroup        string
	numGroupSwitches int
}

// AddGroup adds a new (or an already existing) group to the pool
func (tgp *TextGroupPool) AddGroup(groupID string, m *mapping.Mapping) {
	if groupID != tgp.lastGroup {
		tgp.numGroupSwitches++
		tgp.lastGroup = groupID
	}
	g, ok := tgp.data[groupID]
	if ok {
		g.mappings = append(g.mappings, m)
		g.stepLast = tgp.numGroupSwitches

	} else {
		tgp.data[groupID] = NewTextGroup(groupID, m, tgp.numGroupSwitches)
	}
}

// PopNextReady removes and returns the oldest text group which last
// change is more than 3 group changes old (which should be
// OK for how ictools generate numeric alignments).
//
// In case no group matches the crieria, nil is returned
func (tgp *TextGroupPool) PopNextReady() *TextGroup {
	minFound := tgp.numGroupSwitches
	var minKey string
	for k, v := range tgp.data {
		if tgp.numGroupSwitches-v.stepLast >= 3 && v.stepFound <= minFound {
			minFound = v.stepFound
			minKey = k
		}
	}
	if minKey != "" {
		ans := tgp.data[minKey]
		delete(tgp.data, minKey)
		return ans
	}
	return nil
}

// Size returns number of text groups
func (tgp *TextGroupPool) Size() int {
	return len(tgp.data)
}

// NewTextGroupPool is a factory function for creating a new pool
func NewTextGroupPool() *TextGroupPool {
	return &TextGroupPool{
		data:             make(map[string]*TextGroup),
		numGroupSwitches: -1,
	}
}
