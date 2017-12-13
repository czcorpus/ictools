//  Copyright (c) 2004-2016  Pavel Rychly, Milos Jakubicek

#ifndef VIRTCORP_HH
#define VIRTCORP_HH

#include <finlib/fstream.hh>
#include <finlib/range.hh>
#include <vector>
#include <string>

class Corpus;
class PosAttr;
template <class NormClass=MapBinFile<int64_t>,
          class FreqClass=MapBinFile<uint32_t>,
          class FloatFreqClass=MapBinFile<float> >
class VirtualPosAttr;

class VirtualCorpus {
public:
    struct PosTrans {
        // transition from original positions to new positions,
        // the size of the included region is computed as a difference of
        // next newpos and current newpos (in the postrans array)
        Position orgpos;
        Position newpos;
        PosTrans (Position o, Position n) : orgpos(o), newpos(n){}
    };
    struct Segment {
        Corpus *corp;
        // last item in postrans indicates the end of the last
        // transition region and the start position of the next segment
        std::vector<PosTrans> postrans;
    };
    std::vector<Segment> segs;
    //string path;

    virtual Position size() {return segs.back().postrans.back().newpos;}
    virtual ~VirtualCorpus() {}
    FastStream *combine_poss (VirtualPosAttr<> *vpa, std::vector<FastStream*> &fsv);
    //VirtualCorpus(const string &path) : path(path){}
};

VirtualCorpus* setup_virtcorp (const std::string &filename);
PosAttr* setup_virtposattr (VirtualCorpus *vc, const std::string &path,
                            const std::string &name, const std::string &locale,
                            const std::string &enc, bool ownedByCorpus=true,
                            const std::string &def="", const std::string &doc="");
ranges* setup_virtstructrng (VirtualCorpus *vc, const std::string &name);
VirtualCorpus* virtcorp2virtstruc (VirtualCorpus *vc, const std::string &name);

#endif

// vim: ts=4 sw=4 sta et sts=4 si cindent tw=80:
