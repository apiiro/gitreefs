package fs

import (
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"sync"
)

type RootInode struct {
	Inode
	clonesPath         string
	repositoriesByName map[string]*RepositoryInode
	mutex              *sync.Mutex
}

func NewRootInode(clonesPath string) (root *RootInode, err error) {
	return &RootInode{
		clonesPath:         clonesPath,
		repositoriesByName: make(map[string]*RepositoryInode),
		mutex:              &sync.Mutex{},
	}, nil
}

func (in *RootInode) GetOrAddChild(name string) (child Inode, err error) {
	repository, found := in.repositoriesByName[name]
	if !found {
		in.mutex.Lock()
		repository, found = in.repositoriesByName[name]
		if !found {
			repository, err = NewRepositoryInode(in.clonesPath, name)
			if err != nil {
				in.mutex.Unlock()
				return
			}
			in.repositoriesByName[name] = repository
		}
		in.mutex.Unlock()
	}
	child = repository
	return
}

func (in *RootInode) Id() fuseops.InodeID {
	return fuseops.RootInodeID
}

func (in *RootInode) ListChildren() ([]*fuseutil.Dirent, error) {
	// ListChildren isn't implemented as there is no use case to list all possible commitishes
	return []*fuseutil.Dirent{}, nil
}

func (in *RootInode) Attributes() fuseops.InodeAttributes {
	// default implementation
	return DirAttributes()
}

func (in *RootInode) Contents() (string, error) {
	// default implementation
	return "", nil
}
