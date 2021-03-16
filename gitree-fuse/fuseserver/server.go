package fuseserver

import (
	"github.com/jacobsa/fuse"
	"github.com/jacobsa/fuse/fuseops"
	"github.com/jacobsa/fuse/fuseutil"
	"gitreefs/logger"
	"gitreefs/virtualfs/inodefs"
	"golang.org/x/net/context"
	"sync"
)

type fuseFs struct {
	fuseutil.NotImplementedFileSystem
	clonesPath string
	inodes     *sync.Map
}

func NewFsServer(clonesPath string) (server fuse.Server, err error) {
	var rootInode *inodefs.RootInode
	rootInode, err = inodefs.NewRootInode(clonesPath)
	if err != nil {
		return
	}
	inodes := &sync.Map{}
	inodes.Store(rootInode.Id(), rootInode)
	server = fuseutil.NewFileSystemServer(&fuseFs{
		clonesPath: clonesPath,
		inodes:     inodes,
	})
	return
}

func (fs *fuseFs) StatFS(
	ctx context.Context,
	op *fuseops.StatFSOp) error {
	return nil
}

func (fs *fuseFs) lookUpInode(parentId fuseops.InodeID, name string) (inode inodefs.Inode, err error) {
	parent, found := fs.inodes.Load(parentId)
	if !found {
		return nil, nil
	}
	inode, err = parent.(inodefs.Inode).GetOrAddChild(name)
	if err != nil || inode == nil {
		return
	}
	fs.inodes.LoadOrStore(inode.Id(), inode)
	return
}

func (fs *fuseFs) LookUpInode(
	ctx context.Context,
	op *fuseops.LookUpInodeOp) error {
	inode, err := fs.lookUpInode(op.Parent, op.Name)
	if err != nil {
		logger.Error("fuseFs.LookUpInode for %v on %v: %v", inode, op.Name, err)
		return fuse.EIO
	}
	if inode == nil {
		return fuse.ENOENT
	}
	outputEntry := &op.Entry
	outputEntry.Child = inode.Id()
	outputEntry.Attributes = inode.Attributes()
	return nil
}

func (fs *fuseFs) GetInodeAttributes(
	ctx context.Context,
	op *fuseops.GetInodeAttributesOp) error {
	var inode, found = fs.inodes.Load(op.Inode)
	if !found {
		return fuse.ENOENT
	}
	op.Attributes = inode.(inodefs.Inode).Attributes()
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
	var inode, found = fs.inodes.Load(op.Inode)
	if !found {
		return fuse.ENOENT
	}
	children, err := inode.(inodefs.Inode).ListChildren()
	if err != nil {
		logger.Error("fuseFs.ReadDir for %v: %v", inode, err)
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
	var inode, found = fs.inodes.Load(op.Inode)
	if !found {
		return fuse.ENOENT
	}
	contents, err := inode.(inodefs.Inode).Contents()
	if err != nil {
		logger.Error("fuseFs.ReadFile for %v: %v", inode, err)
		return fuse.EIO
	}

	if op.Offset > int64(len(contents)) {
		return fuse.EIO
	}

	contents = contents[op.Offset:]
	op.BytesRead = copy(op.Dst, contents)
	return nil
}

func (fs *fuseFs) ReleaseDirHandle(
	ctx context.Context,
	op *fuseops.ReleaseDirHandleOp) error {
	return nil
}

func (fs *fuseFs) GetXattr(
	ctx context.Context,
	op *fuseops.GetXattrOp) error {
	return nil
}

func (fs *fuseFs) ListXattr(
	ctx context.Context,
	op *fuseops.ListXattrOp) error {
	return fuse.ENOSYS
}

func (fs *fuseFs) ForgetInode(
	ctx context.Context,
	op *fuseops.ForgetInodeOp) error {
	return nil
}

func (fs *fuseFs) ReleaseFileHandle(
	ctx context.Context,
	op *fuseops.ReleaseFileHandleOp) error {
	return nil
}

func (fs *fuseFs) FlushFile(
	ctx context.Context,
	op *fuseops.FlushFileOp) error {
	return nil
}

