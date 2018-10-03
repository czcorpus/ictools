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

package calign

import "fmt"

// IgnorableError is represents an error
// in parsing which does not affect resulting
// data - e.g. additional tags in source file
type IgnorableError struct {
	message string
}

func (err IgnorableError) Error() string {
	return err.message
}

// NewIgnorableError creates a new IngorableError instance
// with formatted string (just like fmt.Errorf creates error).
func NewIgnorableError(msg string, args ...interface{}) IgnorableError {
	return IgnorableError{message: fmt.Sprintf(msg, args...)}
}

// -------------------------

// FileImportError is a general error in source XML import.
// It provides a line number where the error was encountered.
type FileImportError struct {
	line    int
	message string
}

func (err FileImportError) Error() string {
	return fmt.Sprintf("%s (line: %d)", err.message, err.line)
}

// NewFileImportError is the default factory function for FileImportError
func NewFileImportError(err error, line int) FileImportError {
	return FileImportError{message: err.Error(), line: line}
}
