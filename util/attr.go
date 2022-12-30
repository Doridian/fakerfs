package util

import (
	"time"

	"github.com/hanwen/go-fuse/v2/fuse"
)

var startTime = uint64(time.Now().Unix())

func FillAttr(attr *fuse.AttrOut) {
	attr.Atime = startTime
	attr.Mtime = startTime
	attr.Ctime = startTime
	attr.Atimensec = 0
	attr.Mtimensec = 0
	attr.Ctimensec = 0

	attr.Nlink = 1

	attr.Uid = 0
	attr.Gid = 0

	attr.Blocks = 1
	attr.Blksize = 4096
	attr.Padding = 0
}
