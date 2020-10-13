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

func createGroupTag(lang1, lang2, ident string) string {
	g1 := fmt.Sprintf("%s.%s-00.xml", ident, lang1)
	g2 := fmt.Sprintf("%s.%s-00.xml", ident, lang2)
	return fmt.Sprintf("<linkGrp toDoc=\"%s\" fromDoc=\"%s\">", g2, g1)
}

type Export struct {
	RegPath1    string
	Corp1       attrib.GoCorpus
	Attr1       attrib.GoPosAttr
	RegPath2    string
	Corp2       attrib.GoCorpus
	Attr2       attrib.GoPosAttr
	MappingPath string
	groupFilter GroupFilter
	pool        *gpool.TextGroupPool
}

func (e *Export) createPosRange(rng *mapping.PosRange, attr attrib.GoPosAttr, itemize bool) []string {
	if itemize {
		items := make([]string, rng.Last-rng.First+1)
		for i := 0; i < len(items); i++ {
			items[i] = attr.ID2Str(rng.First + i)
		}
		return items
	}
	return []string{attr.ID2Str(rng.First), attr.ID2Str(rng.Last)}
}

func (e *Export) createTag(item *mapping.Mapping, exportType string) []string {
	var lft, rgt []string
	var lftArity, rgtArity int

	if item.From.First == -1 {
		lft = []string{}
		lftArity = 0

	} else if item.From.First != item.From.Last {
		lftArity = item.From.Last - item.From.First + 1

		if item.To.First == -1 {
			lft = e.createPosRange(&item.From, e.Attr1, exportType == "intercorp")

		} else {
			lft = e.createPosRange(&item.From, e.Attr1, false)
		}

	} else {
		lft = []string{e.Attr1.ID2Str(item.From.First)}
		lftArity = 1
	}

	if item.To.First == -1 {
		rgt = []string{}
		rgtArity = 0

	} else if item.To.First != item.To.Last {
		rgtArity = item.To.Last - item.To.First + 1

		if item.From.First == -1 {
			rgt = e.createPosRange(&item.To, e.Attr2, exportType == "intercorp")

		} else {
			rgt = e.createPosRange(&item.To, e.Attr2, false)
		}

	} else {
		rgt = []string{e.Attr2.ID2Str(item.To.First)}
		rgtArity = 1
	}

	if item.From.First == -1 {
		ans := make([]string, len(rgt))
		for i, rgtItem := range rgt {
			ans[i] = fmt.Sprintf("<link type=\"0-1\" xtargets=\";%s\" status=\"man\" />", rgtItem)
		}
		return ans
	}
	if item.To.First == -1 {
		ans := make([]string, len(lft))
		for i, lftItem := range lft {
			ans[i] = fmt.Sprintf("<link type=\"1-0\" xtargets=\"%s;\" status=\"man\" />", lftItem)
		}
		return ans
	}
	if len(lft) <= 2 && len(rgt) <= 2 {
		return []string{fmt.Sprintf("<link type=\"%d-%d\" xtargets=\"%s;%s\" status=\"man\" />",
			lftArity, rgtArity, strings.Join(lft, " "), strings.Join(rgt, " "))}
	}
	log.Print("WARNING: returning empty range - this should not happen ", item)
	return []string{}
}

