package main

import (
	"io"
	"os"
	"fmt"
	"log"
	"os/user"
	"archive/tar"
	"compress/gzip"
	"path/filepath"
)

type Action func(*tar.Reader) error

func main() {
	var action byte
	args := os.Args
	zflag := false

	if len(args) < 2 || len(args[1]) > 2 || len(args[1]) < 1 { usage() }
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

	c_file := func(loc string, info os.FileInfo, _ error) error {
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
		} else { panic(err) }
	}
}

func targo(action Action, zflag bool) {

	var f io.Reader

	if zflag {
		gr, err := gzip.NewReader(os.Stdin)
		defer gr.Close()
		if err != nil { log.Fatal(err) }
		f = gr
	} else {
		f = os.Stdin
	}
	tr := tar.NewReader(f)
	if err := action(tr); err != nil { log.Fatal(err) }
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
	for {
		hdr, err := tr.Next()
		fi := hdr.FileInfo()
		u, err := user.Current()
		if err != nil { return err }
		if err == io.EOF { return nil }
		if err != nil { return err }

		switch hdr.Typeflag {
		case tar.TypeReg, tar.TypeRegA:
			f, err := os.OpenFile(hdr.Name, os.O_CREATE | os.O_WRONLY, fi.Mode())
			if err != nil {
				return err
			}
			defer f.Close()
			buf := make([]byte, hdr.Size) //bufio?
			if _, err := tr.Read(buf); err != nil { return err }
			if _, err := f.Write(buf); err != nil { return err }
		case tar.TypeLink:
			//hard link
		case tar.TypeSymlink:
		case tar.TypeChar:
		case tar.TypeBlock:
		case tar.TypeDir:
			if err := os.MkdirAll(hdr.Name, fi.Mode()); err != nil { return err }
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
		//if err := os.Chown(hdr.Name, hdr.Uid, hdr.Gid); err != nil { return err }
		if err := os.Chown(hdr.Name, u.Uid, u.Gid); err != nil { return err }
		if err := os.Chmod(hdr.Name, fi.Mode()); err != nil { return err }
	}
	return nil //shouldn't get here
}
