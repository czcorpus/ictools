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

package export

import (
	"path/filepath"
	"regexp"
)

const (
	GroupFilterTypeIntercorp = "intercorp"
)

var (
	intercorpPattern = regexp.MustCompile("^(\\w{2}):([\\w\\d_-]+):(\\d+):(\\d+):(\\d+)$")
)

type GroupFilter interface {
	ExtractGroupId(recId string) string
	ExtractLangFromRegistry(regPath string) string
}

// ------

type FilterIntercorp struct {
	srch *regexp.Regexp
}

func (f *FilterIntercorp) ExtractGroupId(recId string) string {
	srch := f.srch.FindStringSubmatch(recId)
	if len(srch) > 0 {
		return srch[2]
	}
	return ""
}

func (f *FilterIntercorp) ExtractLangFromRegistry(regPath string) string {
	return regPath[len(regPath)-2:]
}

// ------

type FilterEmpty struct {
}

func (f *FilterEmpty) ExtractGroupId(recId string) string {
	return ""
}

func (f *FilterEmpty) ExtractLangFromRegistry(regPath string) string {
	return filepath.Base(regPath)
}

// ------

func NewGroupFilter(ftype string) GroupFilter {
	switch ftype {
	case GroupFilterTypeIntercorp:
		return &FilterIntercorp{srch: intercorpPattern}
	}
	return &FilterEmpty{}
}
