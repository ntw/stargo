package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

type Action func(*tar.Reader) error

func main() {
	var action byte
	args := os.Args
	zflag := false

	if len(args) < 2 || len(args[1]) > 2 || len(args[1]) < 1 {
		usage()
	}
	if args[1][0] == 'z' {
		if len(args[1]) != 2 {
			usage()
		} else {
			zflag = true
			action = args[1][1]
		}
	} else {
		action = args[1][0]
	}

	switch action {
	case 'c':
		c(args[2:], zflag)
	case 'x':
		targo(x, zflag)
	case 't':
		targo(t, zflag)
	default:
		usage()
	}
}

func usage() {
	log.Fatal("stargo [z][cxt] [files]\n")
}

func c(locs []string, zflag bool) {
	var f io.Writer
	if zflag {
		gw := gzip.NewWriter(os.Stdout)
		defer gw.Close()
		f = gw
	} else {
		f = os.Stdout
	}

	tw := tar.NewWriter(f)
	defer tw.Close()

	c_file := func(loc string, fi os.FileInfo, _ error) error {
		/*
		hdr, err := tar.FileInfoHeader(fi, loc)
		if err != nil {
			return err
		}
		
		target, err := os.Open(loc)
		if err != nil {
			return err
		}
		*/

		return nil
	}

	//error checking! also, what about symlinks?  Lstat or Stat?
	for _, loc := range locs {
		if fi, err := os.Stat(loc); err == nil {
			if fi.IsDir() {
				filepath.Walk(loc, c_file)
			} else {
				c_file(loc, fi, nil)
			}
		} else if err == filepath.SkipDir {
			log.Println(err)
		} else {
			panic(err)
		}
	}
}

func targo(action Action, zflag bool) {

	var f io.Reader

	if zflag {
		gr, err := gzip.NewReader(os.Stdin)
		defer gr.Close()
		if err != nil {
			log.Fatal(err)
		}
		f = gr
	} else {
		f = os.Stdin
	}
	tr := tar.NewReader(f)
	if err := action(tr); err != nil {
		log.Fatal(err)
	}
}

func t(tr *tar.Reader) error {

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", hdr.Name)
	}
	return nil //shouldn't get here
}

func x(tr *tar.Reader) error {
	u, err := user.Current() // as of know we overwrite the original user with the current.
	if err != nil {
		return err
	}
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fi := hdr.FileInfo()

		switch hdr.Typeflag {
		case tar.TypeReg, tar.TypeRegA:
			f, err := os.OpenFile(hdr.Name, os.O_CREATE|os.O_WRONLY, fi.Mode()) // don't clobber by default
			if err != nil {
				return err
			}
			if _, err := io.CopyN(f, tr, hdr.Size); err != nil {
				return err
			}
			if err := f.Close(); err != nil {
				return err
			}
		case tar.TypeLink:
			//hard link
		case tar.TypeSymlink:
			if err := os.Symlink(hdr.Linkname, hdr.Name); err != nil {
				return err
			}
			continue
		case tar.TypeChar:
		case tar.TypeBlock:
		case tar.TypeDir:
			if err := os.MkdirAll(hdr.Name, fi.Mode()); err != nil {
				return err
			}
		case tar.TypeFifo:
		case tar.TypeCont:
			//reserved
		case tar.TypeXHeader:
			//extended header
		case tar.TypeXGlobalHeader:
			//extended global header
		default:
			log.Printf("Unknown type for %s\n", hdr.Name)
		}
		//TODO: error checking
		uid, _ := strconv.Atoi(u.Uid)
		gid, _ := strconv.Atoi(u.Gid)
		if err := os.Chown(hdr.Name, uid, gid); err != nil {
			return err
		}
		//TODO: fix ModTime
		if err := os.Chmod(hdr.Name, fi.Mode()); err != nil {
			return err
		}
		if err := os.Chtimes(hdr.Name, fi.ModTime(), hdr.ModTime); err != nil {
			return err
		}
	}
	return nil //shouldn't get here
}
