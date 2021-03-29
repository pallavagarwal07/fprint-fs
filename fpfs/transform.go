package fpfs

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pallavagarwal07/libfprint-go/fprint"
	"github.com/pallavagarwal07/libfprint-go/prompt"
	"github.com/pallavagarwal07/mirror-fs/mfs"
)

type Transformer struct {
	Pass string
}

type FileInfo struct {
	name string
	mode os.FileMode
}

func (fi *FileInfo) Name() string {
	return fi.name
}
func (fi *FileInfo) Mode() os.FileMode {
	return fi.mode
}

var _ mfs.AttrTransformer = (*Transformer)(nil)
var _ mfs.DataTransformer = (*Transformer)(nil)
var _ mfs.ReverseTransformer = (*Transformer)(nil)
var _ mfs.FileInfo = (*FileInfo)(nil)

func (*Transformer) AttrTransform(c mfs.Ctx, fi os.FileInfo) (mfs.FileInfo, error) {
	return &FileInfo{
		name: strings.TrimSuffix(fi.Name(), ".crypt"),
		mode: fi.Mode(),
	}, nil
}

func (t *Transformer) DataTransform(c mfs.Ctx, data []byte) (out []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Println(e)
			err = syscall.EACCES
		}
	}()

	if p := filepath.Join(c.Basepath(), "encKey"); p == c.Realpath() {
		return []byte{}, nil
	}

	conn, err := fprint.NewConn("pallav")
	if err != nil {
		log.Println(err)
		return nil, syscall.EBUSY
	}
	defer conn.Close()

	fpChan, err := conn.StartVerification(5, 10)
	if err != nil {
		log.Println(err)
		return nil, syscall.EBUSY
	}
	passwdChan := make(chan string, 10)
	p := prompt.NewPrompt(passwdChan)

loop:
	for {
		fpRes := fprint.VERIFY_FAILED
		select {
		case passwd, ok := <-passwdChan:
			if !ok {
				break loop
			}
			if passwd == t.Pass {
				fpRes = fprint.VERIFY_SUCCESS
			}
		case result, ok := <-fpChan:
			if !ok {
				fpChan = nil
				continue loop
			}
			fpRes = result
		}
		p.Result <- fpRes
		if fpRes == fprint.VERIFY_SUCCESS {
			return []byte(Decrypt([]byte(t.Pass), string(data))), nil
		}
	}
	return nil, syscall.EACCES
}

func (t *Transformer) ReverseTransform(c mfs.Ctx, data []byte) (out []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Println(err)
			err = fmt.Errorf("%v", e)
		}
	}()

	if p := filepath.Join(c.Basepath(), "encKey"); p == c.Realpath() {
		t.Pass = string(data)
		return []byte{}, nil
	}

	return []byte(Encrypt([]byte(t.Pass), string(data))), nil
}
