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
		rgt = attr2.ID2Str(item.To.First)
		rgtNum = 1
	}
	return fmt.Sprintf("<link type=\"%d-%d\" xtargets=\"%s;%s\" status=\"man\" />", lftNum, rgtNum, lft, rgt)
}

func createGroupTag(lang1, lang2, ident string) string {
	g1 := fmt.Sprintf("%s.%s-00.xml", ident, lang1)
	g2 := fmt.Sprintf("%s.%s-00.xml", ident, lang2)
	return fmt.Sprintf("<linkGrp toDoc=\"%s\" fromDoc=\"%s\">", g2, g1)
}

type ExportArgs struct {
	RegPath1        string
	Corp1           attrib.GoCorpus
	Attr1           attrib.GoPosAttr
	RegPath2        string
	Corp2           attrib.GoCorpus
	Attr2           attrib.GoPosAttr
	MappingPath     string
	GroupFilterType string
}

func Run(args ExportArgs) {
	log.Print(args.Corp1, args.Corp2)

	srcFile, err := os.Open(args.MappingPath)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}
	groupFilter := NewGroupFilter(args.GroupFilterType)
	lang1 := groupFilter.ExtractLangFromRegistry(args.RegPath1)
	lang2 := groupFilter.ExtractLangFromRegistry(args.RegPath2)

	fmt.Println("<?xml version=\"1.0\" encoding=\"utf-8\"?>")
	fr := bufio.NewScanner(srcFile)
	var currGroup1, newGroup1 string
	for i := 0; fr.Scan(); i++ {
		item, err := mapping.NewMappingFromString(fr.Text())
		if err != nil {
			log.Print("ERROR: ", err)
		}
		newGroup1 = groupFilter.ExtractGroupId(args.Attr1.ID2Str(item.From.First))
		if newGroup1 == "" {
			newGroup1 = groupFilter.ExtractGroupId(args.Attr2.ID2Str(item.To.First))
		}
		if newGroup1 != "" && newGroup1 != currGroup1 {
			if currGroup1 != "" {
				fmt.Println("</linkGrp>")
			}
			fmt.Println(createGroupTag(lang1, lang2, newGroup1))
			currGroup1 = newGroup1
		}
		fmt.Println(createTag(args.Corp1, args.Attr1, args.Corp2, args.Attr2, item))
	}
	fmt.Println("</linkGrp>")
}
