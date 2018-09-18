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

// #cgo CXXFLAGS: -std=c++11
// #cgo CFLAGS: -I${SRCDIR}/attrib -I${SRCDIR}/attrib/corp
// #cgo LDFLAGS:  -lmanatee -L${SRCDIR} -Wl,-rpath='$ORIGIN'
// #include <stdlib.h>
// #include "attrib.h"
import "C"

import (
	"fmt"
	"unsafe"
)

// GoCorpus is a wrapper for Manatee Corpus instance
type GoCorpus struct {
	corp C.CorpusV
}

// GetStructSize returns a number of occurences
// for a specific structure in a provided corpus.
func GetStructSize(corpus GoCorpus, name string) (int, error) {
	ans := (C.get_struct_size(corpus.corp, C.CString(name)))
	if ans.err != nil {
		err := fmt.Errorf(C.GoString(ans.err))
		defer C.free(unsafe.Pointer(ans.err))
		return -1, err
	}
	return int(ans.value), nil
}

// GoPosAttr is a wrapper for Manatee PosAttr
// (note: structural attributes belong here too)
type GoPosAttr struct {
	attr C.PosAttrV
}

// Str2ID transforms a string value of the attribute
// to its numeric form (= an index).
func (gpa GoPosAttr) Str2ID(value string) int {
	return int(C.attr_str2id(gpa.attr, C.CString(value)))
}

// OpenCorpus is a factory function creating
// a Manatee corpus wrapper.
func OpenCorpus(path string) (GoCorpus, error) {
	ret := GoCorpus{}
	var err error
	ans := C.open_corpus(C.CString(path))
	if ans.err != nil {
		err = fmt.Errorf(C.GoString(ans.err))
		defer C.free(unsafe.Pointer(ans.err))
		return ret, err
	}
	ret.corp = ans.value
	return ret, nil
}

// OpenAttr gets an instance of a specific
// structural attribute of a provided corpus.
func OpenAttr(corpus GoCorpus, name string) (GoPosAttr, error) {
	ret := GoPosAttr{}
	var err error
	ans := C.get_attr(corpus.corp, C.CString(name))
	if ans.err != nil {
		err = fmt.Errorf(C.GoString(ans.err))
		defer C.free(unsafe.Pointer(ans.err))
		return ret, err
	}
	ret.attr = ans.value
	return ret, nil
}
