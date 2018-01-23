import re
import sys
import os
import subprocess
import argparse


KNOWN_VERSIONS = {
    '2.151.5': '2.36.5',
    '2.150': '2.36.5',
    '2.139.3': '2.35.2',
    '2.137.2': '2.35.1',
    '2.130.6': '2.33.1'
}

class NoManatee(object):

    @staticmethod
    def version():
        return None


def autodetect_version():
    try:
        import manatee
        v = manatee.version()
        return v[v.index('open-') + 5:]
    except ImportError:
        manatee = NoManatee()
        return manatee.version()


def _download_file(url, target):
    print('trying: {0}'.format(url))
    with open(target, 'wb') as fw:
        p = subprocess.Popen(['curl', '-fL', url], stdout=fw)
        p.wait()
    return p.returncode


def download_manatee_src(version):
    print('Looking for manatee sources ...')
    out_file = '/tmp/manatee-open-{0}.tar.gz'.format(version)
    if not os.path.exists(out_file):
        url = 'https://corpora.fi.muni.cz/noske/src/manatee-open/manatee-open-{0}.tar.gz'.format(version)
        ans = _download_file(url, out_file)
        if ans > 0:
            url = 'https://corpora.fi.muni.cz/noske/src/manatee-open/archive/manatee-open-{0}.tar.gz'.format(version)
            print('...failed.')
            ans = _download_file(url, out_file)
    else:
        print('...found in /tmp')
        ans = 0
    if ans == 0:
        ans = unpack_archive(out_file)
    if ans == 0:
        return '/tmp/manatee-open-{0}'.format(version)
    else:
        raise Exception('Failed to download and extract manatee. Please do this manually and run the script with --finlib-src ...')
    return ans


def download_finlib_src(version):
    print('Looking for finlib sources...')
    out_file = '/tmp/finlib-{0}.tar.gz'.format(version)
    if not os.path.exists(out_file):
        url = 'https://corpora.fi.muni.cz/noske/src/finlib/finlib-{0}.tar.gz'.format(version)
        ans = _download_file(url, out_file)
        if ans > 0:
            url = 'https://corpora.fi.muni.cz/noske/src/finlib/archive/finlib-{0}.tar.gz'.format(version)
            print('...failed.')
            ans = _download_file(url, out_file)
    else:
        print('...found in /tmp')
        ans = 0
    if ans == 0:
        ans = unpack_archive(out_file)
    if ans == 0:
        return '/tmp/finlib-{0}'.format(version)
    else:
        raise Exception('Failed to download and extract finlib. Please do this manually and run the script with --manatee-src ...')


def unpack_archive(path):
    p = subprocess.Popen(['tar', 'xzf', path, '-C', '/tmp'])
    p.wait()
    return p.returncode


def build_project(manatee_src, finlib_src, manatee_lib):
    p = subprocess.Popen([
        'CGO_CPPFLAGS="-I{0} -I{1}"'.format(manatee_src, finlib_src),
        'CGO_LDFLAGS="-lmanatee -L{0}"'.format(manatee_lib),
        'go build'
    ])
    p.wait()


def read_str_ans(msg):
    sys.stdout.write('{0}: '.format(msg))
    return raw_input()


def read_binary_ans(msg):
    sys.stdout.write('{0} (Y/N):'.format(msg))
    v = raw_input()
    if v in ('y', 'Y'):
        return True
    elif v in ('n', 'N'):
        return False
    else:
        print('Unknown answer - use Y/N')
        read_binary_ans('')

if __name__ == '__main__':
    argparser = argparse.ArgumentParser(description=None)
    argparser.add_argument('version', metavar="VERSION", help="Manatee version used along with ictools")
    argparser.add_argument('-f', '--finlib-src', type=str, help='Location of Finlib header files')
    argparser.add_argument('-m', '--manatee-src', type=str, help='Location of Manatee header files')
    argparser.add_argument('-M', '--manatee-lib', type=str, help='Location of libmanatee.so')
    args = argparser.parse_args()
    if args.version not in KNOWN_VERSIONS:
        print('Unsupported version: {0}. Please use one of: {1}'.format(args.version, ', '.join(sorted(KNOWN_VERSIONS.keys()))))
        sys.exit(1)
    print('')
    if not args.manatee_src:
        manatee_src = download_manatee_src(args.version)
    else:
        manatee_src = args.manatee_src
    if not args.finlib_src:
        finlib_src = download_finlib_src(KNOWN_VERSIONS[args.version])
    else:
        finlib_src = args.finlib_src
    if not args.manatee_lib:
        manatee_lib = autodetect_version()
        if manatee_lib is None:
            print('Manatee not found in system searched paths. Please run the script with --manatee-lib argument')
        elif args.version != manatee_lib:
            print('\nFound Manatee {0}, you require {1}.'.format(manatee_lib, args.version))
            print('Please specify a custom path with proper Manatee version (--manatee-lib)')
            print('or run the script with the version of the detected Manatee')
    else:
        manatee_lib = args.manatee_lib
    build_project(manatee_src, finlib_src, manatee_lib)

