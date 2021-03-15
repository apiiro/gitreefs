package fs

import (
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"github.com/orcaman/concurrent-map"
	"gitreefs/git"
	"gitreefs/logger"
	"gitreefs/common"
	"path"
)

type RepositoryInode struct {
	id              fuseops.InodeID
	clonePath       string
	provider        *git.RepositoryProvider
	commitishByName cmap.ConcurrentMap
}

var _ Inode = &RepositoryInode{}

func NewRepositoryInode(clonesPath string, name string) (inode *RepositoryInode, err error) {
	clonePath := path.Join(clonesPath, name)
	err = common.ValidateDirectory(clonePath, false)
	if err != nil {
		return
	}
	var provider *git.RepositoryProvider
	provider, err = git.NewRepositoryProvider(clonePath)
	if err != nil || provider == nil {
		return nil, err
	}
	inode = &RepositoryInode{
		id:              NextInodeID(),
		provider:        provider,
		clonePath:       clonePath,
		commitishByName: cmap.New(),
	}
	logger.Debug("NewRepositoryInode: %v", inode.clonePath)
	return
}

func (in *RepositoryInode) Id() fuseops.InodeID {
	return in.id
}

func (in *RepositoryInode) GetOrAddChild(name string) (child Inode, err error) {
	wrapped :=
		in.commitishByName.Upsert(name, nil, func(found bool, existingValue interface{}, _ interface{}) interface{} {
			if found && existingValue != nil && existingValue.(*CommitishInode) != nil {
				return existingValue
			}
			var commitish *CommitishInode
			commitish, err = NewCommitishInode(in, name)
			return commitish
		})
	if wrapped.(*CommitishInode) == nil {
		return nil, err
	}
	return wrapped.(*CommitishInode), err
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
