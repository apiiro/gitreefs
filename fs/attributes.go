package fs

import (
	"github.com/jacobsa/fuse/fuseops"
	"os"
	"time"
)

var (
	uid uint32 = uint32(os.Getuid())
	gid uint32 = uint32(os.Getgid())
)

func FileAttributes(size int64) fuseops.InodeAttributes {
	return fuseops.InodeAttributes{
		Size:  uint64(size),
		Nlink: 1,
		Mode:  os.ModePerm,
		Atime: time.Now(),
		Mtime: time.Now(),
		Ctime: time.Now(),
		Uid:   uid,
		Gid:   gid,
	}
}

func DirAttributes() fuseops.InodeAttributes {
	return fuseops.InodeAttributes{
		Size:  0,
		Nlink: 1,
		Mode:  os.ModeDir | os.ModePerm,
		Atime: time.Now(),
		Mtime: time.Now(),
		Ctime: time.Now(),
		Uid:   uid,
		Gid:   gid,
	}
}
