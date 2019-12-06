// Copyright 2018 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2018 Charles University, Faculty of Arts,
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
	"os"
	"path/filepath"

	"github.com/czcorpus/ictools/mapping"

	"github.com/stretchr/testify/assert"

	"testing"
)

type MockAttr1 struct {
}

func (ma *MockAttr1) Str2ID(value string) int {
	switch value {
	case "foo:0":
		return 0
	case "foo:1":
		return 1
	case "foo:2":
		return 2
	case "foo:3":
		return 3
	case "foo:4":
		return 4
	case "foo:5":
		return 5
	default:
		return -1
	}
}

func (ma *MockAttr1) ID2Str(ident int) string {
	switch ident {
	case 0:
		return "foo:0"
	case 1:
		return "foo:1"
	case 2:
		return "foo:2"
	case 3:
		return "foo:3"
	case 4:
		return "foo:4"
	case 5:
		return "foo:5"
	default:
		return ""
	}
}

type MockAttr2 struct {
}

func (ma *MockAttr2) Str2ID(value string) int {
	switch value {
	case "bar:0":
		return 0
	case "bar:1":
		return 1
	case "bar:2":
		return 2
	case "bar:3":
		return 3
	case "bar:4":
		return 4
	case "bar:5":
		return 5
	default:
		return -1
	}
}

func (ma *MockAttr2) ID2Str(ident int) string {
	switch ident {
	case 0:
		return "bar:0"
	case 1:
		return "bar:1"
	case 2:
		return "bar:2"
	case 3:
		return "bar:3"
	case 4:
		return "bar:4"
	case 5:
		return "bar:5"
	default:
		return ""
	}
}

func createProcessor() *Processor {
	return &Processor{
		attr1: &MockAttr1{},
		attr2: &MockAttr2{},
	}
}

func createFullProcessor() *Processor {
	ans := createProcessor()
	ans.valPrefix = "xtargets='"
	ans.valSuffix = "'"
	ans.valOffset = len(ans.valPrefix)
	return ans
}

func TestProcessColElementSingle(t *testing.T) {
	p := createProcessor()
	r, err := p.processColElm("foo:0", p.attr1, 0)
	assert.Nil(t, err)
	assert.Equal(t, 0, r.First)
}

func TestProcessColElementRange(t *testing.T) {
	p := createProcessor()
	r, err := p.processColElm("foo:0 foo:3", p.attr1, 0)
	assert.Nil(t, err)
	assert.Equal(t, 0, r.First)
	assert.Equal(t, 3, r.Last)
}

func TestProcessColElementBadSyntax(t *testing.T) {
	p := createProcessor()
	r, err := p.processColElm("foo:0-foo:3", p.attr1, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, r.First)
	assert.Equal(t, 0, r.Last)
}

func TestProcessColElementNonExistent(t *testing.T) {
	p := createProcessor()
	r, err := p.processColElm("foo:123", p.attr1, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, r.First)
	assert.Equal(t, 0, r.Last)
}

// TestProcessColElementNonExistentRightHalf
// the function should auto-correct
func TestProcessColElementNonExistentRightHalf(t *testing.T) {
	p := createProcessor()
	r, err := p.processColElm("foo:1 foo:20", p.attr1, 0)
	assert.Nil(t, err)
	assert.Equal(t, 1, r.First)
	assert.Equal(t, 1, r.Last) // we set num value of foo:1 here
}

func TestProcessColElementNonExistentLeftHalf(t *testing.T) {
	p := createProcessor()
	r, err := p.processColElm("foo:20 foo:2", p.attr1, 0)
	assert.Nil(t, err)
	assert.Equal(t, 2, r.First) // we set num value of foo:2 here
	assert.Equal(t, 2, r.Last)
}

func TestProcessColElementNonExistentBothSides(t *testing.T) {
	p := createProcessor()
	r, err := p.processColElm("foo:20 foo:21", p.attr1, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, r.First)
	assert.Equal(t, 0, r.Last)
}

func TestParseLine(t *testing.T) {
	line := "<link type='1-1' xtargets='pl:_ACQUIS:jrc2;cs:_ACQUIS:jrc3' status='auto'/>"
	p := createProcessor()
	p.valPrefix = "ets='"
	p.valSuffix = "'"
	p.valOffset = len("ets='")
	v := p.parseLine(line)
	assert.Equal(t, "pl:_ACQUIS:jrc2;cs:_ACQUIS:jrc3", v)
}

func TestParseLineInvalid(t *testing.T) {
	line := "<link type='1-1' xstuff='pl:_ACQUIS:jrc2;cs:_ACQUIS:jrc3' status='auto'/>"
	p := createProcessor()
	p.valPrefix = "ets='"
	p.valSuffix = "'"
	p.valOffset = len("ets='")
	v := p.parseLine(line)
	assert.Equal(t, "", v)
}

func TestProcessLine(t *testing.T) {
	line := "<link type='1-1' xtargets='foo:1 foo:2;bar:1 bar:3' status='auto'/>"
	p := createFullProcessor()
	m, err := p.processLine(line, 0)
	assert.Nil(t, err)
	assert.Equal(t, 1, m.From.First)
	assert.Equal(t, 2, m.From.Last)
	assert.Equal(t, 1, m.To.First)
	assert.Equal(t, 3, m.To.Last)
}

