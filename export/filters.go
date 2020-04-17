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
	// ExportTypeIntercorp represents a filter for CNC'c InterCorp
	// corpora with their well established rules on how to identify idividual
	// sentences.
	ExportTypeIntercorp = "intercorp"
)

var (
	intercorpPattern = regexp.MustCompile("^(\\w{2}):([\\w\\d_-]+):(\\d+):(\\d+):(\\d+)$")
)

// GroupFilter specifies an object able to extract group identifiers
// (which can be typically seen as a text/document ID) and language information
// out of provided individual record IDs.
// The actual way how to obtain this is not specified here - sometimes
// it can be generates directly from the provided string (e.g. as a substring),
// sometimes the implementation may need an external database.
type GroupFilter interface {

	// ExtractGroupID extract text/doc (aka group here) identifier
	// from the individual record ID
	ExtractGroupID(recID string) string

	// ExtractLangFromRegistry extract language code from a corpus registry
	// path.
	ExtractLangFromRegistry(regPath string) string
}

// ------

// FilterIntercorp is a concrete GroupFilter implementation
// for CNC's InterCorp corpora. The extraction in this case
// can be done just by substringing the identifiers/paths.
type FilterIntercorp struct {
	srch *regexp.Regexp
}

// ExtractGroupID - please see the GroupFilter interface
func (f *FilterIntercorp) ExtractGroupID(recID string) string {
	srch := f.srch.FindStringSubmatch(recID)
	if len(srch) > 0 {
		return srch[2]
	}
	return ""
}

// ExtractLangFromRegistry - please see the GroupFilter interface
func (f *FilterIntercorp) ExtractLangFromRegistry(regPath string) string {
	return regPath[len(regPath)-2:]
}

// ------

// FilterEmpty is a 'null' implementation for a general case where
// we don't know how to extract groups. It means that the export
// in this case produces one big chunk of sentences.
type FilterEmpty struct {
}

// ExtractGroupID - please see the GroupFilter interface
func (f *FilterEmpty) ExtractGroupID(recID string) string {
	return ""
}

// ExtractLangFromRegistry - please see the GroupFilter interface
func (f *FilterEmpty) ExtractLangFromRegistry(regPath string) string {
	return filepath.Base(regPath)
}

// ------

// NewGroupFilter is a factory for GroupFilter instances
func NewGroupFilter(ftype string) GroupFilter {
	switch ftype {
	case ExportTypeIntercorp:
		return &FilterIntercorp{srch: intercorpPattern}
	}
	return &FilterEmpty{}
}
