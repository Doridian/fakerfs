package ffs

import (
	"context"
	"syscall"

	"github.com/Doridian/fakerfs/util"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type ffsDir struct {
	fs.LoopbackNode

	children  map[string]NodeInterface
	childList []string
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
		return d.LoopbackNode.Lookup(ctx, name, out)
	}

	subAttr := fuse.AttrOut{}
	child.Getattr(ctx, nil, &subAttr)
	out.Attr = subAttr.Attr
	attr := fs.StableAttr{
		Mode: subAttr.Attr.Mode,
	}

	return d.NewInode(ctx, child, attr), fs.OK
}

func (*ffsDir) Opendir(ctx context.Context) syscall.Errno {
	return fs.OK
}

func (d *ffsDir) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	superDir, superErr := d.LoopbackNode.Readdir(ctx)

	if superErr != fs.OK {
		superDir = nil
	}

	return newLister(superDir, d), fs.OK
}

func (*ffsDir) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	return syscall.EPERM
}

func (*ffsDir) Setxattr(ctx context.Context, attr string, data []byte, flags uint32) syscall.Errno {
	return syscall.EPERM
}

func (*ffsDir) Getxattr(ctx context.Context, attr string, dest []byte) (uint32, syscall.Errno) {
	return 0, syscall.ENODATA
}

func (*ffsDir) Listxattr(ctx context.Context, dest []byte) (uint32, syscall.Errno) {
	return 0, syscall.ENODATA
}

func (*ffsDir) Removexattr(ctx context.Context, attr string) syscall.Errno {
	return syscall.EPERM
}

func (*ffsDir) Readlink(ctx context.Context) ([]byte, syscall.Errno) {
	return nil, syscall.EINVAL
}
