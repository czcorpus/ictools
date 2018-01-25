# ictools - a set of tools for generating corpora alignments

This is a faster integrated replacement for "classic" *calign.py*, *compressrng.py*, *fixgaps.py*,
*transalign.py* scripts used to prepare corpora alignment numeric data from lists of structural
attribute values mapping between languages.

## Contents

* [How to build ictools](#how_to_build_ictools)
  * [build helper script](#how_to_build_ictools_helper_script)
  * [manual variant](#how_to_build_ictools_manual_variant)
* [Using ictools](#using_ictools)
* [Benchmark](#benchmark)

<a name="how_to_build_ictools"></a>
## How to build ictools

To build ictools, a working installation of Manatee(-open) must be installed on your system.
This includes not just *libmanatee.so* shared library but also source header files for both
Manatee(-open) and Finlib. A working Go language environment is also required.

Download ictools package:

```bash
go get -d https://github.com/czcorpus/ictools
```

<a name="how_to_build_ictools_helper_script"></a>
### build helper script

The *build* script handles all the intricacies regarding miscellaneous command line arguments
and environment variables required to build the project and link it with *libmanatee.so* properly.
It is written in Python 2 and requires Manatee Python wrapper libraries (these are created by default
when building Manatee).

Let's say you have Manatee-open 2.150 installed on your system. Then just enter:

```bash
./build 2.150
```

The script looks for *libmanatee.so* in typical locations (/usr/lib and /usr/local/lib) and
downloads sources for *Manatee-open* and matching *Finlib* version.

In case your Manatee installation resides in a custom directory, you must specify it yourself:

```
./build 2.150 --manatee-lib /opt/manatee-2.150/lib
```

In both cases, the *build* script assumes that the version of the specified (or automatically found)
Manatee library and the version entered as the first argument (here 2.150) match together.

Script finishes in one of two possible result states:

* Two created files: *ictools* (a bash script), *ictools.bin* (a binary executable) in case *LD_LIBRARY_PATH* must be
  set because of a non-standard location of *libmanatee.so*. You can just copy these files to */usr/local/bin*
  and refer the program as *ictools*.
* One file: *ictools* (a binary executable) in case of standard system installation of *libmanatee.so* (i.e. the
  operating system is able to locate *libmanatee.so* by itself).

<a name="how_to_build_ictools_manual_variant"></a>
### manual variant

In many cases, simple `go build` won't work because of missing header files and/or non-standard
*libmanatee.so* location. In such case you have to specify all the locations by yourself when
building the project:

```bash
CGO_CPPFLAGS="-I/path/to/manatee/src -I/path/to/finlib/src" CGO_LDFLAGS="-lmanatee -L/path/to/manatee/lib/dir" go build
```

In case you have installed *manatee-open* to a non-standard directory, then you have to tell the OS where to look
for *libmanatee*:

```
LD_LIBRARY_PATH="/path/to/libmanatee.so/dir" ./ictools
```

or you can write a simple start script:

```bash
#!/usr/bin/env bash
export LD_LIBRARY_PATH="/path/to/libmanatee.so/dir"
`dirname $0`/ictools "${@:1}"
```

In case your have `$GOPATH/bin` in your `$PATH` you are ready to go. Otherwise you can copy the
compiled binary to a location like `/usr/local/bin` to be able to call it without referring its full
path.

<a name="using_ictools"></a>
## Using ictools


### The default usage style

This usage style is the recommended one as it handles whole import of XML data
(= *calign* -> *fixgaps* -> *compress*) in one step. The individual transformations the import
is composed of run concurrently to be able to keep up with the classic scripts connected via
pipes (where all the processes run concurrently too). Ictools' approach is a little bit more
efficient as there is no process overhead, no repeated data serialization/deserialization.

To prepare alignment data, two actions are necessary:

1. importing of two or more XML files containing mappings between structures (typically sentences) of
   two languages (one of them is considered a *pivot*) identified by their string IDs.
1. create a new mapping between two or more non-pivot languages

Please note that the parser does not care about XML validity - it just looks for tags with the
following form (actually, only *xtargets* attribute is significant):

```xml
<link type='0-3' xtargets=';cs:Adams-Holisticka_det_k:0:7:1 cs:Adams-Holisticka_det_k:0:7:2 cs:Adams-Holisticka_det_k:0:7:3' status='man'/>
```

Let's say we have two files with mappings between Polish and Czech (*intercorp_pl2cs*) and between
English and Czech (*intercorp.en2cs*) where Czech is a pivot.

```bash
ictools -registry-path /var/local/corpora/registry import intercorp_v10_pl intercorp_v10_cs s.id /var/local/corpora/aligndef/intercorp_pl2cs > intercorp.pl2cs

ictools -registry-path /var/local/corpora/registry import intercorp_v10_en intercorp_v10_cs s.id /var/local/corpora/aligndef/intercorp_en2cs > intercorp.en2cs

ictools transalign ./intercorp.pl2cs ./intercorp.en2cs > intercorp.pl2en
```

For the *import* action, you may want to *tweak line buffer size* (by default *bufio.MaxScanTokenSize* = 64 * 1024
is used which may fail in case of some complex alignments):

```bash
ictools -line-buffer 250000 -registry-path /var/local/corpora/registry import ....etc...
```

In case you do not want a result file to be compressed, use *no-compress* arg:

```bash
ictools -no-compress ....etc....
```

### The "old" usage style

This is for legacy (and debugging) reasons and it should work in a similar way to the Python scripts
*calign.py*, *compressrng.py*, *fixgaps.py* and *transalign.py*.

```bash
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

* A (a server)
  * CPU: Intel Xeon E5-2640 v3 @ 2.60GHz
  * 64GB RAM
* B (a common Dell desktop)
  * CPU: Intel Core) i5-2400 @ 3.10GHz
  * 8GB RAM

| Setup | Used program    | calign+fixgaps+compress [sec] | transalign [sec] | total [sec]  |
|-------|-----------------|------------------------------:|-----------------:|-------------:|
| A     | classic scripts |  255                          | 191              | 446          |
| A     | ictools         |  **180**                      | **57**           | **237**      |
| B     | classic scripts |  312                          | DNF (RAM)        | DNF          |
| B     | ictools         |  **175**                      | **63**           | **238**      |

In terms of memory usage, there were no thorough measurements performed but according to the *top*
utility the *transalign* function in *ictools* consumes less than half of the memory compared
with the classic scripts. The import function (i.e. calign+fixgaps+compress) in both programs
consumes only a little memory because data read from an input file are (almost) immediately written
to the output without any unnecessary memory allocation.
