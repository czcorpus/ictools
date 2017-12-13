//  Copyright (c) 1999-2012  Pavel Rychly, Milos Jakubicek

#ifndef EXCEP_HH
#define EXCEP_HH

#include <string>
#include <sstream>

class CorpInfoNotFound: public std::exception {
    const std::string _what;
public:
    const std::string name;
    CorpInfoNotFound (const std::string &name)
        :_what ("CorpInfoNotFound (" + name + ")"), name (name) {}
    virtual const char* what () const throw () {return _what.c_str();}
    virtual ~CorpInfoNotFound() throw() {}
};

class AttrNotFound: public std::exception {
    const std::string _what;
public:
    const std::string name;
    AttrNotFound (const std::string &name)
        :_what ("AttrNotFound (" + name + ")"), name (name) {}
    virtual const char* what () const throw () {return _what.c_str();}
    virtual ~AttrNotFound() throw() {}
};

class NotImplemented: public std::exception {
    std::string _what;
public:
    NotImplemented (const std::string func, const std::string file, int line) {
        std::stringstream ss;
        ss << func << " not implemented (" << file << ": " << line << ")";
        _what = ss.str();
    }
    virtual const char* what () const throw () {return _what.c_str();}
    virtual ~NotImplemented() throw () {}
};
#define NOTIMPLEMENTED throw NotImplemented(__func__,__FILE__,__LINE__);

#endif

// vim: ts=4 sw=4 sta et sts=4 si cindent tw=80:
