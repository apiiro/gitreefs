package git

import "path"

type Entry struct {
	Size          int64
	IsDir         bool
	EntriesByName map[string]*Entry
}

type RootEntry struct {
	Entry
	EntriesByPath map[string]*Entry
}

func FileEntry(parent *Entry, name string, size int64) (entry *Entry) {
	entry = &Entry{
		Size:          size,
		EntriesByName: nil,
	}
	parent.addEntry(entry, name)
	return
}

func DirEntry(parent *Entry, name string) (entry *Entry) {
	entry = &Entry{
		EntriesByName: make(map[string]*Entry),
		IsDir:         true,
	}
	if parent != nil {
		parent.addEntry(entry, name)
	}
	return
}

func (entry *Entry) IsRoot() bool {
	return false
}

func (root *RootEntry) IsRoot() bool {
	return true
}

func (entry *Entry) addEntry(child *Entry, name string) {
	entry.EntriesByName[name] = child
}

func ExtractBaseName(fromPath string) string {
	return path.Base(fromPath)
}

func ExtractDirPath(fromPath string) (dirPath string) {
	dirPath = path.Dir(fromPath)
	if dirPath == "." {
		return ""
	}
	return
}

func (root *RootEntry) ensurePath(fullPath string) (entry *Entry) {

	var found bool
	entry, found = root.EntriesByPath[fullPath]
	if found {
		return entry
	}

	dirPath := ExtractDirPath(fullPath)
	parentEntry := root.ensurePath(dirPath)

	entry = DirEntry(parentEntry, ExtractBaseName(fullPath))
	root.EntriesByPath[fullPath] = entry
	return
}
