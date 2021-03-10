package fs

import (
	"github.com/jacobsa/fuse/fuseutil"
	"gitreefs/git"
	"path"
	"sync"
)

type CommitishInode struct {
	Inode
	commitish  string
	repository *RepositoryInode
	isFetched  bool
	rootEntry  *EntryInode
	mutex      *sync.Mutex
}

func NewCommitishInode(parent *RepositoryInode, commitish string) (inode *CommitishInode, err error) {
	inode = &CommitishInode{
		Inode: Inode{
			Id:      NextInodeID(),
			OwnerId: parent.Id,
		},
		commitish:  commitish,
		repository: parent,
		isFetched:  false,
		mutex:      &sync.Mutex{},
	}
	return
}

func (in *CommitishInode) inodeTreeFromGitTree(gitEntry *git.Entry, entryPath string) (entry *EntryInode, err error) {
	var entries []*EntryInode = nil
	if gitEntry.IsDir {
		entries = make([]*EntryInode, len(gitEntry.EntriesByName))
		i := 0
		for name, childGitEntry := range gitEntry.EntriesByName {
			var childEntry *EntryInode
			childEntry, err = in.inodeTreeFromGitTree(childGitEntry, path.Join(entryPath, name))
			if err != nil {
				return nil, err
			}
			entries[i] = childEntry
			i++
		}
	}
	entry, err = NewEntryInode(
		in,
		entryPath,
		gitEntry,
		entries,
	)
	return
}

func (in *CommitishInode) ListChildren(buffer []byte, offset int) ([]fuseutil.Dirent, error) {
	if !in.isFetched {
		in.mutex.Lock()
		if !in.isFetched {
			root, err := in.repository.provider.ListTree(in.commitish)
			if err != nil {
				return nil, err
			}
			in.rootEntry, err = in.inodeTreeFromGitTree(&root.Entry, "")
			if err != nil {
				return nil, err
			}
			in.isFetched = true
		}
		in.mutex.Unlock()
	}
	return in.rootEntry.ListChildren(buffer, offset)
}
