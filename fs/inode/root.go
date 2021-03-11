package fs

import (
	"github.com/jacobsa/fuse/fuseops"
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
		Inode: Inode{
			Id:      fuseops.RootInodeID,
			OwnerId: fuseops.RootInodeID,
		},
		clonesPath:         clonesPath,
		repositoriesByName: make(map[string]*RepositoryInode),
		mutex:              &sync.Mutex{},
	}, nil
}

func (in *RootInode) GetOrAddChild(name string) (child *Inode, err error) {
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
	child = &repository.Inode
	return
}
