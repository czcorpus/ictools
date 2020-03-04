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

package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/czcorpus/ictools/attrib"
	"github.com/czcorpus/ictools/calign"
	"github.com/czcorpus/ictools/export"
	"github.com/czcorpus/ictools/fixgaps"
	"github.com/czcorpus/ictools/mapping"
	"github.com/czcorpus/ictools/transalign"
)

const (
	defaultChanBufferSize = 5000
)

type calignArgs struct {
	registryPath1   string
	registryPath2   string
	attrName        string
	mappingFilePath string
	bufferSize      int
	quoteStyle      int
}

type corpusPair struct {
	corp1 attrib.GoCorpus
	attr1 attrib.GoPosAttr
	corp2 attrib.GoCorpus
	attr2 attrib.GoPosAttr
}

func openCorpusPair(args calignArgs) *corpusPair {
	var err error

	c1, err := attrib.OpenCorpus(args.registryPath1)
	if err != nil {
		log.Fatalf("FATAL: Failed to open corpus %s: %s", args.registryPath1, err)
	}
	attr1, err := attrib.OpenAttr(c1, args.attrName)
	if err != nil {
		log.Fatalf("FATAL: Failed to open attribute %s: %s", args.attrName, err)
	}
	c2, err := attrib.OpenCorpus(args.registryPath2)
	if err != nil {
		log.Fatalf("FATAL: Failed to open corpus %s: %s", args.registryPath1, err)
	}
	attr2, err := attrib.OpenAttr(c2, args.attrName)
	if err != nil {
		log.Fatalf("FATAL: Failed to open attribute %s: %s", args.attrName, err)
	}
	return &corpusPair{
		corp1: c1,
		attr1: attr1,
		corp2: c2,
		attr2: attr2,
	}
}

func openAttribute(registryPath, attrName string) attrib.GoPosAttr {
	var err error

	corp, err := attrib.OpenCorpus(registryPath)
	if err != nil {
		log.Fatalf("FATAL: Failed to open corpus %s: %s", registryPath, err)
	}
	attr, err := attrib.OpenAttr(corp, attrName)
	if err != nil {
		log.Fatalf("FATAL: Failed to open attribute %s: %s", attrName, err)
	}
	return attr
}

func getStructSize(corp attrib.GoCorpus, structAttr string) (int, error) {
	structName := strings.Split(structAttr, ".")[0]
	return attrib.GetStructSize(corp, structName)
}

func prepareCalign(corps *corpusPair, mappingFilePath string, quoteStyle int) (*os.File, *calign.Processor) {
	var file *os.File
	var err error

	if mappingFilePath == "" {
		file = os.Stdin

	} else {
		file, err = os.Open(mappingFilePath)
		if err != nil {
			log.Fatalf("FATAL: Failed to open file %s", mappingFilePath)
		}
	}
	return file, calign.NewProcessor(corps.attr1, corps.attr2, quoteStyle)
}

func runTransalign(filePath1 string, filePath2 string) {
	var file1, file2 *os.File
	var err error

	file1, err = os.Open(filePath1)
	if err != nil {
		log.Fatalf("FATAL: Failed to open file %s", filePath1)
	}
	file2, err = os.Open(filePath2)
	if err != nil {
		log.Fatalf("FATAL: Failed to open file %s", filePath2)
	}
	if file2 != file2 {

	}
	hm1, err := transalign.NewPivotMapping(file1)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}
	err = hm1.Load()
	if err != nil {
		log.Fatal("FATAL: Failed to load pivot mapping 1: ", err)
	}
	hm2, err := transalign.NewPivotMapping(file2)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}
	err = hm2.Load()
	if err != nil {
		log.Fatal("FATAL: Failed to load pivot mapping 2: ", err)
	}

	ch1 := make(chan []mapping.Mapping, 5)
	buff1 := make([]mapping.Mapping, 0, defaultChanBufferSize)
	go func() {
		transalign.Run(hm1, hm2, func(item mapping.Mapping) {
			if !item.IsEmpty() {
				buff1 = append(buff1, item)
				if len(buff1) == defaultChanBufferSize {
					ch1 <- buff1
					buff1 = make([]mapping.Mapping, 0, defaultChanBufferSize)
				}
			}
		})
		if len(buff1) > 0 {
			ch1 <- buff1
		}
		close(ch1)
	}()
	calign.CompressFromChan(ch1, false, func(item mapping.Mapping) {
		item.IsGap = false
		fmt.Println(item)
	})
	log.Print("INFO: ...Done")
}

