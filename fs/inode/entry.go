package fs

import (
	"fmt"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"gitreefs/fs"
	"gitreefs/git"
	"sort"
)

type EntryInode struct {
	Inode
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
		Inode: Inode{
			Id:      NextInodeID(),
			OwnerId: commitish.Id,
		},
		commitish:     commitish,
		size:          gitEntry.Size,
		isDir:         gitEntry.IsDir,
		entriesByName: entriesByName,
		path:          path,
	}, nil
}
func (in *EntryInode) Attributes() fuseops.InodeAttributes {
	if in.isDir {
		return fs.DirAttributes()
	}
	return fs.FileAttributes(in.size)
}

func (in *EntryInode) ListChildren(buffer []byte, offset int) (dirents []fuseutil.Dirent, err error) {
	if !in.isDir {
		return in.Inode.ListChildren(buffer, offset)
	}
	if offset >= len(in.entriesByName) {
		return nil, fmt.Errorf("offset is too big: %v >= %v", offset, len(in.entriesByName))
	}
	names := make([]string, len(in.entriesByName)-offset)
	i := 0
	for name, _ := range in.entriesByName {
		names[i] = name
		i++
	}
	sort.Strings(names)
	dirents = make([]fuseutil.Dirent, len(in.entriesByName))
	for i := offset; i < len(names); i++ {
		name := names[i]
		childEntry := in.entriesByName[name]
		var childType fuseutil.DirentType
		if childEntry.isDir {
			childType = fuseutil.DT_Directory
		} else {
			childType = fuseutil.DT_File
		}
		fuseutil.WriteDirent(buffer, fuseutil.Dirent{
			Offset: fuseops.DirOffset(i),
			Inode:  childEntry.Id,
			Name:   name,
			Type:   childType,
		})
	}
	return
}

func (in *EntryInode) Contents() (string, error) {
	if in.isDir {
		return in.Inode.Contents()
	}
	return in.commitish.repository.provider.FileContents(in.commitish.commitish, in.path)
}
