package ffs

import (
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
	rootNode *fsNode
	rootDir  *ffsDir

	server *fuse.Server
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

	ffs.rootNode = ffs.newNode(ffs.rootFS, nil, "", &st).(*fsNode)
	ffs.rootDir = &ffsDir{
		children:  map[string]*fsNode{},
		childList: []*fsNode{},
	}
	ffs.rootNode.handler = ffs.rootDir
	ffs.rootDir.node = ffs.rootNode

	ffs.rootNode.name = "[ROOT]"

	return ffs, nil
}

func (ffs *FakerFS) newNode(rootData *fs.LoopbackRoot, parent *fs.Inode, name string, st *syscall.Stat_t) fs.InodeEmbedder {
	return &fsNode{
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

	parent := ffs.rootDir

	lastIdx := len(pathElems) - 1

	for i := 0; i < len(pathElems); i++ {
		name := pathElems[i]
		newNode := parent.children[name]
		var newDir *ffsDir
		if newNode == nil {
			newDir = &ffsDir{
				children:  map[string]*fsNode{},
				childList: []*fsNode{},
			}
			newNode = &fsNode{
				handler: newDir,
				name:    name,
				LoopbackNode: fs.LoopbackNode{
					RootData: ffs.rootFS,
				},
			}
			newDir.node = newNode

			if i == lastIdx {
				newNode.handler = file
			}

			parent.children[name] = newNode
			parent.childList = append(parent.childList, newNode)
		} else {
			newDir = newNode.handler.(*ffsDir)
		}

		parent = newDir
	}
}
