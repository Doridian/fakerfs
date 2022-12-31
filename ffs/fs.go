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

	sfs := &FakerFS{
		rootPath: rootPath,
	}

	sfs.rootFS = &fs.LoopbackRoot{
		Path:    rootPath,
		Dev:     uint64(st.Dev),
		NewNode: sfs.newNode,
	}

	sfs.rootNode = sfs.newNode(sfs.rootFS, nil, "", &st).(*fsNode)
	sfs.rootDir = &ffsDir{
		children:  map[string]*fsNode{},
		childList: []*fsNode{},
	}
	sfs.rootNode.handler = sfs.rootDir
	sfs.rootDir.node = sfs.rootNode

	sfs.rootNode.name = "[ROOT]"

	return sfs, nil
}

func (sfs *FakerFS) newNode(rootData *fs.LoopbackRoot, parent *fs.Inode, name string, st *syscall.Stat_t) fs.InodeEmbedder {
	return &fsNode{
		LoopbackNode: fs.LoopbackNode{
			RootData: rootData,
		},
	}
}

func (sfs *FakerFS) Mount(target string) error {
	opts := &fs.Options{}
	opts.AllowOther = true
	opts.DirectMount = true

	opts.MountOptions.Options = append(opts.MountOptions.Options, "default_permissions")

	opts.FsName = "sfs"
	opts.Name = "sfs"
	opts.MountOptions.Name = "sfs"

	server, err := fs.Mount(target, sfs.rootNode, opts)
	if err != nil {
		return err
	}

	sfs.server = server

	return nil
}

func (sfs *FakerFS) Wait() {
	sfs.server.Wait()
}

func (sfs *FakerFS) AddHandler(name string, file NodeInterface) {
	path := filepath.Clean(name)

	pathElems := strings.Split(path, string(os.PathSeparator))

	parent := sfs.rootDir

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
					RootData: sfs.rootFS,
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
