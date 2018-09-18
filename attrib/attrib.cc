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
#include <string.h>
#include <stdio.h>
#include <iostream>

using namespace std;

// a bunch of wrapper functions we need to get data
// from Manatee


AttrRetval get_attr(CorpusV corp, const char* attrName) {
    string tmp(attrName);
    AttrRetval ans;
    try {
        PosAttrV attr = ((Corpus*)corp)->get_attr(tmp);
        ans.value = attr;
        return ans;

    } catch (std::exception &e) {
        ans.err = strdup(e.what());
        return ans;
    }
}

long attr_str2id(PosAttrV attr, const char* str) {
    return ((PosAttr *)attr)->str2id(str);
}

StructSizeRetval get_struct_size(CorpusV corpus, const char* structName) {
    string tmp(structName);
    StructSizeRetval ans;
    try {
        StructV strct = ((Corpus*)corpus)->get_struct(tmp);
        ans.value = ((Structure *)strct)->size();
        return ans;

    } catch (std::exception &e) {
        ans.err = strdup(e.what());
        return ans;
    }
}

CorpusRetval open_corpus(const char* corpusPath) {
    string tmp(corpusPath);
    CorpusRetval ans;
    try {
        ans.value = new Corpus(tmp);
        return ans;

    } catch (std::exception &e) {
        ans.err = strdup(e.what());
        return ans;
    }
}

void close_corpus(CorpusV corpus) {
    delete (Corpus *)corpus;
}

void close_attr(PosAttrV attr) {
    delete (PosAttr *)attr;
}
