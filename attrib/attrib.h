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

#ifdef __cplusplus
extern "C" {
#endif

typedef void* PosAttrV;
typedef void* CorpusV;
typedef void* StructV;

/**
 * CorpusRetval wraps both
 * a returned Manatee corpus object
 * and possible error
 */
typedef struct CorpusRetval {
    CorpusV value;
    const char * err;
} CorpusRetval;


/**
 * AttrRetval wraps both
 * a returned Manatee attribute object
 * and possible error
 */
typedef struct AttrRetval {
    PosAttrV value;
    const char * err;
} AttrRetval;

/**
 * StructSizeRetval wraps both
 * a returned size of a structure object
 * and possible error
 */
typedef struct StructSizeRetval {
    long value;
    const char * err;
} StructSizeRetval;

/**
 * Provide number of structures of a given name
 */
StructSizeRetval get_struct_size(CorpusV corpus, const char* structName);

/**
 * Return a Manatee PosAttr instance
 */
AttrRetval get_attr(CorpusV corp, const char* attrName);

/**
 * Get numeric identifier of a provided PosAttr's original string value.
 */
long attr_str2id(PosAttrV attr, const char* str);

/**
 * Get original string identifier of a provided PosAttr's numeric value.
 */
const char* attr_id2str(PosAttrV attr, long ident);

/**
 * Create a Manatee corpus instance
 */
CorpusRetval open_corpus(const char* corpusPath);

/**
 * Note: currently not used but we probably should
 * to make the code correct (yet there should be
 * no memory issues as we always open two corpora
 * and then the program exits).
 */
void close_corpus(CorpusV corpus);

/**
 * Note: currently not used; for more info see
 * close_corpus.
 */
void close_attr(PosAttrV attr);

#ifdef __cplusplus
}
#endif
