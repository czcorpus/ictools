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
	"strings"

	"github.com/czcorpus/ictools/attrib"
	"github.com/czcorpus/ictools/export/gpool"
	"github.com/czcorpus/ictools/mapping"
)

func createTag(corp1 attrib.GoCorpus, attr1 attrib.GoPosAttr, corp2 attrib.GoCorpus, attr2 attrib.GoPosAttr, item *mapping.Mapping) string {
	var lft, rgt string
	var lftNum, rgtNum int

	if item.From.First == -1 {
		lft = ""
		lftNum = 0

	} else if item.From.First != item.From.Last {
		lftNum = item.From.Last - item.From.First + 1

		if item.To.First == -1 {
			lft = attr1.ID2Str(item.From.First) + " " + attr1.ID2Str(item.From.Last)

		} else {
			items := make([]string, lftNum)
			for i := 0; i < len(items); i++ {
				items[i] = attr1.ID2Str(item.From.First + i)
			}
			lft = strings.Join(items, " ")
		}

	} else {
		lft = attr1.ID2Str(item.From.First)
		lftNum = 1
	}

	if item.To.First == -1 {
		rgt = ""
		rgtNum = 0

	} else if item.To.First != item.To.Last {
		rgtNum = item.To.Last - item.To.First + 1
		if item.From.First == -1 {
			rgt = attr2.ID2Str(item.To.First) + " " + attr2.ID2Str(item.To.Last)

		} else {
			items := make([]string, rgtNum)
			for i := 0; i < len(items); i++ {
				items[i] = attr2.ID2Str(item.To.First + i)
			}
			rgt = strings.Join(items, " ")
		}

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

// RunArgs wraps all the values needed to configure the export.
type RunArgs struct {
	RegPath1        string
	Corp1           attrib.GoCorpus
	Attr1           attrib.GoPosAttr
	RegPath2        string
	Corp2           attrib.GoCorpus
	Attr2           attrib.GoPosAttr
	MappingPath     string
	GroupFilterType string
}

// ungroupAndAdd  ungroups (if needed) items encoded in a numeric interval specified by "item".
// All the resulting groupIDs and *mapping.Mapping instances are then added to the pool
func ungroupAndAdd(item *mapping.Mapping, groupFilter GroupFilter, pool *gpool.TextGroupPool, attr1 attrib.GoPosAttr, attr2 attrib.GoPosAttr) {
	var newGroup, currGroup string
	if item.From.First == -1 {
		currGroupStartIdx := item.To.First
		for i := item.To.First; i <= item.To.Last; i++ {
			newGroup = groupFilter.ExtractGroupID(attr2.ID2Str(i))
			if newGroup != "" {
				if currGroup != "" && newGroup != currGroup {
					pool.AddGroup(currGroup, &mapping.Mapping{
						From: mapping.PosRange{First: -1, Last: -1},
						To:   mapping.PosRange{First: currGroupStartIdx, Last: i - 1},
					})
					currGroupStartIdx = i
				}
				currGroup = newGroup
			}
		}
		if newGroup != "" {
			pool.AddGroup(newGroup, &mapping.Mapping{
				From: mapping.PosRange{First: -1, Last: -1},
				To:   mapping.PosRange{First: currGroupStartIdx, Last: item.To.Last},
			})
		}

	} else if item.To.First == -1 {
		currGroupStartIdx := item.From.First
		for i := item.From.First; i <= item.From.Last; i++ {
			newGroup = groupFilter.ExtractGroupID(attr1.ID2Str(i))
			if newGroup != "" {
				if currGroup != "" && newGroup != currGroup {
					pool.AddGroup(currGroup, &mapping.Mapping{
						From: mapping.PosRange{First: currGroupStartIdx, Last: i - 1},
						To:   mapping.PosRange{First: -1, Last: -1},
					})
					currGroupStartIdx = i
				}
				currGroup = newGroup
			}
		}
		if newGroup != "" {
			pool.AddGroup(newGroup, &mapping.Mapping{
				From: mapping.PosRange{First: currGroupStartIdx, Last: item.From.Last},
				To:   mapping.PosRange{First: -1, Last: -1},
			})
		}

	} else {
		currGroup = groupFilter.ExtractGroupID(attr1.ID2Str(item.From.First))
		pool.AddGroup(currGroup, item)
	}
}

// getGroupIdent extracts a respective string identifier either from attr1 (i.e. first language)
// or attr2 (i.e. the second language).
func getGroupIdent(item *mapping.Mapping, groupFilter GroupFilter, attr1 attrib.GoPosAttr, attr2 attrib.GoPosAttr) string {
	var group string
	if item.From.First != -1 {
		group = groupFilter.ExtractGroupID(attr1.ID2Str(item.From.First))
	}
	if group == "" && item.To.First != -1 {
		group = groupFilter.ExtractGroupID(attr2.ID2Str(item.To.First))
	}
	return group
}

// Run generates a XML-ish output with the same format as the one
// used as input format for generating numerical alignment files.
// The algorithm is able to ungroup 'compressed' numeric intervals
// so if an interval contains multiple texts - all of them should
// be written to the output.
func Run(args RunArgs) {
	srcFile, err := os.Open(args.MappingPath)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}
	groupFilter := NewGroupFilter(args.GroupFilterType)
	lang1 := groupFilter.ExtractLangFromRegistry(args.RegPath1)
	lang2 := groupFilter.ExtractLangFromRegistry(args.RegPath2)

	fmt.Println("<?xml version=\"1.0\" encoding=\"utf-8\"?>")
	fr := bufio.NewScanner(srcFile)
	var newGroup1 string
	currGroups := gpool.NewTextGroupPool()
	for i := 0; fr.Scan(); i++ {
		item, err := mapping.NewMappingFromString(fr.Text())
		if err != nil {
			log.Print("ERROR: ", err)
		}
		newGroup1 = getGroupIdent(&item, groupFilter, args.Attr1, args.Attr2)
		if newGroup1 != "" {
			ungroupAndAdd(&item, groupFilter, currGroups, args.Attr1, args.Attr2)
			for nxt := currGroups.PopNextReady(); nxt != nil; nxt = currGroups.PopNextReady() {
				log.Print("NXT: ", nxt)
				fmt.Println(createGroupTag(lang1, lang2, nxt.ID))
				nxt.ForEach(func(mp *mapping.Mapping) {
					fmt.Println(createTag(args.Corp1, args.Attr1, args.Corp2, args.Attr2, mp))
				})
				fmt.Println("</linkGrp>")
			}
		}
	}
}
