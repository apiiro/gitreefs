package fs

import (
	"github.com/go-git/go-billy/v5"
	"os"
)

const (
	ERROR_MESSAGE = "operation is not readonly and not supported"
)

func (fs *GitFileSystem) Create(filename string) (billy.File, error) {
	panic(ERROR_MESSAGE)
}
func (fs *GitFileSystem) Rename(oldpath, newpath string) error {
	panic(ERROR_MESSAGE)
}

func (fs *GitFileSystem) Remove(filename string) error {
	panic(ERROR_MESSAGE)
}

func (fs *GitFileSystem) TempFile(dir, prefix string) (billy.File, error) {
	panic(ERROR_MESSAGE)
}

func (fs *GitFileSystem) MkdirAll(filename string, perm os.FileMode) error {
	panic(ERROR_MESSAGE)
}

func (fs *GitFileSystem) Lstat(filename string) (os.FileInfo, error) {
	panic(ERROR_MESSAGE)
}

func (fs *GitFileSystem) Symlink(target, link string) error {
	panic(ERROR_MESSAGE)
}

func (fs *GitFileSystem) Readlink(link string) (string, error) {
	panic(ERROR_MESSAGE)
}

func (fs *GitFileSystem) Chroot(path string) (billy.Filesystem, error) {
	panic(ERROR_MESSAGE)
}

func (file *File) Write(p []byte) (n int, err error) {
	panic(ERROR_MESSAGE)
}
func (file *File) Lock() error {
	panic(ERROR_MESSAGE)
}

func (file *File) Unlock() error {
	panic(ERROR_MESSAGE)
}

func (file *File) Truncate(size int64) error {
	panic(ERROR_MESSAGE)
}
