package fs

import (
	"fmt"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"gitreefs/git"
)

type EntryInode struct {
	Inode
	id               fuseops.InodeID
	size             int64
	isDir            bool
	entries          []*EntryInode
	entryNameToIndex map[string]int
	path             string
	commitish        *CommitishInode
}

func NewEntryInode(
	commitish *CommitishInode,
	path string,
	gitEntry *git.Entry,
	entries []*EntryInode,
	entryNameToIndex map[string]int,
) (inode *EntryInode, err error) {
	if gitEntry.IsDir {
		// directory paths aren't used hence aren't saved
		path = ""
	}
	return &EntryInode{
		id:               NextInodeID(),
		commitish:        commitish,
		size:             gitEntry.Size,
		isDir:            gitEntry.IsDir,
		entries:          entries,
		entryNameToIndex: entryNameToIndex,
		path:             path,
	}, nil
}

func (in *EntryInode) Id() fuseops.InodeID {
	return in.id
}

func (in *EntryInode) Attributes() fuseops.InodeAttributes {
	if in.isDir {
		return DirAttributes()
	}
	return FileAttributes(in.size)
}

func (in *EntryInode) GetOrAddChild(name string) (child Inode, err error) {
	childIndex, found := in.entryNameToIndex[name]
	if !found {
		return nil, fmt.Errorf("no child with name %v for %v", name, in)
	}
	return in.entries[childIndex], err
}

func (in *EntryInode) ListChildren() (children []*fuseutil.Dirent, err error) {
	if !in.isDir {
		return in.Inode.ListChildren()
	}
	children = make([]*fuseutil.Dirent, len(in.entries))
	for name, i := range in.entryNameToIndex {
		childEntry := in.entries[i]
		var childType fuseutil.DirentType
		if childEntry.isDir {
			childType = fuseutil.DT_Directory
		} else {
			childType = fuseutil.DT_File
		}
		children[i] = &fuseutil.Dirent{
			Offset: fuseops.DirOffset(i + 1),
			Inode:  childEntry.id,
			Name:   name,
			Type:   childType,
		}
	}
	return
}

func (in *EntryInode) Contents() (string, error) {
	if in.isDir {
		return in.Inode.Contents()
	}
	return in.commitish.repository.provider.FileContents(in.commitish.commitish, in.path)
}
