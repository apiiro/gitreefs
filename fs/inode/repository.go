package fs

import (
	"github.com/jacobsa/fuse/fuseops"
	"gitreefs/fs"
	"gitreefs/git"
	"gitreefs/logger"
	"path"
	"sync"
)

type RepositoryInode struct {
	Inode
	clonePath       string
	provider        *git.RepositoryProvider
	commitishByName map[string]*CommitishInode
	mutex           *sync.Mutex
}

func NewRepositoryInode(clonesPath string, name string) (inode *RepositoryInode, err error) {
	clonePath := path.Join(clonesPath, name)
	err = fs.ValidateDirectory(clonePath)
	if err != nil {
		return
	}
	var provider *git.RepositoryProvider
	provider, err = git.NewRepositoryProvider(clonePath)
	if err != nil {
		return
	}
	inode = &RepositoryInode{
		Inode: Inode{
			Id:      NextInodeID(),
			OwnerId: fuseops.RootInodeID,
		},
		provider:        provider,
		clonePath:       clonePath,
		commitishByName: map[string]*CommitishInode{},
		mutex:           &sync.Mutex{},
	}
	logger.Debug("NewRepositoryInode: %v", inode.clonePath)
	return
}

func (in *RepositoryInode) GetOrAddChild(name string) (child *Inode, err error) {
	commitish, found := in.commitishByName[name]
	if !found {
		in.mutex.Lock()
		commitish, found = in.commitishByName[name]
		if !found {
			commitish, err = NewCommitishInode(in, name)
			if err != nil {
				in.mutex.Unlock()
				return
			}
			in.commitishByName[name] = commitish
		}
		in.mutex.Unlock()
	}
	child = &commitish.Inode
	return
}

// ListChildren isn't implemented as there is no use case to list all possible commitishes
