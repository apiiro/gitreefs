package fs

import (
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"sync/atomic"
)

var (
	allocatedInodeId uint64 = fuseops.RootInodeID
)

func NextInodeID() (next fuseops.InodeID) {
	nextInodeId := atomic.AddUint64(&allocatedInodeId, 1)
	return fuseops.InodeID(nextInodeId)
}

type Inode interface {
	Id() fuseops.InodeID
	GetOrAddChild(name string) (Inode, error)
	Attributes() fuseops.InodeAttributes
	ListChildren() (children []*fuseutil.Dirent, err error)
	Contents() (string, error)
}
