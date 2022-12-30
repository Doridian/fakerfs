package ffs

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/Doridian/fakerfs/dev"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type FakerFS struct {
	rootPath string
	rootFS   *fs.LoopbackRoot
	rootNode *fsNode

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

	sfs.rootNode.isFake = true
	sfs.rootNode.name = "[ROOT]"
	sfs.rootNode.children = map[string]*fsNode{}
	sfs.rootNode.childList = []*fsNode{}

	return sfs, nil
}

func (sfs *FakerFS) newNode(rootData *fs.LoopbackRoot, parent *fs.Inode, name string, st *syscall.Stat_t) fs.InodeEmbedder {
	return &fsNode{
		isFake: false,

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

func (sfs *FakerFS) AddHandler(file *dev.FileFuse) {
	path := filepath.Clean(file.GetName())

	pathElems := strings.Split(path, string(os.PathSeparator))

	parent := sfs.rootNode

	lastIdx := len(pathElems) - 1

	for i := 0; i < len(pathElems); i++ {
		name := pathElems[i]
		newNode := parent.children[name]
		if newNode == nil {
			newNode = &fsNode{
				isFake:  true,
				handler: nil,
				name:    name,

				children:  map[string]*fsNode{},
				childList: []*fsNode{},

				LoopbackNode: fs.LoopbackNode{
					RootData: sfs.rootFS,
				},
			}

			if i == lastIdx {
				newNode.handler = file
				newNode.children = nil
				newNode.childList = nil
			}

			parent.children[name] = newNode
			parent.childList = append(parent.childList, newNode)
		}

		parent = newNode
	}
}
