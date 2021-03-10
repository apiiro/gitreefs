package fs

import (
	"github.com/jacobsa/fuse/fuseops"
	"gitreefs/fs"
	"gitreefs/git"
	"path"
	"sync"
)

type RepositoryInode struct {
	Inode
	clonePath string
	provider  *git.RepositoryProvider
	children  map[string]*CommitishInode
	mutex     *sync.Mutex
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
		provider:  provider,
		clonePath: clonePath,
		children:  map[string]*CommitishInode{},
		mutex:     &sync.Mutex{},
	}
	return
}

func (in *RepositoryInode) GetOrAddChild(name string) (child *Inode, err error) {
	var found bool
	var commitish *CommitishInode
	if !found {
		in.mutex.Lock()
		commitish, found = in.children[name]
		if !found {
			commitish, err = NewCommitishInode(in, name)
			if err != nil {
				return
			}
			in.children[name] = commitish
		}
		in.mutex.Unlock()
	}
	child = &commitish.Inode
	return
}

// ListChildren isn't implemented as there is no use case to list all possible commitishes
