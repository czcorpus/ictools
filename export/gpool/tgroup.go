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
	"fmt"

	"github.com/czcorpus/ictools/mapping"
)

// TextGroup represents a reconstructed list of
// items (~ sentences) belonging to a single
// group (~ document).
type TextGroup struct {
	ID        string
	mappings  []*mapping.Mapping
	stepFound int
	stepLast  int
}

func (tg *TextGroup) String() string {
	return fmt.Sprintf("TextGroup [ID: %v, StepFound: %d, StepLast: %d, Num of mappings: %v", tg.ID, tg.stepFound, tg.stepLast, len(tg.mappings))
}

// ForEach applies a function for all the mappings in the group
func (tg *TextGroup) ForEach(fn func(mp *mapping.Mapping)) {
	for _, v := range tg.mappings {
		fn(v)
	}
}

// NewTextGroup is a factory for creating a text group instance
func NewTextGroup(ID string, m *mapping.Mapping, stepFound int) *TextGroup {
	mlist := make([]*mapping.Mapping, 1, 3000)
	mlist[0] = m
	return &TextGroup{
		ID:        ID,
		mappings:  mlist,
		stepFound: stepFound,
		stepLast:  stepFound,
	}
}
