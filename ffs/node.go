package ffs

import (
	"context"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type NodeInterface interface {
	fs.NodeSetattrer
	fs.NodeGetattrer
	fs.NodeOpener
	fs.NodeLookuper

	fs.NodeOpendirer
	fs.NodeReaddirer

	fs.NodeSetxattrer
	fs.NodeGetxattrer
	fs.NodeListxattrer
	fs.NodeRemovexattrer
}

type fsNode struct {
	fs.LoopbackNode

	name    string
	handler NodeInterface
}

func (n *fsNode) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	if n.handler != nil {
		return n.handler.Lookup(ctx, name, out)
	}

	return n.LoopbackNode.Lookup(ctx, name, out)
}

func (n *fsNode) Open(ctx context.Context, flags uint32) (fh fs.FileHandle, fuseFlags uint32, errno syscall.Errno) {
	if n.handler != nil {
		return n.handler.Open(ctx, flags)
	}

	return n.LoopbackNode.Open(ctx, flags)
}

func (n *fsNode) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if n.handler != nil {
		return n.handler.Getattr(ctx, fh, out)
	}

	return n.LoopbackNode.Getattr(ctx, fh, out)
}

func (n *fsNode) Setattr(ctx context.Context, fh fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	if n.handler != nil {
		return n.handler.Setattr(ctx, fh, in, out)
	}

	return n.LoopbackNode.Setattr(ctx, fh, in, out)
}

func (n *fsNode) Getxattr(ctx context.Context, attr string, dest []byte) (uint32, syscall.Errno) {
	if n.handler != nil {
		return n.handler.Getxattr(ctx, attr, dest)
	}

	return n.LoopbackNode.Getxattr(ctx, attr, dest)
}

func (n *fsNode) Setxattr(ctx context.Context, attr string, data []byte, flags uint32) syscall.Errno {
	if n.handler != nil {
		return n.handler.Setxattr(ctx, attr, data, flags)
	}

	return n.LoopbackNode.Setxattr(ctx, attr, data, flags)
}

func (n *fsNode) Removexattr(ctx context.Context, attr string) syscall.Errno {
	if n.handler != nil {
		return n.handler.Removexattr(ctx, attr)
	}

	return n.LoopbackNode.Removexattr(ctx, attr)
}

func (n *fsNode) Listxattr(ctx context.Context, dest []byte) (uint32, syscall.Errno) {
	if n.handler != nil {
		return n.handler.Listxattr(ctx, dest)
	}

	return n.LoopbackNode.Listxattr(ctx, dest)
}

func (n *fsNode) Opendir(ctx context.Context) syscall.Errno {
	if n.handler != nil {
		return n.handler.Opendir(ctx)
	}

	return n.LoopbackNode.Opendir(ctx)
}

func (n *fsNode) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	if n.handler != nil {
		return n.handler.Readdir(ctx)
	}

	return n.LoopbackNode.Readdir(ctx)
}
