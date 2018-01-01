// Copyright 2012 Milos Jakubicek
// Copyright 2017 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2017 Charles University, Faculty of Arts,
//                Institute of the Czech National Corpus
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

package fixgaps

import (
	"bufio"
	"fmt"
	"github.com/czcorpus/ictools/mapping"
	"os"
	"strings"
)

func FixGapsFromFile(file *os.File, onItem func(item mapping.Mapping)) {
	fr := bufio.NewScanner(file)
	lastL1 := -1
	lastL2 := -1
	for fr.Scan() {
		line := fr.Text()
		items := strings.Split(line, "\t")
		l1t := strings.Split(items[0], ",")
		l2t := strings.Split(items[1], ",")
		r1 := mapping.NewPosRange(l1t)
		r2 := mapping.NewPosRange(l2t)
		for r1.First > lastL1+1 {
			lastL1++
			onItem(mapping.NewMapping(lastL1, lastL1, -1, -1))
		}
		for r2.First > lastL2+1 {
			lastL2++
			onItem(mapping.NewMapping(-1, -1, lastL2, lastL2))
		}
		if r1.Last != -1 {
			lastL1 = r1.Last
		}
		if r2.Last != -1 {
			lastL2 = r2.Last
		}
		fmt.Println(line)
	}
}
