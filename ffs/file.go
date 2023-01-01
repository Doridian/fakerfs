package ffs

import (
	"context"
	"syscall"

	"github.com/Doridian/fakerfs/util"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type FileNode struct {
	fs.Inode
}

func (*FileNode) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if fh == nil {
		util.FillAttr(out)
		out.Mode = fuse.S_IFREG | 0644
		return fs.OK
	}
	return fh.(fs.FileGetattrer).Getattr(ctx, out)
}

func (f *FileNode) Setattr(ctx context.Context, fh fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	if fh == nil {
		return f.Getattr(ctx, nil, out)
	}
	return fh.(fs.FileSetattrer).Setattr(ctx, in, out)
}

func (*FileNode) GetStableMode() uint32 {
	return fuse.S_IFREG
}
