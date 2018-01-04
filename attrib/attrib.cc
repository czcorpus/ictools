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

#include "corp/corpus.hh"
#include "attrib.h"
#include <stdio.h>
#include <iostream>

using namespace std;


PosAttrV get_attr(CorpusV corp, const char* attrName) {
    string tmp(attrName);
    PosAttrV attr = ((Corpus*)corp)->get_attr(tmp);
    return attr;
}

long attr_str2id(PosAttrV attr, const char* str) {
    return ((PosAttr *)attr)->str2id(str);
}

CorpusV open_corpus(const char* corpusPath) {
    string tmp(corpusPath);
    return new Corpus(tmp);
}

void close_corpus(CorpusV corpus) {
    delete (Corpus *)corpus;
}

void close_attr(PosAttrV attr) {
    delete (PosAttr *)attr;
}

/*
int main() {

	CorpusV corp = open_corpus("/var/local/corpora/registry/syn2015");
	PosAttrV attr = get_attr(corp, "s.id");
	cout << "corp: " << ((PosAttr *)attr)->str2id("5_elefan:1:4641:4")  << endl;

	close_attr(attr);
	//close_corpus(corp);

}
*/
