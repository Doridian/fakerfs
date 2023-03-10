package sysfs

import (
	"context"
	"syscall"
	"time"

	"github.com/Doridian/fakerfs/ffs"
	"github.com/Doridian/fakerfs/util"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type FileHandler interface {
	GetData() ([]byte, syscall.Errno)
	SetData([]byte) syscall.Errno
	LoadConfig(map[string]interface{}) error
}

type fsNode struct {
	ffs.FileNode

	handler FileHandler
	mtime   uint64

	readData  []byte
	readErrno syscall.Errno
	cacheData bool
}

var _ ffs.NodeInterface = &fsNode{}

type fileHandle struct {
	readData  []byte
	readErrno syscall.Errno

	fs           *fsNode
	currentState []byte
}

func MakeFile(handler FileHandler, cacheData bool) *fsNode {
	file := &fsNode{
		handler:   handler,
		mtime:     0,
		cacheData: cacheData,
	}

	file.Refresh()

	return file
}

func (f *fsNode) Refresh() {
	if !f.cacheData {
		return
	}
	f.readData, f.readErrno = f.handler.GetData()
}

func (f *fsNode) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	return f.MakeFileHandle(), fuse.FOPEN_KEEP_CACHE | fuse.FOPEN_NONSEEKABLE, fs.OK
}

func (f *fsNode) MakeFileHandle() *fileHandle {
	fh := &fileHandle{
		currentState: []byte{},
		fs:           f,
	}

	if f.cacheData {
		fh.readData = f.readData
		fh.readErrno = f.readErrno
	} else {
		fh.readData, fh.readErrno = f.handler.GetData()
	}

	return fh
}

func (f *fileHandle) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	if f.readErrno != fs.OK {
		return nil, f.readErrno
	}

	end := int(off) + len(dest)
	if end > len(f.readData) {
		end = len(f.readData)
	}

	return fuse.ReadResultData(f.readData[off:end]), fs.OK
}

func (f *fileHandle) Write(ctx context.Context, data []byte, off int64) (uint32, syscall.Errno) {
	if off != 0 {
		return 0, syscall.EINVAL
	}

	res := f.fs.handler.SetData(data)
	if res != fs.OK {
		return 0, res
	}

	f.fs.mtime = uint64(time.Now().Unix())

	f.fs.Refresh()

	return uint32(len(data)), fs.OK
}

func (f *fileHandle) Flush(ctx context.Context) syscall.Errno {
	return fs.OK
}

func (f *fileHandle) Getattr(ctx context.Context, out *fuse.AttrOut) syscall.Errno {
	util.FillAttr(out)
	out.Mode = fuse.S_IFREG | 0644

	if f.fs.mtime > 0 {
		out.Mtime = f.fs.mtime
		out.Mtimensec = 0
		out.Atime = f.fs.mtime
		out.Atimensec = 0
	}

	if f.readErrno != fs.OK {
		return fs.OK
	}

	out.Size = uint64(len(f.readData))
	out.Blocks = (out.Size / uint64(out.Blksize)) + 1
	return fs.OK
}

func (f *fileHandle) Setattr(ctx context.Context, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {
	return f.Getattr(ctx, out)
}
