import argparse
import sys
import os.path

sys.path.insert(0, '/opt/manatee/2.130.6/lib/python2.7/site-packages/')
import manatee

class Alignment(object):

    def __init__(self, l1, l2, r1, r2):
        self.l1 = l1
        self.l2 = l2
        self.r1 = r1
        self.r2 = r2

    def __repr__(self):
        return '{0},{1} --> {2},{3}'.format(self.l1, self.l2, self.r1, self.r2)

class AligndefFile(object):

    def __init__(self, file_path):
        self._data = []
        with open(file_path, 'rb') as fr:
            for line in fr:
                tmp = line.strip()
                if len(tmp) > 0:
                    aitem = self._parse_line(tmp)
                    if aitem.l1 > -1:
                        self._data.append(aitem)

    def _parse_line(self, s):
        items = s.split('\t')
        lft_v = [int(x) for x in items[0].split(',')]
        rgt_v = [int(x) for x in items[1].split(',')]

        if len(lft_v) == 1:
            lft_v.append(lft_v[0])
        if len(rgt_v) == 1:
            rgt_v.append(rgt_v[0])
        return  Alignment(lft_v[0], lft_v[1], rgt_v[0], rgt_v[1])

    def find_left_val(self, v):
        return self._find_left_val(v, 0, len(self._data))

    def _find_left_val(self, v, lft, rgt):
        pivot = (rgt + lft) / 2
        #print('range [{0} .... {1} .... {2}]'.format(lft, pivot, rgt))
        #print('data around: {0}'.format(self._data[pivot-2:pivot+2]))
        if v < self._data[lft].l1 or v > self._data[rgt-1].l2:
            return None
        if v < self._data[pivot].l1:
            ans = self._find_left_val(v, lft, pivot)
        elif v > self._data[pivot].l2:
            ans = self._find_left_val(v, pivot, rgt)
        else:
            ans = self._data[pivot]
        return ans


class CorpusProvider(object):

    def __init__(self, reg_path=''):
        self._reg_path = reg_path

    def open_corpus(self, ident):
        p = os.path.join(self._reg_path, ident)
        return manatee.Corpus(p)


class Validator(object):

    def __init__(self, corp1, corp2, aligndef_file, struct_name):
        self._corp1 = corp1
        self._corp2 = corp2
        self._attr1 = corp1.get_attr('s.id')  # TODO configurability
        self._aligndef_file = aligndef_file

    def run(self):
        attr = self._corp1.get_struct('doc')
        print(attr)
        for i in range(0, attr.size()):
            find_struct_begin(self._corp1, self._aligndef_file, self._attr1, 'doc', i)


def _parse_refs(s):
    items = s.split(',')
    struct_id = items[0].split('=')[1]
    token_idx = items[1][1:]
    return struct_id, int(token_idx)


def _find_refs(conc, attr, alignment, idx):
    limit = 1
    leftcontext = '-1'
    rightcontext = '1'
    attrs = ''
    attrs_allpos = ''
    structs = ''
    refs = 's.id,#'
    maxcontext = 10
    kw = manatee.KWICLines(conc.corp(), conc.RS(True, 0, limit), leftcontext, rightcontext, attrs, attrs_allpos, structs, refs, maxcontext)
    while kw.nextline():
        refs = kw.get_refs()
        struct_id, token_idx = _parse_refs(refs)
        sent_order = attr.str2id(struct_id)
        srch = alignment.find_left_val(sent_order)
        print('#{0} -- {1} -- 1st sentence in corp: {2} -- aligndef line: {3}'.format(idx, struct_id, sent_order, srch))


def find_struct_begin(corp, alignment, sentence_attr, struct_name, struct_idx):
    conc = manatee.Concordance(corp, '<{0} #{1}>[]'.format(struct_name, struct_idx), 0, -1)
    conc.sync()
    if conc.size() != 1:
        print('ERROR: <{0} #{1}> not found'.format(struct_name, struct_idx))
    _find_refs(conc, sentence_attr, alignment, struct_idx)
    return None


def run(reg_path, lang1, lang2, aligndef_file, struct_name):
    cp = CorpusProvider(reg_path=reg_path)
    corp1 = cp.open_corpus(lang1)
    corp2 = cp.open_corpus(lang2)
    af = AligndefFile(aligndef_file)
    v = Validator(corp1=corp1, corp2=corp2, aligndef_file=af, struct_name=struct_name)
    v.run()


if __name__ == '__main__':
    argparser = argparse.ArgumentParser(description="Alignment validator")
    argparser.add_argument('corp1', metavar="CORP1", help="")
    argparser.add_argument('corp2', metavar="CORP2", help="")
    argparser.add_argument('aligndef_file', metavar="ALIGNDEF_FILE", help="")
    argparser.add_argument('-r', '--registry-path', type=str)
    argparser.add_argument('-s', '--struct-name', type=str, default='doc')
    args = argparser.parse_args()
    run(reg_path=args.registry_path, lang1=args.corp1, lang2=args.corp2,
        aligndef_file=args.aligndef_file, struct_name=args.struct_name)
