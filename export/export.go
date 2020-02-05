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
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/czcorpus/ictools/attrib"
	"github.com/czcorpus/ictools/mapping"
)

func createTag(corp1 attrib.GoCorpus, attr1 attrib.GoPosAttr, corp2 attrib.GoCorpus, attr2 attrib.GoPosAttr, item mapping.Mapping) string {
	var lft, rgt string
	var lftNum, rgtNum int

	if item.From.First == -1 {
		lft = ""
		lftNum = 0

	} else if item.From.First != item.From.Last {
		lft = attr1.ID2Str(item.From.First) + " " + attr1.ID2Str(item.From.Last)
		lftNum = item.From.Last - item.From.First + 1

	} else {
		lft = attr1.ID2Str(item.From.First)
		lftNum = 1
	}

	if item.To.First == -1 {
		rgt = ""
		rgtNum = 0

	} else if item.To.First != item.To.Last {
		rgt = attr2.ID2Str(item.To.First) + " " + attr2.ID2Str(item.To.Last)
		rgtNum = item.To.Last - item.To.First + 1

	} else {
		rgt = attr2.ID2Str(item.From.First)
		rgtNum = 1
	}
	return fmt.Sprintf("<link type=\"%d-%d\" xtargets=\"%s;%s\" status=\"man\" />", lftNum, rgtNum, lft, rgt)
}

func Run(corp1 attrib.GoCorpus, attr1 attrib.GoPosAttr, corp2 attrib.GoCorpus, attr2 attrib.GoPosAttr, mappingPath string) {
	log.Print(corp1, corp2)

	srcFile, err := os.Open(mappingPath)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}
	fmt.Println("<?xml version=\"1.0\" encoding=\"utf-8\"?>")
	fmt.Println("<linkGrp>")
	fr := bufio.NewScanner(srcFile)
	for i := 0; fr.Scan(); i++ {
		item, err := mapping.NewMappingFromString(fr.Text())
		if err != nil {
			log.Print("ERROR: ", err)
		}
		fmt.Println(createTag(corp1, attr1, corp2, attr2, item))
	}
	fmt.Println("</linkGrp>")
}
