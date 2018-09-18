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

// Str2Int converts a string-represented integer to int.
// In case of an error the function returns -1 and logs
// and error message.
func Str2Int(v string) int {
	ans, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("ERROR: Failed to import string-encoded integer '%s'", v)
		return -1
	}
	return ans
}

// FileSize returns a size of a specified file.
// In case of an error the function returns
// size -1 and a respective error
func FileSize(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	st, err := f.Stat()
	if err != nil {
		return -1, err
	}
	return int(st.Size()), nil
}
