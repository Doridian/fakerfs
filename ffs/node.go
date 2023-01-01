package ffs

import (
	"context"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type NodeInterface interface {
	fs.InodeEmbedder
	fs.NodeGetattrer
	GetStableMode() uint32
}

type fakerFSFileNode struct {
	fs.LoopbackNode
}

func (*fakerFSFileNode) Mknod(ctx context.Context, name string, mode, rdev uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, syscall.EPERM
}

func (*fakerFSFileNode) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, syscall.EPERM
}

func (*fakerFSFileNode) Rmdir(ctx context.Context, name string) syscall.Errno {
	return syscall.EPERM
}

func (*fakerFSFileNode) Unlink(ctx context.Context, name string) syscall.Errno {
	return syscall.EPERM
}

func (*fakerFSFileNode) Rename(ctx context.Context, name string, newParent fs.InodeEmbedder, newName string, flags uint32) syscall.Errno {
	return syscall.EPERM
}

func (*fakerFSFileNode) Create(ctx context.Context, name string, flags uint32, mode uint32, out *fuse.EntryOut) (inode *fs.Inode, fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	return nil, nil, 0, syscall.EPERM
}

func (*fakerFSFileNode) Symlink(ctx context.Context, target, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, syscall.EPERM
}

func (*fakerFSFileNode) Link(ctx context.Context, target fs.InodeEmbedder, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	return nil, syscall.EPERM
}
