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

package attrib

// #cgo CFLAGS: -I${SRCDIR}/attrib -I${SRCDIR}/attrib/corp
// #cgo LDFLAGS:  -lmanatee -L${SRCDIR} -Wl,-rpath='$ORIGIN'
// #include "attrib.h"
import "C"

type GoPosAttr struct {
	attr C.PosAttrV
}

func (gpa GoPosAttr) Str2ID(value string) int64 {
	return int64(C.attr_str2id(gpa.attr, C.CString(value)))
}

type GoCorpus struct {
	corp C.CorpusV
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
