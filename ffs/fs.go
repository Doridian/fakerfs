package ffs

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type FakerFS struct {
	rootPath string
	rootFS   *fs.LoopbackRoot
	rootNode *ffsDir

	server *fuse.Server
}

type fakerFSFileNode struct {
	fs.LoopbackNode
}

func (n *fakerFSFileNode) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	fh, fuseFlags, errno = n.LoopbackNode.Open(ctx, flags)
	if errno == fs.OK {
		fuseFlags |= fuse.FOPEN_DIRECT_IO
	}
	return
}

func NewFakerFS(rootPath string) (*FakerFS, error) {
	var st syscall.Stat_t
	err := syscall.Stat(rootPath, &st)
	if err != nil {
		return nil, err
	}

	ffs := &FakerFS{
		rootPath: rootPath,
	}

	ffs.rootFS = &fs.LoopbackRoot{
		Path:    rootPath,
		Dev:     uint64(st.Dev),
		NewNode: ffs.newNode,
	}

	ffs.rootNode = &ffsDir{
		LoopbackNode: fs.LoopbackNode{
			RootData: ffs.rootFS,
		},
		children:  map[string]NodeInterface{},
		childList: []string{},
	}

	return ffs, nil
}

func (ffs *FakerFS) newNode(rootData *fs.LoopbackRoot, parent *fs.Inode, name string, st *syscall.Stat_t) fs.InodeEmbedder {
	return &fakerFSFileNode{
		LoopbackNode: fs.LoopbackNode{
			RootData: rootData,
		},
	}
}

func (ffs *FakerFS) Mount(target string) error {
	opts := &fs.Options{}
	opts.AllowOther = true
	opts.DirectMount = true

	opts.MountOptions.Options = append(opts.MountOptions.Options, "default_permissions")

	opts.FsName = "ffs"
	opts.Name = "ffs"
	opts.MountOptions.Name = "ffs"

	server, err := fs.Mount(target, ffs.rootNode, opts)
	if err != nil {
		return err
	}

	ffs.server = server

	return nil
}

func (ffs *FakerFS) Wait() {
	ffs.server.Wait()
}

func (ffs *FakerFS) AddHandler(name string, file NodeInterface) {
	path := filepath.Clean(name)

	pathElems := strings.Split(path, string(os.PathSeparator))

	parent := ffs.rootNode

	lastIdx := len(pathElems) - 1

	for i := 0; i < len(pathElems); i++ {
		name := pathElems[i]
		newNode := parent.children[name]
		var newDir *ffsDir
		if newNode == nil {
			if i == lastIdx {
				newDir = nil
				newNode = file
			} else {
				newDir = &ffsDir{
					LoopbackNode: fs.LoopbackNode{
						RootData: ffs.rootFS,
					},
					children:  map[string]NodeInterface{},
					childList: []string{},
				}
				newNode = newDir
			}

			parent.children[name] = newNode
			parent.childList = append(parent.childList, name)
		} else {
			newDir = newNode.(*ffsDir)
		}

		parent = newDir
	}
}
