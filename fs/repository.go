package fs

import (
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"gitreefs/git"
	"gitreefs/logger"
	"path"
	"sync"
)

type RepositoryInode struct {
	Inode
	id              fuseops.InodeID
	clonePath       string
	provider        *git.RepositoryProvider
	commitishByName map[string]*CommitishInode
	mutex           *sync.Mutex
}

func NewRepositoryInode(clonesPath string, name string) (inode *RepositoryInode, err error) {
	clonePath := path.Join(clonesPath, name)
	err = ValidateDirectory(clonePath, false)
	if err != nil {
		return
	}
	var provider *git.RepositoryProvider
	provider, err = git.NewRepositoryProvider(clonePath)
	if err != nil {
		return
	}
	inode = &RepositoryInode{
		id:              NextInodeID(),
		provider:        provider,
		clonePath:       clonePath,
		commitishByName: map[string]*CommitishInode{},
		mutex:           &sync.Mutex{},
	}
	logger.Debug("NewRepositoryInode: %v", inode.clonePath)
	return
}

func (in *RepositoryInode) Id() fuseops.InodeID {
	return in.id
}

func (in *RepositoryInode) GetOrAddChild(name string) (child Inode, err error) {
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
	child = commitish
	return
}

func (in *RepositoryInode) ListChildren() ([]*fuseutil.Dirent, error) {
	// ListChildren isn't implemented as there is no use case to list all possible commitishes
	return []*fuseutil.Dirent{}, nil
}

func (in *RepositoryInode) Attributes() fuseops.InodeAttributes {
	// default implementation
	return DirAttributes()
}

func (in *RepositoryInode) Contents() (string, error) {
	// default implementation
	return "", nil
}
