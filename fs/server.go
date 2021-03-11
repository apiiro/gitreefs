package fs

import (
	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"gitreefs/logger"
	"golang.org/x/net/context"
)

type fuseFs struct {
	fuseutil.NotImplementedFileSystem
	clonesPath string
	inodes     map[fuseops.InodeID]*Inode
}

func NewFsServer(clonesPath string) (server fuse.Server, err error) {
	var rootInode *RootInode
	rootInode, err = NewRootInode(clonesPath)
	if err != nil {
		return
	}
	server = fuseutil.NewFileSystemServer(&fuseFs{
		clonesPath: clonesPath,
		inodes: map[fuseops.InodeID]*Inode{
			rootInode.Id: &rootInode.Inode,
		},
	})
	return
}

func (fs *fuseFs) StatFS(
	ctx context.Context,
	op *fuseops.StatFSOp) error {
	return nil
}

func (fs *fuseFs) lookUpInode(parentId fuseops.InodeID, name string) (inode *Inode, err error) {
	parent, found := fs.inodes[parentId]
	if !found {
		return nil, nil
	}
	inode, err = parent.GetOrAddChild(name)
	if err != nil {
		return
	}
	if _, found = fs.inodes[inode.Id]; !found {
		fs.inodes[inode.Id] = inode
	}
	return
}

func (fs *fuseFs) LookUpInode(
	ctx context.Context,
	op *fuseops.LookUpInodeOp) error {
	inode, err := fs.lookUpInode(op.Parent, op.Name)
	if err != nil {
		logger.Error("fuseFs.LookUpInode for %v on %v: %w", inode, op.Name, err)
		return fuse.EIO
	}
	if inode == nil {
		return fuse.ENOENT
	}
	outputEntry := &op.Entry
	outputEntry.Child = inode.Id
	outputEntry.Attributes = inode.Attributes()
	return nil
}

func (fs *fuseFs) GetInodeAttributes(
	ctx context.Context,
	op *fuseops.GetInodeAttributesOp) error {
	var inode, found = fs.inodes[op.Inode]
	if !found {
		return fuse.ENOENT
	}
	op.Attributes = inode.Attributes()
	return nil
}

func (fs *fuseFs) OpenDir(
	ctx context.Context,
	op *fuseops.OpenDirOp) error {
	// Allow opening any directory.
	return nil
}

func (fs *fuseFs) ReadDir(
	ctx context.Context,
	op *fuseops.ReadDirOp) error {
	var inode, found = fs.inodes[op.Inode]
	if !found {
		return fuse.ENOENT
	}
	children, err := inode.ListChildren()
	if err != nil {
		logger.Error("fuseFs.ReadDir for %v: %w", inode, err)
		return fuse.EIO
	}

	if op.Offset > fuseops.DirOffset(len(children)) {
		return fuse.EIO
	}

	children = children[op.Offset:]

	for _, child := range children {
		bytesWritten := fuseutil.WriteDirent(op.Dst[op.BytesRead:], *child)
		if bytesWritten == 0 {
			break
		}
		op.BytesRead += bytesWritten
	}
	return nil
}

func (fs *fuseFs) OpenFile(
	ctx context.Context,
	op *fuseops.OpenFileOp) error {
	// Allow opening any file.
	return nil
}

func (fs *fuseFs) ReadFile(
	ctx context.Context,
	op *fuseops.ReadFileOp) error {
	var inode, found = fs.inodes[op.Inode]
	if !found {
		return fuse.ENOENT
	}
	contents, err := inode.Contents()
	if err != nil {
		logger.Error("fuseFs.ReadFile for %v: %w", inode, err)
		return fuse.EIO
	}

	if op.Offset > int64(len(contents)) {
		return fuse.EIO
	}

	contents = contents[op.Offset:]
	op.BytesRead = copy(op.Dst, contents[op.Offset:])
	return nil
}
