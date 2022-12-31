package ffs

import (
	"github.com/hanwen/go-fuse/v2/fs"
)

type NodeInterface interface {
	fs.InodeEmbedder
	fs.NodeGetattrer
}
