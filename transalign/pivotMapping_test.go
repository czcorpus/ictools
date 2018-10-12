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

package transalign

import (
	"os"
	"path/filepath"

	"github.com/stretchr/testify/assert"

	"testing"
)

func TestInitialization(t *testing.T) {

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	f, err := os.Open(filepath.Join(cwd, "..", "testdata", "foo2.txt"))
	if err != nil {
		panic(err)
	}

	pm, err := NewPivotMapping(f)
	pm.Load()

	assert.Equal(t, 9, pm.Size())
	assert.True(t, pm.HasGapAtRow(8))

}
