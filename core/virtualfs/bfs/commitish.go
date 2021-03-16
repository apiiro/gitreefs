package bfs

import (
	"gitreefs/core/git"
	"gitreefs/core/logger"
	"os"
	"sync"
)

type Commitish struct {
	name      string
	provider  *git.RepositoryProvider
	rootEntry *git.RootEntry
	mutex     *sync.Mutex
}

func NewCommitish(name string, provider *git.RepositoryProvider) (commitish *Commitish, err error) {
	logger.Debug("NewCommitish: %v", name)
	return &Commitish{
		name:      name,
		provider:  provider,
		rootEntry: nil,
		mutex:     &sync.Mutex{},
	}, err
}

func (commitish *Commitish) fetchContentIfNeeded() (err error) {
	commitish.mutex.Lock()
	defer commitish.mutex.Unlock()
	if commitish.rootEntry != nil {
		return
	}
	commitish.rootEntry, err = commitish.provider.ListTree(commitish.name)
	return
}

func (commitish *Commitish) GetEntry(subPath string) (entry *git.Entry, err error) {
	err = commitish.fetchContentIfNeeded()
	if err != nil {
		return
	}
	entry, _ = commitish.rootEntry.EntriesByPath[subPath]
	return
}

func (commitish *Commitish) ListDir(subPath string) ([]os.FileInfo, error) {
	entry, err := commitish.GetEntry(subPath)
	if err != nil || entry == nil {
		return nil, err
	}
	children := make([]os.FileInfo, len(entry.EntriesByName))
	i := 0
	for name, child := range entry.EntriesByName {
		if child.IsDir {
			children[i], err = statDir(name)
		} else {
			children[i], err = statFile(name, child.Size)
		}
		if err != nil {
			return nil, err
		}
		i++
	}
	return children, nil
}

func (commitish *Commitish) FileContents(subPath string) (string, error) {
	entry, err := commitish.GetEntry(subPath)
	if err != nil || entry == nil {
		return "", err
	}
	return commitish.provider.FileContents(commitish.name, subPath)
}
