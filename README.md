# ictools - a set of tools for generating corpora alignments

This is a replacement for "classic" *calign.py*, *compressrng.py*, *fixgaps.py*, *transalign.py*.

:construction: Current status: a working prototype, misc. tested must be written and performed

## How to build ictools

First, check whether [manatee-open](https://nlp.fi.muni.cz/trac/noske/wiki/Downloads) is installed on your system (no Python, Ruby, Java etc. APIs are needed - just a native library). E.g. by `ldconfig -p | grep manatee`. A working installation of [Go](https://www.golang.org) must be also available. 

```
go get https://github.com/czcorpus/ictools
```

:construction_worker: Please note that currently, *ictools* come with required *manatee-open* header files which is convenient but it can be a problem in case *manatee-open* on your system differs from the one the headers were copied from.
This issue will be solved once a first release is ready.

In case your have `$GOPATH/bin` in your `$PATH` you are ready to go. Otherwise you can copy the compiled binary to a location like `/usr/local/bin` to be able to call it without referring its full path.

## Using ictools

### The "new way"

This approach handles whole import of XML data (= *calign* > *fixgaps* > *compress*) in one step. 
The individual transformations the import is composed of run concurrently to be able to keep up with 
the classic scripts connected via pipes (where all the processes run concurrently too). Ictools' 
approach is a little bit more efficient as there is no process overhead, no repeated data serialization 
and deserialization.

To prepare your data, two actions are necessary:

1. import two or more XML files containing mappings between structures (typically sentences) of two languages (one of them is considered a *pivot*) identified by their string IDs.
2. create a new mapping between two or more non-pivot languages


Let's say we have two files with mappings between Polish and Czech (*intercorp_pl2cs*) and between English and Czech (*intercorp.en2cs*) where Czech is a pivot.

```
ictools -registry-path /var/local/corpora/registry import intercorp_v10_pl intercorp_v10_cs s.id /var/local/corpora/aligndef/intercorp_pl2cs > intercorp.pl2cs

ictools -registry-path /var/local/corpora/registry import intercorp_v10_en intercorp_v10_cs s.id /var/local/corpora/aligndef/intercorp_en2cs > intercorp.en2cs

ictools transalign ./intercorp.pl2cs ./intercorp.en2cs > intercorp.pl2en
```

For the *import* action, you may want to *tweak line buffer size* (by default *bufio.MaxScanTokenSize* = 64 * 1024
is used which may fail in case of some complex alignments):

```
ictools -line-buffer 250000 -registry-path /var/local/corpora/registry import ....etc...
```

In case you do not want a result file to be compressed, use *no-compress* arg:

```
ictools -no-compress ....etc....
```

### The "old way"

This is for legacy (and debugging) reasons and it should work in a similar way to the Python scripts  *calign.py*, *compressrng.py*, *fixgaps.py* and *transalign.py*.

```
ictools calign -registry-path /var/local/corpora/registry import intercorp_v10_pl intercorp_v10_cs s.id /var/local/corpora/aligndef/intercorp_pl2cs | ictools fixgaps | ictools compressrng > intercorp.pl2cs

ictools calign -registry-path /var/local/corpora/registry import intercorp_v10_en intercorp_v10_cs s.id /var/local/corpora/aligndef/intercorp_en2cs | ictools fixgaps | ictools compressrng > intercorp.en2cs

ictools transalign ./intercorp.pl2cs ./intercorp.en2cs > intercorp.pl2en
```
