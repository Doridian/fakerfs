package ffs

import (
	"context"
	"syscall"

	"github.com/Doridian/fakerfs/util"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type SimpleNode struct {
}

func (*SimpleNode) Setattr(ctx context.Context, f fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	return syscall.EPERM
}

func (*SimpleNode) Getattr(ctx context.Context, f fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	util.FillAttr(out)
	out.Mode = fuse.S_IFREG | 0644
	return fs.OK
}

func (*SimpleNode) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, syscall.ENOTDIR
}

func (*SimpleNode) Opendir(ctx context.Context) syscall.Errno {
	return syscall.ENOTDIR
}

func (*SimpleNode) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	return nil, syscall.ENOTDIR
}

func (*SimpleNode) Setxattr(ctx context.Context, attr string, data []byte, flags uint32) syscall.Errno {
	return syscall.EPERM
}

func (*SimpleNode) Getxattr(ctx context.Context, attr string, dest []byte) (uint32, syscall.Errno) {
	return 0, syscall.ENODATA
}

func (*SimpleNode) Listxattr(ctx context.Context, dest []byte) (uint32, syscall.Errno) {
	return 0, syscall.ENODATA
}

func (*SimpleNode) Removexattr(ctx context.Context, attr string) syscall.Errno {
	return syscall.EPERM
}

func (*SimpleNode) Readlink(ctx context.Context) ([]byte, syscall.Errno) {
	return nil, syscall.ENOLINK
}
