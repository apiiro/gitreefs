package fs

import (
	"github.com/jacobsa/fuse/fuseops"
	"gitreefs/fs"
)

type VersionedEntry struct {
	Inode
	size    int64
	isDir   bool
	entries []fuseops.InodeID
}

func (in *VersionedEntry) Attributes() fuseops.InodeAttributes {
	if in.isDir {
		return fs.DirAttributes()
	}
	return fs.FileAttributes(in.size)
}
