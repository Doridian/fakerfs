package ffs

import (
	"context"
	"syscall"

	"github.com/Doridian/fakerfs/util"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type ffsDir struct {
	SimpleNode

	children  map[string]*fsNode
	childList []*fsNode

	node *fsNode
}

var _ NodeInterface = &ffsDir{}

func (*ffsDir) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	util.FillAttr(out)
	out.Mode = fuse.S_IFDIR | 0755
	return fs.OK
}

func (*ffsDir) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	return nil, 0, syscall.EISDIR
}

func (d *ffsDir) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	child := d.children[name]
	if child == nil {
		return d.node.LoopbackNode.Lookup(ctx, name, out)
	}

	subAttr := fuse.AttrOut{}
	child.handler.Getattr(ctx, nil, &subAttr)
	attr := fs.StableAttr{
		Mode: subAttr.Attr.Mode,
	}
	out.Attr = subAttr.Attr

	return d.node.NewInode(ctx, child, attr), fs.OK
}

func (*ffsDir) Opendir(ctx context.Context) syscall.Errno {
	return fs.OK
}

func (d *ffsDir) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	superDir, superErr := d.node.LoopbackNode.Readdir(ctx)

	if superErr != fs.OK {
		superDir = nil
	}

	return newLister(superDir, d), fs.OK
}
