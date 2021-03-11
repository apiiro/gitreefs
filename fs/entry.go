package fs

import (
	"fmt"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"gitreefs/git"
	"sort"
)

type EntryInode struct {
	Inode
	id fuseops.InodeID
	size          int64
	isDir         bool
	entriesByName map[string]*EntryInode
	path          string
	commitish     *CommitishInode
}

func NewEntryInode(commitish *CommitishInode, path string, gitEntry *git.Entry, entries []*EntryInode) (inode *EntryInode, err error) {
	entriesByName := make(map[string]*EntryInode)
	for _, entry := range entries {
		entriesByName[git.ExtractBaseName(entry.path)] = entry
	}
	if gitEntry.IsDir {
		// directory paths aren't used hence aren't saved
		path = ""
	}
	return &EntryInode{
		id: NextInodeID(),
		commitish:     commitish,
		size:          gitEntry.Size,
		isDir:         gitEntry.IsDir,
		entriesByName: entriesByName,
		path:          path,
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
	childEntry, found := in.entriesByName[name]
	if !found {
		return nil, fmt.Errorf("no child with name %v for %v", name, in)
	}
	return childEntry, err
}

func (in *EntryInode) ListChildren() (children []*fuseutil.Dirent, err error) {
	if !in.isDir {
		return in.Inode.ListChildren()
	}
	names := make([]string, len(in.entriesByName))
	i := 0
	for name, _ := range in.entriesByName {
		names[i] = name
		i++
	}
	sort.Strings(names)
	children = make([]*fuseutil.Dirent, len(in.entriesByName))
	for i := 0; i < len(names); i++ {
		name := names[i]
		childEntry := in.entriesByName[name]
		var childType fuseutil.DirentType
		if childEntry.isDir {
			childType = fuseutil.DT_Directory
		} else {
			childType = fuseutil.DT_File
		}
		children[i] = &fuseutil.Dirent{
			Offset: fuseops.DirOffset(i),
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
