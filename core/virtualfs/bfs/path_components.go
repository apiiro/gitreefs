package bfs

import (
	"fmt"
	"path/filepath"
)

type pathComponents struct {
	repositoryName string
	commitishName  string
	subPath        string
}

func (components *pathComponents) hasRepository() bool {
	return len(components.repositoryName) > 0
}

func (components *pathComponents) hasCommitish() bool {
	return len(components.commitishName) > 0
}

func breakdown(fullPath string) (components *pathComponents, err error) {
	if len(fullPath) == 0 {
		return nil, fmt.Errorf("breakdown: path is empty")
	}
	parts := filepath.SplitList(fullPath)
	components = &pathComponents{
		repositoryName: "",
		commitishName:  "",
		subPath:        "",
	}
	if len(parts) >= 1 {
		components.repositoryName = parts[0]
	}
	if len(parts) >= 2 {
		components.commitishName = parts[1]
	}
	if len(parts) >= 3 {
		components.subPath = filepath.Join(parts[2:]...)
	}
	return
}
