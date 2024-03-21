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

ICTools come with [manabuild](https://github.com/czcorpus/manabuild) as its dependency. So in case you have
 `~/go/bin` in your `$PATH`, everything needed to build `ictools` is:

```
manabuild
```

In case Manabuild finds Manatee-open in a non-standard location where system does not look for libraries,
it produces `ictools.bin` with actual ICTools binary and `ictools` which is a short Bash script
to set `LD_LIBRARY_PATH` to the path Manabuild found Manatee in and to start the binary. So in this case,
two files must be moved (or copied) to a target installation location (e.g. `/usr/local/bin`).s



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

Run

```
manabuild -no-build
```

and copy `CGO_CPPFLAGS=...`, `CGO_CPPFLAGS=...` and `CGO_CXXFLAGS=...`.

Open *debug* environment (left column) and click the "gear" button to edit *launch.json*. Then
set proper environment variables (just like in the previous paragraph).

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      " ....  parts are omitted here ... " : " ... ",
      "env": {
          "CGO_LDFLAGS": "...",
          "CGO_CPPFLAGS": "...",
          "CGO_CXXFLAGS": "..."
      },
      " ....  parts are omitted here ... " : " ... ",
    }
  ]
}
```

Where the env. variables part is the one copied in the previous step.


<a name="for_developers_running_tests"></a>
### Running tests

```
manabuild -test
```