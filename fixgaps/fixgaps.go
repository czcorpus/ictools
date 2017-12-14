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
	"github.com/czcorpus/ictools/common"
	"os"
	"strings"
)

func FixGaps(file *os.File) {
	fr := bufio.NewScanner(file)
	lastL1 := -1
	lastL2 := -1
	for fr.Scan() {
		line := fr.Text()
		items := strings.Split(line, "\t")
		l1t := strings.Split(items[0], ",")
		l2t := strings.Split(items[1], ",")
		l11 := common.Str2Int(l1t[0])
		l12 := common.Str2Int(l1t[len(l1t)-1])
		l21 := common.Str2Int(l2t[0])
		l22 := common.Str2Int(l2t[len(l2t)-1])
		for l11 > lastL1+1 {
			lastL1++
			fmt.Printf("%d\t-1\n", lastL1)
		}
		for l21 > lastL2+1 {
			lastL2++
			fmt.Printf("-1\t%d\n", lastL2)
		}
		if l12 != -1 {
			lastL1 = l12
		}
		if l22 != -1 {
			lastL2 = l22
		}
		fmt.Println(line)
	}
}
