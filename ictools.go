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
	"strings"
	"time"

	"github.com/czcorpus/ictools/attrib"
	"github.com/czcorpus/ictools/calign"
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

func openCorpora(args calignArgs) *corpusPair {
	var err error

	c1, err := attrib.OpenCorpus(args.registryPath1)
	if err != nil {
		log.Fatalf("FATAL: Failed to open corpus %s", args.registryPath1)
	}
	attr1, err := attrib.OpenAttr(c1, args.attrName)
	if err != nil {
		log.Fatalf("FATAL: Failed to open attribute %s", args.attrName)
	}
	c2, err := attrib.OpenCorpus(args.registryPath2)
	if err != nil {
		log.Fatalf("FATAL: Failed to open corpus %s", args.registryPath1)
	}
	attr2, err := attrib.OpenAttr(c2, args.attrName)
	if err != nil {
		log.Fatalf("FATAL: Failed to open attribute %s", args.attrName)
	}
	return &corpusPair{
		corp1: c1,
		attr1: attr1,
		corp2: c2,
		attr2: attr2,
	}
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
	hm1.Load()
	hm2, err := transalign.NewPivotMapping(file2)
	if err != nil {
		log.Fatal("FATAL: ", err)
	}
	hm2.Load()

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
	corps := openCorpora(args)
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

		err = fixgaps.FromChan(ch1, true, s1Size, s2Size, func(item mapping.Mapping) {
			buff2 = append(buff2, item)
			if len(buff2) == defaultChanBufferSize {
				ch2 <- buff2
				buff2 = make([]mapping.Mapping, 0, defaultChanBufferSize)
			}
		})
		if err != nil {
			log.Fatal("FATAL: ", err)

		} else if len(buff2) > 0 {
			ch2 <- buff2

		}
		close(ch2)
	}()
	calign.CompressFromChan(ch2, true, func(item mapping.Mapping) {
		fmt.Println(item)
	})

}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s [options] import [LANG registry] [PIVOT registry] [attr] [LANG-PIVOT mapping file]?\n", filepath.Base(os.Args[0]))
		fmt.Fprintf(os.Stderr, "\t%s [options] transalign [LANG1-PIVOT alignment file] [LANG2-PIVOT alignment file]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	var lineBufferSize int
	flag.IntVar(&lineBufferSize, "line-buffer", bufio.MaxScanTokenSize, "Max line buffer size")
	var registryPath string
	flag.StringVar(&registryPath, "registry-path", "", "Path to Manatee registry files (allows using just filenames for registry values in 'import')")
	var quoteStyle int
	flag.IntVar(&quoteStyle, "quote-style", 1, "Input XML quote style: 1 - single, 2 - double")

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
		default:
			log.Fatalf("FATAL: Unknown action '%s'", flag.Arg(0))
		}
		log.Printf("INFO: Finished in %01.2f sec.", float64(time.Now().UnixNano()-t1)/1e9)
	}
}
