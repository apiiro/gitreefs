package fs

import (
	"github.com/jacobsa/fuse/fuseops"
)

type Commitish struct {
	Inode
	repositoryId fuseops.InodeID
	isFetched    bool
	rootId       fuseops.InodeID
}
