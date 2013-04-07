# Stargo
## Stargo is a minimalist tar implementation in go.

Stargo depends on an updated version of the tar library.  You can find more info here: https://code.google.com/p/go/source/detail?r=de92672228d3
It aims to provide the basic functionality of listing the contents of a tar file, as well as creating and extracting tar files.  It currently has support for gzipped files.  Bzip2 support is planned... once we have a bzip2 writer library for go.

Basic usage:
- To list the contents of a tar file: stargo t < tarfile.tar
- To extract a tar file to the current directory: stargo x < tarfile.tar

if the file is gzipped use "stargo zt" or stargo "xt".

At this point tar file creation has yet to be implemented, although most of the scaffolding is in place.  If you attempt to extract tar files with any of the more unusual filetypes supported by tar (such as block devices) you won't have much luck yet.
