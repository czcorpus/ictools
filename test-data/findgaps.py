import sys
import manatee


def process_range(s):
    tmp = s.split(',')
    if len(tmp) == 2:
        return int(tmp[0]), int(tmp[1])
    return int(tmp[0]), int(tmp[0])


def analyze_range(v1, v2):
    if v2 - v1 > 100:
        return v2 - v1
    return 0


def process_line(s):
    items = s.split('\t')
    l1, l2 = process_range(items[0])
    r1, r2 = process_range(items[1])
    return (l1, l2, r1, r2)

def analyze_line(l1, l2, r1, r2):
    lft = analyze_range(l1, l2)
    rgt = analyze_range(r1, r2)
    if (lft > 0 or rgt > 0) and l1 > -1 and r1 > -1:
        print('{0},{1} ({2}) --> {3},{4} ({5})'.format(l1, l2, lft, r1, r2, rgt))

class TestAlign(object):

    def __init__(self, corp1, corp2, align_path):
        self.attr1 = corp1.get_attr('s.id')
        self.attr2 = corp2.get_attr('s.id')
        self.input_path = align_path

    def getids1(self, v1, v2):
        return self.attr1.id2str(v1), self.attr1.id2str(v2)

    def getids2(self, v1, v2):
        return self.attr2.id2str(v1), self.attr2.id2str(v2)

    def get_text_id(self, v):
        items = v.split(':')
        return items[1] if len(items) >= 2 else None

    def run(self):
        with open(self.input_path) as fr:
            for line in fr:
                d = process_line(line)
                id11, id12 = self.getids1(d[0], d[1])
                id21, id22 = self.getids2(d[2], d[3])
                if True or d[0] > 0 and d[1] > 0 and d[2] > 0 and d[3] > 0:
                    text11 = self.get_text_id(id11)
                    text12 = self.get_text_id(id12)
                    text21 = self.get_text_id(id21)
                    text22 = self.get_text_id(id22)
                    #v = set(filter(lambda x: x is not None, [text11, text12, text21, text22]))
                    #if len(v) > 1:
                    print(text11, text12, text21, text22)


if __name__ == '__main__':
    with open(sys.argv[1], 'rb') as fr:
        for line in fr:
            process_line(line)
    #c1 = manatee.Corpus('/home/tomas/corpora/registry/intercorp_v11_en')
    #c2 = manatee.Corpus('/home/tomas/corpora/registry/intercorp_v11_pl')
    #t = TestAlign(c1, c2, sys.argv[1])
    #t.run()