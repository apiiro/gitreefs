package fs

import (
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	cmap "github.com/orcaman/concurrent-map"
	"sync"
)

type RootInode struct {
	clonesPath         string
	repositoriesByName cmap.ConcurrentMap
	mutex              *sync.Mutex
}

var _ Inode = &RootInode{}

func NewRootInode(clonesPath string) (root *RootInode, err error) {
	return &RootInode{
		clonesPath:         clonesPath,
		repositoriesByName: cmap.New(),
		mutex:              &sync.Mutex{},
	}, nil
}

func (in *RootInode) GetOrAddChild(name string) (child Inode, err error) {
	wrapped :=
		in.repositoriesByName.Upsert(name, nil, func(found bool, existingValue interface{}, _ interface{}) interface{} {
			if found && existingValue != nil && existingValue.(*RepositoryInode) != nil {
				return existingValue
			}
			var repository *RepositoryInode
			repository, err = NewRepositoryInode(in.clonesPath, name)
			return repository
		})
	if wrapped.(*RepositoryInode) == nil {
		return nil, err
	}
	return wrapped.(*RepositoryInode), err
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
