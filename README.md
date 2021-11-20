# dblp-fetcher

Use DBLP citations for BibTeX.

## Usage

```
$ cat titles.txt
Virtual Secure Platform: A Five-Stage Pipeline Processor over TFHE

$ cat titles.txt | ./dblp-fetcher
@inproceedings{DBLP:conf/uss/MatsuokaBMS021,
  author    = {Kotaro Matsuoka and
               Ryotaro Banno and
               Naoki Matsumoto and
               Takashi Sato and
               Song Bian},
  editor    = {Michael Bailey and
               Rachel Greenstadt},
  title     = {Virtual Secure Platform: {A} Five-Stage Pipeline Processor over {TFHE}},
  booktitle = {30th {USENIX} Security Symposium, {USENIX} Security 2021, August 11-13,
               2021},
  pages     = {4007--4024},
  publisher = {{USENIX} Association},
  year      = {2021},
  url       = {https://www.usenix.org/conference/usenixsecurity21/presentation/matsuoka},
  timestamp = {Thu, 16 Sep 2021 17:32:10 +0200},
  biburl    = {https://dblp.org/rec/conf/uss/MatsuokaBMS021.bib},
  bibsource = {dblp computer science bibliography, https://dblp.org}
}
```
