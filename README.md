# Stargo
## Stargo is a minimalist tar implementation in go.

`stargo` depends on an updated version of the tar library.  You can find more info here: https://code.google.com/p/go/source/detail?r=de92672228d3
It aims to provide the basic functionality of listing the contents of a tar file, as well as creating and extracting tar files.  It currently has support for gzipped files.  Bzip2 support is planned... once we have a bzip2 writer library for go.

Basic usage:
- `stargo [options] [pathname ...]`
- available options are `z` for gzipped files, `c` for create, `x` for extract, and `t` for list.

Example usage:
- To list the contents of a gzipped tar file: `stargo zt < tarfile.tar`
- To extract a tar file to the current directory: `stargo x < tarfile.tar`
- To create a gzipped archive: `stargo zc somefile somedir/ ... > file.tar.gz

If you attempt to extract or create tar files with any of the more unusual filetypes supported by tar (such as block devices) you won't have much luck yet.