// runImport runs [calign] > [fixgaps] > [compress]? functions.
func runImport(args calignArgs) {
	corps := openCorpusPair(args)
	file, processor := prepareCalign(corps, args.mappingFilePath, args.quoteStyle)
	ch1 := make(chan []mapping.Mapping, 5)
	buff1 := make([]mapping.Mapping, 0, defaultChanBufferSize)
	go func() {
		err := processor.ProcessFile(file, args.bufferSize, func(item mapping.Mapping, i int) {
			buff1 = append(buff1, item)
			if len(buff1) == defaultChanBufferSize {
				ch1 <- buff1
				buff1 = make([]mapping.Mapping, 0, defaultChanBufferSize)
			}
		})
		if err != nil {
			log.Fatal("FATAL: ", err)

		} else if len(buff1) > 0 {
			ch1 <- buff1
		}
		close(ch1)
	}()

	ch2 := make(chan []mapping.Mapping, 5)
	go func() {
		var err error
		buff2 := make([]mapping.Mapping, 0, defaultChanBufferSize)
		s1Size, err := getStructSize(corps.corp1, args.attrName)
		if err != nil {
			log.Fatalf("FATAL: Cannot determine size of structure %s (%s)", args.attrName, args.registryPath1)
		}
		s2Size, err := getStructSize(corps.corp2, args.attrName)
		if err != nil {
			log.Fatalf("FATAL: Cannot determine size of structure %s (%s)", args.attrName, args.registryPath2)
		}

		errors := make([]error, 0, 10)
		fixgaps.FromChan(ch1, true, s1Size, s2Size, func(item mapping.Mapping, err *fixgaps.FixGapsError) {
			if err != nil {
				log.Print("ERROR: ", err)
				log.Printf("INFO: original struct idents are: item: [%s, %s -- %s, %s], reached positions: [%s, %s]",
					corps.attr1.ID2Str(err.Item.From.First), corps.attr1.ID2Str(err.Item.From.Last),
					corps.attr2.ID2Str(err.Item.To.First), corps.attr2.ID2Str(err.Item.To.Last),
					corps.attr1.ID2Str(err.Left), corps.attr2.ID2Str(err.Pivot))
				buff2 = append(buff2, mapping.NewErrorMapping())
				errors = append(errors, err)

			} else {
				buff2 = append(buff2, item)
			}
			if len(buff2) == defaultChanBufferSize {
				ch2 <- buff2
				buff2 = make([]mapping.Mapping, 0, defaultChanBufferSize)
			}
		})
		if len(buff2) > 0 {
			ch2 <- buff2
		}
		close(ch2)
		if len(errors) > 0 {
			log.Fatalf("FATAL: Finished with %d errors. The result cannot be used to produce a correct alignment.", len(errors))
		}
	}()
	calign.CompressFromChan(ch2, true, func(item mapping.Mapping) {
		fmt.Println(item)
	})

}

func runSearch(corpusRegistry string, attr string, itemIdx int) {
	attrObj := openAttribute(corpusRegistry, attr)
	fmt.Printf("\n\nPosition %d -> %s\n\n", itemIdx, attrObj.ID2Str(itemIdx))
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s [options] import [LANG registry] [PIVOT registry] [attr] [LANG-PIVOT mapping file]?\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s [options] transalign [LANG1-PIVOT alignment file] [LANG2-PIVOT alignment file]\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s [options] search [LANG registry] [attr] [srch position]\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s [options] export [LANG1 registry] [LANG2 registry] [attr] [LANG1-LANG2 numeric mapping file]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	var lineBufferSize int
	flag.IntVar(&lineBufferSize, "line-buffer", bufio.MaxScanTokenSize, "Max line buffer size")
	var registryPath string
	flag.StringVar(&registryPath, "registry-path", "", "Path to Manatee registry files (allows using just filenames for registry values in 'import')")
	var quoteStyle int
	flag.IntVar(&quoteStyle, "quote-style", 1, "Input XML quote style: 1 - single, 2 - double")
	var exportType string
	flag.StringVar(&exportType, "export-type", "",
		fmt.Sprintf("Select specific tools to export data. Currently supported types: %s", export.GroupFilterTypeIntercorp))

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Missing action, try -h for help")
		os.Exit(1)

	} else {
		t1 := time.Now().UnixNano()
		switch flag.Arg(0) {
		case "transalign":
			runTransalign(flag.Arg(1), flag.Arg(2))
		case "import":
			runImport(calignArgs{
				registryPath1:   filepath.Join(registryPath, flag.Arg(1)),
				registryPath2:   filepath.Join(registryPath, flag.Arg(2)),
				attrName:        flag.Arg(3),
				mappingFilePath: flag.Arg(4),
				bufferSize:      lineBufferSize,
				quoteStyle:      quoteStyle,
			})
		case "search":
			itemIdx, err := strconv.Atoi(flag.Arg(3))
			if err != nil {
				log.Fatalf("FATAL: failed to parse item position: %s. Expected integer number.", flag.Arg(1))
			}
			runSearch(flag.Arg(1), flag.Arg(2), itemIdx)
		case "export":
			regPath1 := filepath.Join(registryPath, flag.Arg(1))
			regPath2 := filepath.Join(registryPath, flag.Arg(2))
			corps := openCorpusPair(calignArgs{
				registryPath1: regPath1,
				registryPath2: regPath2,
				attrName:      flag.Arg(3),
			})
			export := export.Export{
				RegPath1:    regPath1,
				Corp1:       corps.corp1,
				Attr1:       corps.attr1,
				RegPath2:    regPath2,
				Corp2:       corps.corp2,
				Attr2:       corps.attr2,
				MappingPath: flag.Arg(4),
			}
			export.Run(regPath1, regPath2, exportType)
		default:
			log.Fatalf("FATAL: Unknown action '%s'", flag.Arg(0))
		}
		log.Printf("INFO: Finished in %01.2f sec.", float64(time.Now().UnixNano()-t1)/1e9)
	}
}
