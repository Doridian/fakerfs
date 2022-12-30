package ffs

import (
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

func (l *dotAndDotDotLister) Close() {

}

type fsLister struct {
	parent fs.DirStream
	node   *fsNode
	idx    int
}

func newLister(parent fs.DirStream, node *fsNode) fs.DirStream {
	if parent == nil {
		parent = &dotAndDotDotLister{idx: 0}
	}

	return &fsLister{
		parent: parent,
		node:   node,
		idx:    0,
	}
}

func (l *fsLister) hasFakeNext() bool {
	return l.idx < len(l.node.childList)
}

func (l *fsLister) HasNext() bool {
	return l.parent.HasNext() || l.hasFakeNext()
}

func (l *fsLister) Next() (fuse.DirEntry, syscall.Errno) {
	if l.parent.HasNext() {
		return l.parent.Next()
	}

	if !l.hasFakeNext() {
		return fuse.DirEntry{}, syscall.EINVAL
	}

	sNode := l.node.childList[l.idx]
	l.idx++

	dirEnt := fuse.DirEntry{
		Name: sNode.name,
		Mode: fuse.S_IFREG,
	}
	if sNode.handler == nil {
		dirEnt.Mode = fuse.S_IFDIR
	}
	return dirEnt, fs.OK
}

func (l *fsLister) Close() {

}
