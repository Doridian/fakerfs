package ffs

import (
	"context"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type dotAndDotDotLister struct {
	idx int
}

func (l *dotAndDotDotLister) HasNext() bool {
	return l.idx < 2
}

func (l *dotAndDotDotLister) Next() (fuse.DirEntry, syscall.Errno) {
	if l.idx >= 2 {
		return fuse.DirEntry{}, syscall.EINVAL
	}

	name := "."
	if l.idx == 1 {
		name = ".."
	}
	l.idx++

	return fuse.DirEntry{
		Name: name,
		Mode: fuse.S_IFDIR,
	}, fs.OK
}

func (l *dotAndDotDotLister) Close() {}

type ffsDirLister struct {
	parent fs.DirStream
	dir    *ffsDir
	idx    int
}

func newLister(parent fs.DirStream, dir *ffsDir) fs.DirStream {
	if parent == nil {
		parent = &dotAndDotDotLister{idx: 0}
	}

	return &ffsDirLister{
		parent: parent,
		dir:    dir,
		idx:    0,
	}
}

func (l *ffsDirLister) hasFakeNext() bool {
	return l.idx < len(l.dir.childList)
}

func (l *ffsDirLister) HasNext() bool {
	return l.parent.HasNext() || l.hasFakeNext()
}

func (l *ffsDirLister) Next() (fuse.DirEntry, syscall.Errno) {
	for l.parent.HasNext() {
		nextTry, err := l.parent.Next()
		if err != fs.OK {
			break
		}
		// Make sure we only list real files if there is no fakes
		if _, ok := l.dir.children[nextTry.Name]; ok {
			continue
		}
		return nextTry, fs.OK
	}

	if !l.hasFakeNext() {
		return fuse.DirEntry{}, syscall.EINVAL
	}

	sNode := l.dir.childList[l.idx]
	l.idx++

	dirEnt := fuse.DirEntry{
		Name: sNode.name,
		Mode: fuse.S_IFREG,
	}

	attrOut := fuse.AttrOut{}
	sNode.Getattr(context.Background(), nil, &attrOut)
	dirEnt.Mode = attrOut.Attr.Mode

	return dirEnt, fs.OK
}

func (l *ffsDirLister) Close() {
	l.parent.Close()
}
