package ffs

import (
	"context"
	"syscall"

	"github.com/Doridian/fakerfs/util"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type NodeInterface interface {
	GetName() string

	fs.NodeSetattrer
	fs.NodeGetattrer
	fs.NodeOpener
}

type fsNode struct {
	fs.LoopbackNode

	isFake bool

	name      string
	handler   NodeInterface
	children  map[string]*fsNode
	childList []*fsNode
}

func (n *fsNode) isDir() bool {
	return n.handler == nil
}

func (n *fsNode) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	child := n.children[name]
	if child == nil {
		return n.LoopbackNode.Lookup(ctx, name, out)
	}

	attr := fs.StableAttr{
		Mode: fuse.S_IFREG,
	}
	if child.isDir() {
		attr.Mode = fuse.S_IFDIR
	}
	return n.NewInode(ctx, child, attr), fs.OK
}

func (n *fsNode) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	if n.isFake {
		if n.isDir() {
			return nil, 0, syscall.EISDIR
		}
		return n.handler.Open(ctx, flags)
	}

	return n.LoopbackNode.Open(ctx, flags)
}

func (n *fsNode) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if n.isFake && !n.isDir() {
		return n.handler.Getattr(ctx, fh, out)
	}

	superErr := n.LoopbackNode.Getattr(ctx, fh, out)

	if n.isFake && superErr == syscall.ENOENT {
		util.FillAttr(out)
		out.Mode = fuse.S_IFDIR | 0755

		return fs.OK
	}

	return superErr
}

func (n *fsNode) Setattr(ctx context.Context, fh fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	if n.isFake && !n.isDir() {
		return n.handler.Setattr(ctx, fh, in, out)
	}

	superErr := n.LoopbackNode.Setattr(ctx, fh, in, out)
	if n.isFake && superErr == syscall.ENOENT {
		util.FillAttr(out)
		out.Mode = fuse.S_IFDIR | 0755

		return fs.OK
	}

	return superErr
}

func (n *fsNode) Getxattr(ctx context.Context, attr string, dest []byte) (uint32, syscall.Errno) {
	if n.isFake && !n.isDir() {
		return 0, syscall.ENODATA
	}

	superRes, superErr := n.LoopbackNode.Getxattr(ctx, attr, dest)
	if n.isFake && superErr == syscall.ENOENT {
		return 0, syscall.ENODATA
	}

	return superRes, superErr
}

func (n *fsNode) Setxattr(ctx context.Context, attr string, data []byte, flags uint32) syscall.Errno {
	if n.isFake && !n.isDir() {
		return syscall.EPERM
	}

	superErr := n.LoopbackNode.Setxattr(ctx, attr, data, flags)
	if n.isFake && superErr == syscall.ENOENT {
		return syscall.EPERM
	}

	return superErr
}

func (n *fsNode) Removexattr(ctx context.Context, attr string) syscall.Errno {
	if n.isFake && !n.isDir() {
		return syscall.EPERM
	}

	superErr := n.LoopbackNode.Removexattr(ctx, attr)
	if n.isFake && superErr == syscall.ENOENT {
		return syscall.EPERM
	}

	return superErr
}

func (n *fsNode) Listxattr(ctx context.Context, dest []byte) (uint32, syscall.Errno) {
	if n.isFake && !n.isDir() {
		return 0, syscall.ENODATA
	}

	superRes, superErr := n.LoopbackNode.Listxattr(ctx, dest)
	if n.isFake && superErr == syscall.ENOENT {
		return 0, syscall.ENODATA
	}

	return superRes, superErr
}

func (n *fsNode) Opendir(ctx context.Context) syscall.Errno {
	if n.isFake {
		if !n.isDir() {
			return syscall.ENOTDIR
		}
		return fs.OK
	}

	return n.LoopbackNode.Opendir(ctx)
}

func (n *fsNode) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	if n.isFake && !n.isDir() {
		return nil, syscall.ENOTDIR
	}

	superDir, superErr := n.LoopbackNode.Readdir(ctx)
	if !n.isFake {
		return superDir, superErr
	}

	if superErr != fs.OK {
		superDir = nil
	}

	return newLister(superDir, n), fs.OK
}
