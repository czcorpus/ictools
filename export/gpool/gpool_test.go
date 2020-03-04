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
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestStoresProperReference(t *testing.T) {
	p := NewTextGroupPool()
	m1 := &mapping.Mapping{}
	m2 := &mapping.Mapping{}
	p.AddGroup("one", m1)
	p.AddGroup("two", m2)
	assert.Same(t, m1, p.data["one"].mappings[0])
	assert.Same(t, m2, p.data["two"].mappings[0])
}

func TestKeepsSameGroupTogether(t *testing.T) {
	p := NewTextGroupPool()
	m1 := &mapping.Mapping{}
	m2 := &mapping.Mapping{}
	p.AddGroup("one", m1)
	p.AddGroup("one", m2)
	assert.Same(t, m1, p.data["one"].mappings[0])
	assert.Same(t, m2, p.data["one"].mappings[1])
	assert.Equal(t, 2, len(p.data["one"].mappings))
}

func TestPopNextReady(t *testing.T) {
	p := NewTextGroupPool()
	p.AddGroup("one", nil)
	p.AddGroup("two", nil)
	p.AddGroup("three", nil)
	p.AddGroup("four", nil)
	nxt := p.PopNextReady()
	assert.Equal(t, "one", nxt.ID)
}

func TestPopNextReadyMustBeFourOrMore(t *testing.T) {
	p := NewTextGroupPool()
	p.AddGroup("one", nil)
	p.AddGroup("two", nil)
	p.AddGroup("three", nil)
	nxt := p.PopNextReady()
	assert.Nil(t, nxt)
}

func TestPopNextReadyOnEmpty(t *testing.T) {
	p := NewTextGroupPool()
	nxt := p.PopNextReady()
	assert.Nil(t, nxt)
}

func TestPopOldest(t *testing.T) {
	p := NewTextGroupPool()
	p.AddGroup("one", nil)
	p.AddGroup("two", nil)
	nxt := p.PopOldest()
	assert.Equal(t, "one", nxt.ID)
}

func TestPopOldestOnEmpty(t *testing.T) {
	p := NewTextGroupPool()
	nxt := p.PopOldest()
	assert.Nil(t, nxt)
}