// ungroupAndAdd  ungroups (if needed) items encoded in a numeric interval specified by "item".
// All the resulting groupIDs and *mapping.Mapping instances are then added to the pool
func (e *Export) ungroupAndAdd(item *mapping.Mapping) {
	var newGroup, currGroup string
	if item.From.First == -1 {
		currGroupStartIdx := item.To.First
		for i := item.To.First; i <= item.To.Last; i++ {
			newGroup = e.groupFilter.ExtractGroupID(e.Attr2.ID2Str(i))
			if newGroup != "" {
				if currGroup != "" && newGroup != currGroup {
					e.pool.AddGroup(currGroup, &mapping.Mapping{
						From: mapping.PosRange{First: -1, Last: -1},
						To:   mapping.PosRange{First: currGroupStartIdx, Last: i - 1},
					})
					currGroupStartIdx = i
				}
				currGroup = newGroup
			}
		}
		if newGroup != "" {
			e.pool.AddGroup(newGroup, &mapping.Mapping{
				From: mapping.PosRange{First: -1, Last: -1},
				To:   mapping.PosRange{First: currGroupStartIdx, Last: item.To.Last},
			})
		}

	} else if item.To.First == -1 {
		currGroupStartIdx := item.From.First
		for i := item.From.First; i <= item.From.Last; i++ {
			newGroup = e.groupFilter.ExtractGroupID(e.Attr1.ID2Str(i))
			if newGroup != "" {
				if currGroup != "" && newGroup != currGroup {
					e.pool.AddGroup(currGroup, &mapping.Mapping{
						From: mapping.PosRange{First: currGroupStartIdx, Last: i - 1},
						To:   mapping.PosRange{First: -1, Last: -1},
					})
					currGroupStartIdx = i
				}
				currGroup = newGroup
			}
		}
		if newGroup != "" {
			e.pool.AddGroup(newGroup, &mapping.Mapping{
				From: mapping.PosRange{First: currGroupStartIdx, Last: item.From.Last},
				To:   mapping.PosRange{First: -1, Last: -1},
			})
		}

	} else {
		currGroup = e.groupFilter.ExtractGroupID(e.Attr1.ID2Str(item.From.First))
		e.pool.AddGroup(currGroup, item)
	}
}

// getGroupIdent extracts a respective string identifier either from attr1 (i.e. first language)
// or attr2 (i.e. the second language).
func (e *Export) getGroupIdent(item *mapping.Mapping) string {
	var group string

	if item.From.First != -1 {
		group = e.groupFilter.ExtractGroupID(e.Attr1.ID2Str(item.From.First))
	}
	if group == "" && item.To.First != -1 {
		group = e.groupFilter.ExtractGroupID(e.Attr2.ID2Str(item.To.First))
	}
	return group
}

func (e *Export) printGroup(lang1, lang2 string, grp *gpool.TextGroup, ignoreEmpty bool, exportType string) {
	var bld strings.Builder
	grp.ForEach(func(mp *mapping.Mapping) {
		if !ignoreEmpty || (mp.From.First > -1 && mp.To.First > 1) {
			for _, line := range e.createTag(mp, exportType) {
				bld.WriteString(line + "\n")
			}
		}
	})
	if bld.Len() > 0 {
		fmt.Println(createGroupTag(lang1, lang2, grp.ID))
		fmt.Print(bld.String())
		fmt.Println("</linkGrp>")
	}
}

// Run generates a XML-ish output with the same format as the one
// used as input format for generating numerical alignment files.
// The algorithm is able to ungroup 'compressed' numeric intervals
// so if an interval contains multiple texts - all of them should
// be written to the output.
func (e *Export) Run(regPath1, regPath2, exportType string, skipEmpty bool) {
	srcFile, err := os.Open(e.MappingPath)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}
	e.groupFilter = NewGroupFilter(exportType)
	lang1 := e.groupFilter.ExtractLangFromRegistry(regPath1)
	lang2 := e.groupFilter.ExtractLangFromRegistry(regPath2)

	fmt.Println("<?xml version=\"1.0\" encoding=\"utf-8\"?>")
	fr := bufio.NewScanner(srcFile)
	var newGroup1 string
	e.pool = gpool.NewTextGroupPool()
	for i := 0; fr.Scan(); i++ {
		item, err := mapping.NewMappingFromString(fr.Text())
		if err != nil {
			log.Print("ERROR: ", err)
		}
		newGroup1 = e.getGroupIdent(&item)
		if newGroup1 != "" {
			e.ungroupAndAdd(&item)
			for nxt := e.pool.PopNextReady(); nxt != nil; nxt = e.pool.PopNextReady() {
				e.printGroup(lang1, lang2, nxt, skipEmpty, exportType)
			}
		}
	}
	for nxt := e.pool.PopOldest(); nxt != nil; nxt = e.pool.PopOldest() {
		e.printGroup(lang1, lang2, nxt, skipEmpty, exportType)
	}
}
