package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pallavagarwal07/fprint-fs/fpfs"
	"github.com/pallavagarwal07/mirror-fs/mfs"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: fprint-fs SOURCE MOUNTPOINT")
		return
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	root := &mfs.Server{
		Realpath:    os.Args[1],
		Transformer: &fpfs.Transformer{""},
		Debug:       true,
		Options:     []string{"allow_other", "default_permissions"},
	}
	if err := root.Mount(os.Args[2]); err != nil {
		log.Fatalln(err)
	}
}
