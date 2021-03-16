package bfs

import (
	"github.com/go-git/go-billy/v5"
	"gitreefs/core/git"
	"gitreefs/core/logger"
	"os"
	"path/filepath"
)

type GitFileSystem struct {
	root *Root
}

var _ billy.Filesystem = &GitFileSystem{}
var _ billy.Capable = &GitFileSystem{}

func NewGitFileSystem(clonesPath string) (*GitFileSystem, error) {
	root, err := NewRoot(clonesPath)
	if err != nil {
		return nil, err
	}
	return &GitFileSystem{
		root: root,
	}, nil
}

func (fs *GitFileSystem) Capabilities() billy.Capability {
	return billy.ReadCapability | billy.SeekCapability
}

func (fs *GitFileSystem) Open(path string) (billy.File, error) {
	components, err := breakdown(path)
	if err != nil {
		logger.Error("fs.Open: could not find '%v': %v", path, err)
		return nil, os.ErrNotExist
	}

	if !components.hasRepository() || !components.hasCommitish() {
		return nil, os.ErrNotExist
	}

	var repository *Repository
	repository, err = fs.root.getOrAddRepository(components.repositoryName)
	if err != nil || repository == nil {
		logger.Info("fs.Open: could not find repository for %v: %v", path, err)
		return nil, os.ErrNotExist
	}

	var commitish *Commitish
	commitish, err = repository.getOrAddCommitish(components.commitishName)
	if err != nil || commitish == nil {
		logger.Info("fs.Open: could not find commitish for %v: %v", path, err)
		return nil, os.ErrNotExist
	}

	var entry *git.Entry
	entry, err = commitish.GetEntry(components.subPath)
	if err != nil || entry == nil {
		logger.Info("fs.Open: could not find file for %v: %v", path, err)
		return nil, os.ErrNotExist
	}
	file, err := NewFile(path, commitish, components.subPath, entry.Size)
	if err != nil || file == nil {
		logger.Info("fs.Open: could not open file for %v: %v", path, err)
		return nil, os.ErrNotExist
	}
	return file, nil
}

func (fs *GitFileSystem) OpenFile(filename string, _ int, _ os.FileMode) (billy.File, error) {
	return fs.Open(filename)
}

func (fs *GitFileSystem) Stat(path string) (os.FileInfo, error) {
	components, err := breakdown(path)
	if err != nil {
		logger.Error("fs.Stat: could not find '%v': %v", path, err)
		return nil, os.ErrNotExist
	}

	if !components.hasRepository() {
		return statDir("")
	}

	var repository *Repository
	repository, err = fs.root.getOrAddRepository(components.repositoryName)
	if err != nil || repository == nil {
		logger.Info("fs.Stat: could not find repository for %v: %v", path, err)
		return nil, os.ErrNotExist
	}

	if !components.hasCommitish() {
		return statDir(repository.name)
	}

	var commitish *Commitish
	commitish, err = repository.getOrAddCommitish(components.commitishName)
	if err != nil || commitish == nil {
		logger.Info("fs.Stat: could not find commitish for %v: %v", path, err)
		return nil, os.ErrNotExist
	}

	entry, err := commitish.GetEntry(components.subPath)
	if err != nil || entry == nil {
		logger.Info("fs.Stat: could not find git entry for %v: %v", path, err)
		return nil, os.ErrNotExist
	}

	baseName := git.ExtractBaseName(path)
	if entry.IsDir {
		info, err := statDir(baseName)
		if info == nil || err != nil {
			logger.Error("fs.Stat: failed to stat %v: %v", path, err)
			return nil, os.ErrNotExist
		}
		return info, nil
	}
	info, err := statFile(baseName, entry.Size)
	if info == nil || err != nil {
		logger.Error("fs.Stat: failed to stat %v: %v", path, err)
		return nil, os.ErrNotExist
	}
	return info, nil
}

func (fs *GitFileSystem) ReadDir(path string) ([]os.FileInfo, error) {
	components, err := breakdown(path)
	if err != nil {
		logger.Error("fs.ReadDir: could not find '%v': %v", path, err)
		return nil, os.ErrNotExist
	}

	if !components.hasRepository() || !components.hasCommitish() {
		return []os.FileInfo{}, nil
	}

	var repository *Repository
	repository, err = fs.root.getOrAddRepository(components.repositoryName)
	if err != nil || repository == nil {
		logger.Info("fs.ReadDir: could not find repository for %v: %v", path, err)
		return nil, os.ErrNotExist
	}

	var commitish *Commitish
	commitish, err = repository.getOrAddCommitish(components.commitishName)
	if err != nil || commitish == nil {
		logger.Info("fs.ReadDir: could not find commitish for %v: %v", path, err)
		return nil, os.ErrNotExist
	}

	files, err := commitish.ListDir(components.subPath)
	if files == nil || err != nil {
		logger.Error("fs.ReadDir: failed on %v: %v", path, err)
		return nil, os.ErrNotExist
	}
	return files, nil
}

func (fs *GitFileSystem) Join(elem ...string) string {
	return filepath.Join(elem...)
}

func (fs *GitFileSystem) Root() string {
	return "/"
}
