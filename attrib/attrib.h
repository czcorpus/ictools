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

#ifdef __cplusplus
extern "C" {
#endif

typedef void* PosAttrV;
typedef void* CorpusV;

PosAttrV get_attr(CorpusV corp, const char* attrName);
long attr_str2id(PosAttrV attr, const char* str);
CorpusV open_corpus(const char* corpusPath);
void close_corpus(CorpusV corpus);
void close_attr(PosAttrV attr);

#ifdef __cplusplus
}
#endif
