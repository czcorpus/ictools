#ifndef _FINLIB_CONFIG_HH
#define _FINLIB_CONFIG_HH 1
 
/* finlib/config.hh. Generated automatically at end of configure. */
/* finlib/config.hh.  Generated from config.hh.in by configure.  */
/* finlib/config.hh.in.  Generated from configure.ac by autoheader.  */

/* MacOS/DARWIN system */
/* #undef DARWIN */

/* Define to 1 if you have the <dlfcn.h> header file. */
#ifndef FINLIB_HAVE_DLFCN_H
#define FINLIB_HAVE_DLFCN_H 1
#endif

/* Define to 1 if you have the <fcntl.h> header file. */
#ifndef FINLIB_HAVE_FCNTL_H
#define FINLIB_HAVE_FCNTL_H 1
#endif

/* Define to 1 if fseeko (and presumably ftello) exists and is declared. */
#ifndef FINLIB_HAVE_FSEEKO
#define FINLIB_HAVE_FSEEKO 1
#endif

/* Define to 1 if you have the `getopt' function. */
#ifndef FINLIB_HAVE_GETOPT
#define FINLIB_HAVE_GETOPT 1
#endif

/* Define to 1 if you have the `getpagesize' function. */
#ifndef FINLIB_HAVE_GETPAGESIZE
#define FINLIB_HAVE_GETPAGESIZE 1
#endif

/* Define to 1 if you have the <inttypes.h> header file. */
#ifndef FINLIB_HAVE_INTTYPES_H
#define FINLIB_HAVE_INTTYPES_H 1
#endif

/* Define to 1 if you have the `memmove' function. */
#ifndef FINLIB_HAVE_MEMMOVE
#define FINLIB_HAVE_MEMMOVE 1
#endif

/* Define to 1 if you have the <memory.h> header file. */
#ifndef FINLIB_HAVE_MEMORY_H
#define FINLIB_HAVE_MEMORY_H 1
#endif

/* Define to 1 if you have a working `mmap' system call. */
#ifndef FINLIB_HAVE_MMAP
#define FINLIB_HAVE_MMAP 1
#endif

/* Define to 1 if you have the <netinet/in.h> header file. */
#ifndef FINLIB_HAVE_NETINET_IN_H
#define FINLIB_HAVE_NETINET_IN_H 1
#endif

/* Define to 1 if you have the `setlocale' function. */
#ifndef FINLIB_HAVE_SETLOCALE
#define FINLIB_HAVE_SETLOCALE 1
#endif

/* Define to 1 if `stat' has the bug that it succeeds when given the
   zero-length file name argument. */
/* #undef HAVE_STAT_EMPTY_STRING_BUG */

/* Define to 1 if you have the <stdint.h> header file. */
#ifndef FINLIB_HAVE_STDINT_H
#define FINLIB_HAVE_STDINT_H 1
#endif

/* Define to 1 if you have the <stdlib.h> header file. */
#ifndef FINLIB_HAVE_STDLIB_H
#define FINLIB_HAVE_STDLIB_H 1
#endif

/* Define to 1 if you have the `strchr' function. */
#ifndef FINLIB_HAVE_STRCHR
#define FINLIB_HAVE_STRCHR 1
#endif

/* Define to 1 if you have the `strdup' function. */
#ifndef FINLIB_HAVE_STRDUP
#define FINLIB_HAVE_STRDUP 1
#endif

/* Define to 1 if you have the `strerror' function. */
#ifndef FINLIB_HAVE_STRERROR
#define FINLIB_HAVE_STRERROR 1
#endif

/* Define to 1 if you have the <strings.h> header file. */
#ifndef FINLIB_HAVE_STRINGS_H
#define FINLIB_HAVE_STRINGS_H 1
#endif

/* Define to 1 if you have the <string.h> header file. */
#ifndef FINLIB_HAVE_STRING_H
#define FINLIB_HAVE_STRING_H 1
#endif

/* Define to 1 if you have the `strpbrk' function. */
#ifndef FINLIB_HAVE_STRPBRK
#define FINLIB_HAVE_STRPBRK 1
#endif