func TestProcessLineNonAlignmentLine(t *testing.T) {
	line := "<linkGrp>"
	p := createFullProcessor()
	m, err := p.processLine(line, 0)
	assert.IsType(t, IgnorableError{}, err)
	assert.Equal(t, 0, m.From.First)
	assert.Equal(t, 0, m.From.Last)
	assert.Equal(t, 0, m.To.First)
	assert.Equal(t, 0, m.To.Last)
}

func TestProcessLineTooManyItems(t *testing.T) {
	line := "<foo xtargets='foo;and;bar'>"
	p := createFullProcessor()
	m, err := p.processLine(line, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, m.From.First)
	assert.Equal(t, 0, m.From.Last)
	assert.Equal(t, 0, m.To.First)
	assert.Equal(t, 0, m.To.Last)
}

func TestProcessLineSimpleMissing1(t *testing.T) {
	line := "<link type='1-1' xtargets='foo:112 foo:112;bar:1 bar:3' status='auto'/>"
	p := createFullProcessor()
	m, err := p.processLine(line, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, m.From.First)
	assert.Equal(t, 0, m.From.Last)
	assert.Equal(t, 0, m.To.First)
	assert.Equal(t, 0, m.To.Last)
}

func TestProcessLineRangeMissing1(t *testing.T) {
	line := "<link type='1-1' xtargets='foo:112 foo:113;bar:1 bar:3' status='auto'/>"
	p := createFullProcessor()
	m, err := p.processLine(line, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, m.From.First)
	assert.Equal(t, 0, m.From.Last)
	assert.Equal(t, 0, m.To.First)
	assert.Equal(t, 0, m.To.Last)
}

func TestProcessLineSimpleMissing2(t *testing.T) {
	line := "<link type='1-1' xtargets='foo:1 foo:1;bar:293 bar:293' status='auto'/>"
	p := createFullProcessor()
	m, err := p.processLine(line, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, m.From.First)
	assert.Equal(t, 0, m.From.Last)
	assert.Equal(t, 0, m.To.First)
	assert.Equal(t, 0, m.To.Last)
}

func TestProcessLineRangeMissing2(t *testing.T) {
	line := "<link type='1-1' xtargets='foo:1 foo:1;bar:293 bar:294' status='auto'/>"
	p := createFullProcessor()
	m, err := p.processLine(line, 0)
	assert.Error(t, err)
	assert.Equal(t, 0, m.From.First)
	assert.Equal(t, 0, m.From.Last)
	assert.Equal(t, 0, m.To.First)
	assert.Equal(t, 0, m.To.Last)
}

func TestNewProcessor(t *testing.T) {
	attr1 := &MockAttr1{}
	attr2 := &MockAttr2{}
	p := NewProcessor(attr1, attr2, quoteStyleSingle)
	assert.Equal(t, p.valPrefix, "xtargets='")
	assert.Equal(t, p.valSuffix, "'")
	assert.Equal(t, p.valOffset, len("xtargets='"))
	assert.Equal(t, p.lastPos, 0)
	assert.Equal(t, p.lastPivotPos, 0)
	assert.Equal(t, p.attr1, attr1)
	assert.Equal(t, p.attr2, attr2)
}

func TestNewProcessorDoubleQ(t *testing.T) {
	attr1 := &MockAttr1{}
	attr2 := &MockAttr2{}
	p := NewProcessor(attr1, attr2, quoteStyleDouble)
	assert.Equal(t, p.valPrefix, "xtargets=\"")
	assert.Equal(t, p.valSuffix, "\"")
	assert.Equal(t, p.valOffset, len("xtargets=\""))
	assert.Equal(t, p.lastPos, 0)
	assert.Equal(t, p.lastPivotPos, 0)
	assert.Equal(t, p.attr1, attr1)
	assert.Equal(t, p.attr2, attr2)
}

func TestProcessFile(t *testing.T) {
	valData := []mapping.Mapping{
		mapping.NewMapping(0, 0, 0, 0),
		mapping.NewMapping(1, 1, 1, 1),
		mapping.NewMapping(2, 3, 2, 2),
		mapping.NewMapping(-1, -1, 3, 3),
		mapping.NewMapping(4, 4, 4, 5),
		mapping.NewMapping(5, 5, -1, -1),
		mapping.NewGapMapping(-1, -1, 6, 9),
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f, err := os.Open(filepath.Join(cwd, "..", "testdata", "foo-ids.xml"))
	if err != nil {
		panic(err)
	}
	p := createFullProcessor()
	i2 := 0 // we deliberately ignore argument 'i' as possibly flawed
	err = p.ProcessFile(f, 1000, func(item mapping.Mapping, i int) {
		assert.Equal(t, valData[i2], item)
		i2++
	})
	assert.Nil(t, err)
}

func TestProcessFileError(t *testing.T) {

	valData := []mapping.Mapping{
		mapping.NewMapping(0, 0, 0, 0),
		mapping.NewMapping(2, 2, 2, 2),
	}

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f, err := os.Open(filepath.Join(cwd, "..", "testdata", "foo-ids.err.xml"))
	if err != nil {
		panic(err)
	}
	p := createFullProcessor()
	i2 := 0
	p.ProcessFile(f, 1000, func(item mapping.Mapping, i int) {
		assert.Equal(t, valData[i2], item)
		i2++
	})
	assert.Nil(t, err)
}
