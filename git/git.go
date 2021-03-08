package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

const (
	MaxFileSizeMB    int64 = 6
	MaxFileSizeBytes       = MaxFileSizeMB * 1024 * 1024
)

func getCommit(clonePath string, commitish string) (commit *object.Commit, err error) {

	var repository *git.Repository
	repository, err = git.PlainOpen(clonePath)
	if err != nil {
		return
	}

	var hash *plumbing.Hash
	hash, err = repository.ResolveRevision(plumbing.Revision(commitish))
	if err != nil {
		return
	}

	commit, err = repository.CommitObject(*hash)
	return
}

func ensureAncestors(dirPath string, entriesByPath map[string]Entry) {
	if len(dirPath) == 0 || dirPath == "." {
		return
	}
	if _, found := entriesByPath[dirPath]; !found {
		entriesByPath[dirPath] = DirEntry(dirPath)
	}
	ensureAncestors(ExtractDirPath(dirPath), entriesByPath)
}

func ListTree(clonePath string, commitish string) (entriesByPath map[string]Entry, err error) {

	var commit *object.Commit
	commit, err = getCommit(clonePath, commitish)
	if err != nil {
		return
	}

	var tree *object.Tree
	tree, err = commit.Tree()
	if err != nil {
		return
	}

	entriesByPath = map[string]Entry{"": DirEntry("")}

	err = tree.Files().ForEach(func(file *object.File) (err error) {
		mode := file.Mode
		if !mode.IsFile() || mode.IsMalformed() || !mode.IsRegular() {
			return
		}
		filePath := file.Name
		fileName := ExtractBaseName(filePath)
		parentPath := ExtractDirPath(filePath)
		ensureAncestors(parentPath, entriesByPath)
		entriesByPath[file.Name] = FileEntry(fileName, parentPath, file.Size)
		return
	})

	return
}

func FileContents(clonePath string, commitish string, filePath string) (contents string, err error) {
	var commit *object.Commit
	commit, err = getCommit(clonePath, commitish)
	if err != nil {
		return
	}

	var file *object.File
	file, err = commit.File(filePath)
	if err != nil {
		return
	}

	if file.Size >= MaxFileSizeBytes {
		err = fmt.Errorf("file size is too large to load to memory - %v at %v/%v/%v", file.Size, clonePath, commitish, filePath)
		return
	}

	contents, err = file.Contents()
	return
}