/* Define to 1 if you have the <sys/param.h> header file. */
#ifndef FINLIB_HAVE_SYS_PARAM_H
#define FINLIB_HAVE_SYS_PARAM_H 1
#endif

/* Define to 1 if you have the <sys/stat.h> header file. */
#ifndef FINLIB_HAVE_SYS_STAT_H
#define FINLIB_HAVE_SYS_STAT_H 1
#endif

/* Define to 1 if you have the <sys/types.h> header file. */
#ifndef FINLIB_HAVE_SYS_TYPES_H
#define FINLIB_HAVE_SYS_TYPES_H 1
#endif

/* Define to 1 if you have the <unistd.h> header file. */
#ifndef FINLIB_HAVE_UNISTD_H
#define FINLIB_HAVE_UNISTD_H 1
#endif

/* Define to 1 if `lstat' dereferences a symlink specified with a trailing
   slash. */
#ifndef FINLIB_LSTAT_FOLLOWS_SLASHED_SYMLINK
#define FINLIB_LSTAT_FOLLOWS_SLASHED_SYMLINK 1
#endif

/* Define to the sub-directory where libtool stores uninstalled libraries. */
#ifndef FINLIB_LT_OBJDIR
#define FINLIB_LT_OBJDIR ".libs/"
#endif

/* Name of package */
#ifndef FINLIB_PACKAGE
#define FINLIB_PACKAGE "finlib"
#endif

/* Define to the address where bug reports for this package should be sent. */
#ifndef FINLIB_PACKAGE_BUGREPORT
#define FINLIB_PACKAGE_BUGREPORT "pary@fi.muni.cz"
#endif

/* Define to the full name of this package. */
#ifndef FINLIB_PACKAGE_NAME
#define FINLIB_PACKAGE_NAME "finlib"
#endif

/* Define to the full name and version of this package. */
#ifndef FINLIB_PACKAGE_STRING
#define FINLIB_PACKAGE_STRING "finlib 2.36.5"
#endif

/* Define to the one symbol short name of this package. */
#ifndef FINLIB_PACKAGE_TARNAME
#define FINLIB_PACKAGE_TARNAME "finlib"
#endif

/* Define to the home page for this package. */
#ifndef FINLIB_PACKAGE_URL
#define FINLIB_PACKAGE_URL ""
#endif

/* Define to the version of this package. */
#ifndef FINLIB_PACKAGE_VERSION
#define FINLIB_PACKAGE_VERSION "2.36.5"
#endif

/* Define to 1 if you have the ANSI C header files. */
#ifndef FINLIB_STDC_HEADERS
#define FINLIB_STDC_HEADERS 1
#endif

/* Define to 1 if using ICU regexps */
/* #undef USE_ICU */

/* Define to 1 if using pcre regexps */
#ifndef FINLIB_USE_PCRE
#define FINLIB_USE_PCRE 1
#endif

/* Define to 1 if using standard regexps */
/* #undef USE_REGEX */

/* Version number of package */
#ifndef FINLIB_VERSION
#define FINLIB_VERSION "2.36.5"
#endif

/* Enable large inode numbers on Mac OS X 10.5.  */
#ifndef _DARWIN_USE_64_BIT_INODE
# define _DARWIN_USE_64_BIT_INODE 1
#endif

/* Number of bits in a file offset, on hosts where this is settable. */
/* #undef _FILE_OFFSET_BITS */

/* Define to 1 to make fseeko visible on some hosts (e.g. glibc 2.2). */
/* #undef _LARGEFILE_SOURCE */

/* Define for large files, on AIX-style hosts. */
/* #undef _LARGE_FILES */

/* Define to empty if `const' does not conform to ANSI C. */
/* #undef const */

/* Define to `__inline__' or `__inline' if that's what the C compiler
   calls it, or to nothing if 'inline' is not supported under any name.  */
#ifndef __cplusplus
/* #undef inline */
#endif

/* Define to `long int' if <sys/types.h> does not define. */
/* #undef off_t */

/* Define to `unsigned int' if <sys/types.h> does not define. */
/* #undef size_t */
 
/* once: _FINLIB_CONFIG_HH */
#endif
