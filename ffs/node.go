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

func (n *fsNode) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	child := n.children[name]
	if child == nil {
		return n.LoopbackNode.Lookup(ctx, name, out)
	}

	attr := fs.StableAttr{
		Mode: fuse.S_IFREG,
	}
	if child.handler == nil {
		attr.Mode = fuse.S_IFDIR
	}
	return n.NewInode(ctx, child, attr), fs.OK
}

func (n *fsNode) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	if n.isFake && n.handler != nil {
		return n.handler.Open(ctx, flags)
	}

	return n.LoopbackNode.Open(ctx, flags)
}

func (n *fsNode) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if n.isFake {
		if n.handler != nil {
			return n.handler.Getattr(ctx, fh, out)
		}

		util.FillAttr(out)
		out.Mode = fuse.S_IFDIR | 0755

		return fs.OK
	}

	return n.LoopbackNode.Getattr(ctx, fh, out)
}

func (n *fsNode) Setattr(ctx context.Context, fh fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	if n.isFake {
		if n.handler != nil {
			return n.handler.Setattr(ctx, fh, in, out)
		}

		util.FillAttr(out)
		out.Mode = fuse.S_IFDIR | 0755

		return fs.OK
	}

	return n.LoopbackNode.Setattr(ctx, fh, in, out)
}

func (n *fsNode) Getxattr(ctx context.Context, attr string, dest []byte) (uint32, syscall.Errno) {
	if n.isFake {
		return 0, syscall.ENODATA
	}

	return n.LoopbackNode.Getxattr(ctx, attr, dest)
}

func (n *fsNode) Setxattr(ctx context.Context, attr string, data []byte, flags uint32) syscall.Errno {
	if n.isFake {
		return syscall.EPERM
	}

	return n.LoopbackNode.Setxattr(ctx, attr, data, flags)
}

func (n *fsNode) Removexattr(ctx context.Context, attr string) syscall.Errno {
	if n.isFake {
		return syscall.EPERM
	}

	return n.LoopbackNode.Removexattr(ctx, attr)
}

func (n *fsNode) Listxattr(ctx context.Context, dest []byte) (uint32, syscall.Errno) {
	if n.isFake {
		return 0, syscall.ENODATA
	}

	return n.LoopbackNode.Listxattr(ctx, dest)
}

func (n *fsNode) Opendir(ctx context.Context) syscall.Errno {
	if n.isFake && n.handler == nil {
		return fs.OK
	}

	return n.LoopbackNode.Opendir(ctx)
}

func (n *fsNode) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	superDir, superErr := n.LoopbackNode.Readdir(ctx)

	if !n.isFake {
		return superDir, superErr
	}

	if n.handler != nil {
		return nil, syscall.ENOTDIR
	}

	if superErr != fs.OK {
		superDir = nil
	}

	return newLister(superDir, n), fs.OK
}
