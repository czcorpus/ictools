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

// Package attrib contains wrapper functions and types
// used to access Manatee C library.
package attrib

// #cgo CFLAGS: -I${SRCDIR}/attrib -I${SRCDIR}/attrib/corp
// #cgo LDFLAGS:  -lmanatee -L${SRCDIR} -Wl,-rpath='$ORIGIN'
// #include "attrib.h"
import "C"

type GoCorpus struct {
	corp C.CorpusV
}

func GetStructSize(corpus GoCorpus, name string) int {
	return int(C.get_struct_size(corpus.corp, C.CString(name)))
}

// ---

type GoPosAttr struct {
	attr C.PosAttrV
}

func (gpa GoPosAttr) Str2ID(value string) int {
	return int(C.attr_str2id(gpa.attr, C.CString(value)))
}

func OpenCorpus(path string) GoCorpus {
	ret := GoCorpus{}
	ret.corp = C.open_corpus(C.CString(path))
	return ret
}

func OpenAttr(corpus GoCorpus, name string) GoPosAttr {
	ret := GoPosAttr{}
	ret.attr = C.get_attr(corpus.corp, C.CString(name))
	return ret
}
