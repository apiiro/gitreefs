package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gitreefs/logger"
)

const (
	MaxFileSizeMB    int64 = 6
	MaxFileSizeBytes       = MaxFileSizeMB * 1024 * 1024
	ShortShaLength         = 7
	RootEntryPath          = ""
)

type RepositoryProvider struct {
	repository      *git.Repository
	shortShaMapping map[string]string
}

func NewRepositoryProvider(clonePath string) (provider *RepositoryProvider, err error) {
	provider = &RepositoryProvider{
		shortShaMapping: make(map[string]string),
	}
	provider.repository, err = git.PlainOpen(clonePath)
	if err != nil {
		return
	}

	// Manual implementation of short sha mapping, due to bug in go-git: https://github.com/go-git/go-git/issues/148
	var iter object.CommitIter
	iter, err = provider.repository.CommitObjects()
	if err != nil {
		return
	}
	err = iter.ForEach(func(commit *object.Commit) error {
		sha := commit.Hash.String()
		shortSha := sha[:ShortShaLength]
		provider.shortShaMapping[shortSha] = sha
		return nil
	})

	logger.Info("NewRepositoryProvider for %v with total of %v commits detected", clonePath, len(provider.shortShaMapping))
	return
}

func (provider *RepositoryProvider) getCommit(commitish string) (commit *object.Commit, err error) {

	if len(commitish) == ShortShaLength {
		sha, found := provider.shortShaMapping[commitish]
		if found {
			commitish = sha
		}
	}

	var hash *plumbing.Hash
	hash, err = provider.repository.ResolveRevision(plumbing.Revision(commitish))
	if err != nil {
		return
	}

	commit, err = provider.repository.CommitObject(*hash)
	return
}

func (provider *RepositoryProvider) ListTree(commitish string) (root *RootEntry, err error) {

	var commit *object.Commit
	commit, err = provider.getCommit(commitish)
	if err != nil {
		return
	}

	var tree *object.Tree
	tree, err = commit.Tree()
	if err != nil {
		return
	}

	root = &RootEntry{
		Entry: Entry{
			Size:          0,
			IsDir:         true,
			EntriesByName: make(map[string]*Entry),
		},
		EntriesByPath: make(map[string]*Entry),
	}
	root.EntriesByPath[RootEntryPath] = &root.Entry

	err = tree.Files().ForEach(func(file *object.File) (err error) {
		mode := file.Mode
		if !mode.IsFile() || mode.IsMalformed() || !mode.IsRegular() {
			return
		}
		filePath := file.Name
		fileName := ExtractBaseName(filePath)
		parentPath := ExtractDirPath(filePath)
		parent := root.ensurePath(parentPath)
		fileEntry := FileEntry(parent, fileName, file.Size)
		root.EntriesByPath[filePath] = fileEntry
		return
	})

	logger.Info("ListTree for %v with total of %v paths detected", commitish, len(root.EntriesByPath))

	return
}

func (provider *RepositoryProvider) FileContents(commitish string, filePath string) (contents string, err error) {
	var commit *object.Commit
	commit, err = provider.getCommit(commitish)
	if err != nil {
		return
	}

	var file *object.File
	file, err = commit.File(filePath)
	if err != nil {
		return
	}

	if file.Size >= MaxFileSizeBytes {
		err = fmt.Errorf("file size is too large to load to memory - %v at %v/%v", file.Size, commitish, filePath)
		return
	}

	contents, err = file.Contents()
	if err != nil {
		logger.Info("FileContents for %v :: %v with content of size %v", commitish, filePath, len(contents))
	}
	return
}
