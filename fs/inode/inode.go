package fs

import (
	"fmt"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"gitreefs/fs"
	"reflect"
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

type InodeInterface interface {
	GetOrAddChild(name string) (*Inode, error)
	Attributes() fuseops.InodeAttributes
	ListChildren() (children []*fuseutil.Dirent, err error)
	Contents() (string, error)
}

type Inode struct {
	InodeInterface
	Id      fuseops.InodeID
	OwnerId fuseops.InodeID
}

func (in *Inode) Attributes() fuseops.InodeAttributes {
	// default implementation
	return fs.DirAttributes()
}

func (in *Inode) ListChildren() ([]*fuseutil.Dirent, error) {
	// default implementation
	return []*fuseutil.Dirent{}, nil
}

func (in *Inode) Contents() (string, error) {
	// default implementation
	return "", nil
}

func (in *Inode) String() string {
	return fmt.Sprintf("%v:{InodeID=%v}", reflect.TypeOf(in), in.Id)
}
