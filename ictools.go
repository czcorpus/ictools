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

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/czcorpus/ictools/attrib"
	"github.com/czcorpus/ictools/calign"
	"github.com/czcorpus/ictools/fixgaps"
)

func runCalign(registryPath1 string, registryPath2 string, attrName string, mappingFilePath string) {
	c1 := attrib.OpenCorpus(registryPath1)
	attr1 := attrib.OpenAttr(c1, attrName)
	c2 := attrib.OpenCorpus(registryPath2)
	attr2 := attrib.OpenAttr(c2, attrName)

	var file *os.File
	var err error
	if mappingFilePath == "" {
		file = os.Stdin

	} else {
		file, err = os.Open(mappingFilePath)
		if err != nil {
			panic(fmt.Sprintf("Failed to open file %s", mappingFilePath))
		}
	}
	processor := calign.NewProcessor(attr1, attr2)
	processor.ProcessFile(file)
}

func runFixGaps(filePath string) {
	var file *os.File
	var err error
	if filePath == "" {
		file = os.Stdin

	} else {
		file, err = os.Open(filePath)
		if err != nil {
			panic(fmt.Sprintf("Failed to open file %s", filePath))
		}
	}
	fixgaps.FixGaps(file)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s [options] calign [registry path 1] [registry path 2] [attr] [mapping file]?\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Missing action, try -h for help")
		os.Exit(1)

	} else {
		switch flag.Arg(0) {
		case "calign":
			runCalign(flag.Arg(1), flag.Arg(2), flag.Arg(3), flag.Arg(4))
		case "fixgaps":
			runFixGaps(flag.Arg(1))
		}
	}
}
