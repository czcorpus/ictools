# ictools - a set of tools for generating corpora alignments

This is a faster replacement for "classic" *calign.py*, *compressrng.py*, *fixgaps.py*, *transalign.py*.

:construction: Current status: a working prototype, misc. tested must be written and performed

* [How to build ictools](#how_to_build_ictools)
* [Using ictools](#using_ictools)
* [Benchmark](#benchmark)

<a name="how_to_build_ictools"></a>
## How to build ictools

First, check whether [manatee-open](https://nlp.fi.muni.cz/trac/noske/wiki/Downloads) is installed
on your system (no Python, Ruby, Java etc. APIs are needed - just a native library).
E.g. by `ldconfig -p | grep manatee`. A working installation of [Go](https://www.golang.org) must be
also available.

```
go get https://github.com/czcorpus/ictools
```

### Possible issues

In case you have installed *manatee-open* into a directory where OS does not look when searching for libraries
(typically */usr/local/lib* for installations from a source code) then the *go build*
command needs some more arguments to tell compiler and linker where to look for C header files and
compiled *manatee* library:

```
CGO_CPPFLAGS="-I/opt/manatee/2.130.6/include" CGO_LDFLAGS="-lmanatee -L/opt/manatee/2.130.6/lib" go build
```

:construction_worker: Please note that currently, *ictools* come with required *manatee-open* header
files which is convenient but it can be a problem in case *manatee-open* on your system differs from
the one the headers were copied from.
This issue will be solved once a first release is ready.

In case your have `$GOPATH/bin` in your `$PATH` you are ready to go. Otherwise you can copy the
compiled binary to a location like `/usr/local/bin` to be able to call it without referring its full
path.

<a name="using_ictools"></a>
## Using ictools

Note: in case you have installed *manatee-open* to a directory OS does not about when looking for
libraries, then you have to tell where to look for *libmanatee*:

```
LD_LIBRARY_PATH="/opt/manatee/2.130.6/lib" ./ictools
```

or you can write a simple start script:

```bash
#!/bin/bash
export LD_LIBRARY_PATH="/opt/manatee/2.130.6/lib"
`dirname $0`/ictools "${@:1}"
```

### The default usage style

This usage style is the recommended one as it handles whole import of XML data
(= *calign* -> *fixgaps* -> *compress*) in one step. The individual transformations the import
is composed of run concurrently to be able to keep up with the classic scripts connected via
pipes (where all the processes run concurrently too). Ictools' approach is a little bit more
efficient as there is no process overhead, no repeated data serialization and deserialization.

To prepare alignment data, two actions are necessary:

1. importing of two or more XML files containing mappings between structures (typically sentences) of
   two languages (one of them is considered a *pivot*) identified by their string IDs.
2. create a new mapping between two or more non-pivot languages

Please note that the parser does not care about XML validity - it just looks for tags with the
following form (actually, only *xtargets* attribute is significant):

```xml
<link type='0-3' xtargets=';cs:Adams-Holisticka_det_k:0:7:1 cs:Adams-Holisticka_det_k:0:7:2 cs:Adams-Holisticka_det_k:0:7:3' status='man'/>
```


Let's say we have two files with mappings between Polish and Czech (*intercorp_pl2cs*) and between
English and Czech (*intercorp.en2cs*) where Czech is a pivot.

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

### The "old" usage style

This is for legacy (and debugging) reasons and it should work in a similar way to the Python scripts
*calign.py*, *compressrng.py*, *fixgaps.py* and *transalign.py*.

```
ictools calign -registry-path /var/local/corpora/registry import intercorp_v10_pl intercorp_v10_cs s.id /var/local/corpora/aligndef/intercorp_pl2cs | ictools fixgaps | ictools compressrng > intercorp.pl2cs

ictools calign -registry-path /var/local/corpora/registry import intercorp_v10_en intercorp_v10_cs s.id /var/local/corpora/aligndef/intercorp_en2cs | ictools fixgaps | ictools compressrng > intercorp.en2cs

ictools transalign ./intercorp.pl2cs ./intercorp.en2cs > intercorp.pl2en
```

<a name="benchmark"></a>
## Benchmark

Used data files:
  * intercorp_pl2cs (size 1.4GB)
  * intercorp_pl2en (size 1.5GB)

Used hardware:
  * CPU: Intel Xeon E5-2640 v3 @ 2.60GHz
  * 64GB RAM


| Used program  | calign+fixgaps+compress [sec] | transalign [sec] | total [sec]  |
----------------|------------------------------:|-----------------:|-------------:|
classic scripts |  255                          | 191              | 446          |
ictools         |  180                          | 57               | 237          |

In terms of memory usage, there were no thorough measurements performed but according to the *top*
utility the *transalign* function in *ictools* consumes less than half of the memory compared
with the classic scripts. The import function (i.e. calign+fixgaps+compress) in both programs
consumes only a little memory because data read from an input file are (almost) immediately written
to the output without any unnecessary memory allocation.
