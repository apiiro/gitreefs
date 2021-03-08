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
	GetOrAddChild(name string) (Inode, error)
	Attributes() fuseops.InodeAttributes
	ListChildren() ([]fuseutil.Dirent, error)
	Contents() (string, error)
}

type Inode struct {
	InodeInterface
	Id      fuseops.InodeID
	OwnerId fuseops.InodeID
	Name    string
}

func (in *Inode) Attributes() fuseops.InodeAttributes {
	return fs.DirAttributes()
}

func (in *Inode) ListChildren() ([]fuseutil.Dirent, error) {
	return []fuseutil.Dirent{}, nil
}


func (in *Inode) Contents() (string, error) {
	return "", nil
}

func (in *Inode) String() string {
	return fmt.Sprintf("%v:{InodeID=%v,Name='%v'}", reflect.TypeOf(in), in.Id, in.Name)
}
