package fs

import (
	"fmt"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"gitreefs/git"
	"sync"
)

type EntryInode struct {
	id               fuseops.InodeID
	size             int64
	isDir            bool
	entries          []*EntryInode
	entryNameToIndex *sync.Map
	path             string
	commitish        *CommitishInode
}

var _ Inode = &EntryInode{}

func NewEntryInode(
	commitish *CommitishInode,
	path string,
	gitEntry *git.Entry,
	entries []*EntryInode,
	entryNameToIndex *sync.Map,
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
	childIndex, found := in.entryNameToIndex.Load(name)
	if !found {
		return nil, fmt.Errorf("no child with name %v for %v", name, in)
	}
	return in.entries[childIndex.(int)], err
}

func (in *EntryInode) ListChildren() (children []*fuseutil.Dirent, err error) {
	if !in.isDir {
		return []*fuseutil.Dirent{}, nil
	}
	children = make([]*fuseutil.Dirent, len(in.entries))
	in.entryNameToIndex.Range(func(name, i interface{}) bool {
		index := i.(int)
		childEntry := in.entries[index]
		var childType fuseutil.DirentType
		if childEntry.isDir {
			childType = fuseutil.DT_Directory
		} else {
			childType = fuseutil.DT_File
		}
		children[index] = &fuseutil.Dirent{
			Offset: fuseops.DirOffset(index + 1),
			Inode:  childEntry.id,
			Name:   name.(string),
			Type:   childType,
		}
		return true
	})
	return
}

func (in *EntryInode) Contents() (string, error) {
	if in.isDir {
		return "", nil
	}
	return in.commitish.repository.provider.FileContents(in.commitish.commitish, in.path)
}
