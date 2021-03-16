package bfs

import (
	"path/filepath"
	"strings"
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

func split(path string) []string {
	if len(path) == 0 {
		return []string{}
	}
	return strings.Split(path, string(filepath.Separator))
}

func breakdown(fullPath string) (components *pathComponents, err error) {
	parts := split(fullPath)
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
