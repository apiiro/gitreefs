package fs

import (
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"gitreefs/git"
	"gitreefs/logger"
	"path"
	"sync"
)

type CommitishInode struct {
	id         fuseops.InodeID
	commitish  string
	repository *RepositoryInode
	isFetched  bool
	rootEntry  *EntryInode
	mutex      *sync.Mutex
}

var _ Inode = &CommitishInode{}

func NewCommitishInode(parent *RepositoryInode, commitish string) (inode *CommitishInode, err error) {
	var canResolve bool
	canResolve, err = parent.provider.CanResolve(commitish)
	if err != nil || !canResolve {
		return nil, err
	}
	inode = &CommitishInode{
		id:         NextInodeID(),
		commitish:  commitish,
		repository: parent,
		isFetched:  false,
		mutex:      &sync.Mutex{},
	}
	logger.Debug("NewCommitishInode: %v :: %v", commitish, parent.clonePath)
	return
}

func (in *CommitishInode) Id() fuseops.InodeID {
	return in.id
}

func (in *CommitishInode) inodeTreeFromGitTree(gitEntry *git.Entry, entryPath string) (entry *EntryInode, err error) {
	var entries []*EntryInode = nil
	var entryNameToIndex *sync.Map = nil
	if gitEntry.IsDir {
		entries = make([]*EntryInode, len(gitEntry.EntriesByName))
		entryNameToIndex = &sync.Map{}
		i := 0
		for name, childGitEntry := range gitEntry.EntriesByName {
			var childEntry *EntryInode
			childEntry, err = in.inodeTreeFromGitTree(childGitEntry, path.Join(entryPath, name))
			if err != nil {
				return nil, err
			}
			entries[i] = childEntry
			entryNameToIndex.Store(name, i)
			i++
		}
	}
	entry, err = NewEntryInode(
		in,
		entryPath,
		gitEntry,
		entries,
		entryNameToIndex,
	)
	return
}

func (in *CommitishInode) fetchContentIfNeeded() (err error) {
	in.mutex.Lock()
	if !in.isFetched {
		var root *git.RootEntry
		root, err = in.repository.provider.ListTree(in.commitish)
		if err == nil {
			in.rootEntry, err = in.inodeTreeFromGitTree(&root.Entry, "")
			if err == nil {
				in.isFetched = true
			}
		}
	}
	in.mutex.Unlock()
	return
}

func (in *CommitishInode) GetOrAddChild(name string) (child Inode, err error) {
	err = in.fetchContentIfNeeded()
	if err != nil {
		return nil, err
	}
	return in.rootEntry.GetOrAddChild(name)
}

func (in *CommitishInode) ListChildren() (_ []*fuseutil.Dirent, err error) {
	err = in.fetchContentIfNeeded()
	if err != nil {
		return nil, err
	}
	return in.rootEntry.ListChildren()
}

func (in *CommitishInode) Attributes() fuseops.InodeAttributes {
	// default implementation
	return DirAttributes()
}

func (in *CommitishInode) Contents() (string, error) {
	// default implementation
	return "", nil
}
