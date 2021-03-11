package fs

import (
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"sync"
)

var (
	nextInodeIDMutex                 = sync.Mutex{}
	nextInodeID      fuseops.InodeID = fuseops.RootInodeID + 1
)

func NextInodeID() (next fuseops.InodeID) {
	nextInodeIDMutex.Lock()
	next = nextInodeID
	nextInodeID++
	nextInodeIDMutex.Unlock()
	return
}

type Inode interface {
	Id() fuseops.InodeID
	GetOrAddChild(name string) (Inode, error)
	Attributes() fuseops.InodeAttributes
	ListChildren() (children []*fuseutil.Dirent, err error)
	Contents() (string, error)
}
