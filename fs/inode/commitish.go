package fs

import (
	"github.com/jacobsa/fuse/fuseutil"
	"gitreefs/git"
	"sync"
)

type CommitishInode struct {
	Inode
	repository *RepositoryInode
	isFetched  bool
	rootEntry  *VersionedEntry
	mutex      *sync.Mutex
}

func NewCommitishInode(parent *RepositoryInode, name string) (inode *CommitishInode, err error) {
	inode = &CommitishInode{
		Inode: Inode{
			Id:      NextInodeID(),
			OwnerId: parent.Id,
			Name:    name,
		},
		repository: parent,
		isFetched:  false,
		mutex:      &sync.Mutex{},
	}
	return
}

func (in *CommitishInode) constructTree(entriesByPath map[string]*git.Entry) {

}

func (in *CommitishInode) ListChildren() ([]fuseutil.Dirent, error) {
	if !in.isFetched {
		in.mutex.Lock()
		if !in.isFetched {
			entriesByPath, err := in.repository.provider.ListTree(in.Name)
			if err != nil {
				return nil, err
			}
			in.constructTree(entriesByPath)
			in.isFetched = true
		}
		in.mutex.Unlock()
	}
	return in.rootEntry.ListChildren()
}
