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

package common

import (
	"log"
	"os"
	"strconv"
)

func Str2Int(v string) int {
	ans, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("Failed to import string-encoded integer '%s'", v)
		return -1
	}
	return ans
}

func FileSize(path string) int {
	f, err := os.Open(path)
	if err != nil {
		return -1
	}
	st, err := f.Stat()
	if err != nil {
		log.Printf("Failed to get file info: %s", err)
	}
	return int(st.Size())
}
