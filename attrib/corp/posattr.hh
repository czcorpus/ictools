// Copyright (c) 1999-2013  Pavel Rychly, Milos Jakubicek

#ifndef POSATTR_HH
#define POSATTR_HH

#include "excep.hh"
#include <finlib/fsop.hh>
#include <exception>

class CorpInfo;

class IDIterator
{
public:
    virtual int next() =0;
    virtual ~IDIterator() {}
};

class TextIterator
{
public:
    virtual const char *next() =0;
    virtual ~TextIterator() {}
};

class IDPosIterator {
    IDIterator *ids;
    FastStream *poss;
    int curr_id;
public:
    IDPosIterator (IDIterator *it, FastStream *fs):
        ids (it), poss (fs), curr_id (it->next()) {}
    IDPosIterator(): ids (NULL), poss (NULL) {}
    virtual ~IDPosIterator() {delete poss; delete ids;}
    virtual void next() {poss->next(); curr_id = ids->next();}
    virtual Position peek_pos() {return poss->peek();}
    virtual NumOfPos get_delta() {return 0;}
    virtual int peek_id() {return curr_id;}
    virtual bool end() {return poss->peek() >= poss->final();}
};

class DummyIDIter: public IDIterator {
public:
    DummyIDIter () {}
    virtual int next() {return 0;}
};

class DummyTextIter: public TextIterator {
const string s;
public:
    DummyTextIter (const string &s): s(s) {}
    virtual const char *next() {return s.c_str();}
};

class DummyIDPosIter: public IDPosIterator {
    FastStream *poss;
public:
    DummyIDPosIter (FastStream *fs): poss(fs) {}
    virtual ~DummyIDPosIter() {delete poss;}
    virtual void next() {poss->next();}
    virtual Position peek_pos() {return poss->peek();}
    virtual NumOfPos get_delta() {return 0;}
    virtual int peek_id() {return 0;}
    virtual bool end() {return poss->peek() >= poss->final();}
};


class PosAttr 
{
public:
    const std::string attr_path;
    const std::string name;
    const char *locale;
    const char *encoding;
    PosAttr (const std::string &path, const std::string &n, 
             const std::string &loc="", const std::string &enc="");
    virtual ~PosAttr ();
    virtual int id_range () =0;
    virtual const char* id2str (int id) =0;
    virtual int str2id (const char *str) =0;
    virtual int pos2id (Position pos) =0;
    virtual const char* pos2str (Position pos) =0;
    virtual IDIterator *posat (Position pos) =0;
    virtual IDPosIterator *idposat (Position pos) =0;
    virtual TextIterator *textat (Position pos) =0;
    virtual FastStream *id2poss (int id) =0;
    virtual FastStream *dynid2srcids (int id) =0;
    virtual FastStream *regexp2poss (const char *pat, bool ignorecase) =0;
    virtual FastStream *compare2poss (const char *pat, int cmp, bool ignorecase) =0;
    virtual Generator<int> *regexp2ids (const char *pat, bool ignorecase, const char *filter_pat = NULL) =0;
    virtual NumOfPos freq (int id) =0;
    virtual NumOfPos docf (int id) =0;
    virtual float arf (int id) =0;
    virtual float aldf (int id) =0;
    virtual NumOfPos norm (int id) =0;
    virtual NumOfPos size() =0;
};


const char *locale2c_str (const std::string &locale);


PosAttr *createPosAttr (std::string &typecode, const std::string &path, 
                        const std::string &n, const std::string &locale, 
                        const std::string &encoding, NumOfPos text_size=0);
PosAttr *createSubCorpPosAttr (PosAttr *pa, const std::string& subpath,
                               bool complement);

#endif

// vim: ts=4 sw=4 sta et sts=4 si cindent tw=80:

