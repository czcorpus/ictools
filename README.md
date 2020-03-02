# ictools - a program for calculating corpora alignments using a pivot language

This is a faster, less memory-consuming, integrated replacement for legacy *calign.py*,
*compressrng.py*, *fixgaps.py*, *transalign.py* scripts used to prepare corpora alignment
numeric data from lists of structural attribute values mapping between languages. It also fixes
some problems with missing ranges for unaligned structures you can encounter when using the scripts above.
In addition, it also provides an `export` function for performing reversed operations.

Note: you still need *mkalign* tool distributed along with *Manatee-open* to enable corpora alignments
in *KonText* (or NoSkE).

## Contents

* [Using ictools](#using_ictools)
* [How to build ictools](#how_to_build_ictools)
  * [build helper script](#how_to_build_ictools_helper_script)
  * [manual variant](#how_to_build_ictools_manual_variant)
* [Benchmark](#benchmark)
* [For developers](#for_developers)
  * [Setting up VSCode debugging/testing environment](#for_developers_setting_up_vscode)
  * [running tests](#for_developers_running_tests)

<a name="using_ictools"></a>
## Using ictools

*Ictools* provide three operations - import, transalign and export:

### import

Import operation transforms an alignment XML file containing aligned string sentence IDs to a numeric form.
It is able to handle non-existing alignments, gaps between ranges (including the last row range where structure
size is always used to make sure the whole range is filled in).

In terms of the input format, a list of `&lt;link&gt;` elements is expected:

```xml
<link type='0-3' xtargets=';cs:Adams-Holisticka_det_k:0:7:1 cs:Adams-Holisticka_det_k:0:7:2 cs:Adams-Holisticka_det_k:0:7:3' status='man'/>
```

Please note that the parser does not care about XML validity (e.g. there is no need for a root element or even
a proper nesting of elements).

In some cases you may want to *tweak line buffer size* (value is in bytes; by default *bufio.MaxScanTokenSize* = 64 * 1024 is used which may fail in case of some complex alignments and/or long text identifiers). In case the buffer is too
small, ictools will end with fatal log event returning a non-zero value to shell.

```
ictools -line-buffer 250000 -registry-path /var/local/corpora/registry import ....etc...
```

**Example:**

Let's say we have two files with mappings between Polish and Czech (*intercorp_pl2cs*) and between
English and Czech (*intercorp.en2cs*) where Czech is a pivot.

```
ictools -registry-path /var/local/corpora/registry import intercorp_v10_pl intercorp_v10_cs s.id /var/local/corpora/aligndef/intercorp_pl2cs > intercorp.pl2cs

ictools -registry-path /var/local/corpora/registry import intercorp_v10_en intercorp_v10_cs s.id /var/local/corpora/aligndef/intercorp_en2cs > intercorp.en2cs
```

### transalign

Transalign operation takes two numeric alignments against a common pivot language and generates
a new alignment between the two non-pivot languages.

**Example:**

```
ictools transalign ./intercorp.pl2cs ./intercorp.en2cs > intercorp.pl2en
```

### export

The `export` operation is able to reconstruct the XML-ish source used as an input
for the `import` operation using numeric alignment files as produced by
`import -> transalign` operations. Any grouped intervals are split back to the original
text groups.

**Example:**

```
ictools -export-type intercorp export /corpora/registry/intercorp_v12_cs /corpora/registry/intercorp_v12_en s.id /corpora/aligndef/intercorp.cs2en > orig.xml
```


<a name="how_to_build_ictools"></a>
## How to build ictools

To build ictools, a working installation of Manatee(-open) must be installed on your system.
This includes not just *libmanatee.so* shared library but also source header files for both
Manatee(-open) and Finlib. Please note that Manatee starting from 2.158.8 includes Finlib so
there is no need to build Finlib separately. A working Go language environment is also required
(see [install instructions](https://golang.org/doc/install)).

Download ictools package:

```bash
go get -d https://github.com/czcorpus/ictools
```

<a name="how_to_build_ictools_helper_script"></a>
### build helper script

*Ictools* are written to work directly with [Manatee-open](https://nlp.fi.muni.cz/trac/noske) library which
itself is written in C++. This makes the build process a little more complicated then just `go get ...` or `go build`.

*Ictools* come with a simple *build* script (written in Python 2) which is able to handle all the details
for you. In the best scenario, the script requires only Manatee-open version you are building against. It tries to
find the library in system lib paths and download sources from Manatee-open project page. In case you have your Manatee-open installed in a non-standard location, you have to tell the script with `--manatee-lib` parameter. Also, if you already have Manatee-open sources downloaded, you can skip the download by specifying source directory
via `--manatee-src` (or `--finlib-src`).

Let's say you have Manatee-open 2.150 installed on your system. Then just enter:

```bash
./build 2.150
```

The script looks for `libmanatee.so` in typical locations (`/usr/lib` and `/usr/local/lib`) and
downloads sources for `Manatee-open` and matching `Finlib` version.

In case your Manatee installation resides in a custom directory, you must specify it yourself:

```
./build 2.150 --manatee-lib /opt/manatee-2.150/lib
```

In both cases, the *build* script assumes that the version of the specified (or automatically found)
Manatee library and the version entered as the first argument (here 2.150) match together.

Script finishes in one of two possible result states:

* Two created files: *ictools* (a bootstrap script), *ictools.bin* (a binary executable). This means that *LD_LIBRARY_PATH* must be set for *ictools* because of a non-standard location of *libmanatee.so*. You can just copy these files to */usr/local/bin*
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
| A     | ictools         |  **164**                      | **55**           | **219**      |
| B     | classic scripts |  312                          | DNF (RAM)        | DNF          |
| B     | ictools         |  **175**                      | **63**           | **238**      |

Ictools are approximately **twice as fast** as the original Python scripts.

In terms of **memory usage**, there were no thorough measurements performed but according to the *top*
utility the *transalign* function in *ictools* consumes about **30-40% of of the memory** consumed
by the classic scripts. The import function (i.e. calign+fixgaps+compress) in both programs
consumes only a little RAM because data read from an input file are (almost) immediately written
to the output without any unnecessary memory allocation.

<a name="for_developers"></a>
## For developers

<a name="for_developers_setting_up_vscode"></a>
### Setting up VSCode debugging/testing environment

Open *debug* environment (left column) and click the "gear" button to edit *launch.json*. Then
set proper environment variables (just like in the previous paragraph).

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      " ....  parts are omitted here ... " : " ... ",
      "env": {
          "CGO_LDFLAGS": "-lmanatee -L/usr/local/lib",
          "CGO_CPPFLAGS": "-I/tmp/manatee-open-2.158.8"
      },
      " ....  parts are omitted here ... " : " ... ",
    }
  ]
}
```


<a name="for_developers_running_tests"></a>
### Running tests

To run the tests, add `--test` argument when running the *build* script. All the other parameters
must be set in the same way as when building the project. E.g.:

```
./build 2.150 --manatee-lib /opt/manatee-2.150/lib --test
```

To run  tests manually try to use `./build` script first to find out what are the values
of `CGO_LDFLAGS` and `CGO_CPPFLAGS` variables and then use them like this:

```
CGO_LDFLAGS="-lmanatee -L/usr/local/lib" CGO_CPPFLAGS="-I/tmp/manatee-open-2.158.8" go test ./...
```
