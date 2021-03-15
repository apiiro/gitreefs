package fs

import (
	"github.com/go-git/go-billy/v5"
	"os"
)

type GitFileSystem struct {
	clonesPath string
}

var _ billy.Filesystem = &GitFileSystem{}

func NewGitFileSystem(clonesPath string) (*GitFileSystem, error) {
	return &GitFileSystem{clonesPath: clonesPath}, nil
}

func (fs *GitFileSystem) Create(filename string) (billy.File, error) {
	panic("implement me")
}

func (fs *GitFileSystem) Open(filename string) (billy.File, error) {
	panic("implement me")
}

func (fs *GitFileSystem) OpenFile(filename string, flag int, perm os.FileMode) (billy.File, error) {
	panic("implement me")
}

func (fs *GitFileSystem) Stat(filename string) (os.FileInfo, error) {
	panic("implement me")
}

func (fs *GitFileSystem) Rename(oldpath, newpath string) error {
	panic("implement me")
}

func (fs *GitFileSystem) Remove(filename string) error {
	panic("implement me")
}

func (fs *GitFileSystem) Join(elem ...string) string {
	panic("implement me")
}

func (fs *GitFileSystem) TempFile(dir, prefix string) (billy.File, error) {
	panic("implement me")
}

func (fs *GitFileSystem) ReadDir(path string) ([]os.FileInfo, error) {
	panic("implement me")
}

func (fs *GitFileSystem) MkdirAll(filename string, perm os.FileMode) error {
	panic("implement me")
}

func (fs *GitFileSystem) Lstat(filename string) (os.FileInfo, error) {
	panic("implement me")
}

func (fs *GitFileSystem) Symlink(target, link string) error {
	panic("implement me")
}

func (fs *GitFileSystem) Readlink(link string) (string, error) {
	panic("implement me")
}

func (fs *GitFileSystem) Chroot(path string) (billy.Filesystem, error) {
	panic("implement me")
}

func (fs *GitFileSystem) Root() string {
	panic("implement me")
}
