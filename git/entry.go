package git

import "path"

type Entry struct {
	Name       string
	Size       int64
	IsDir      bool
	ParentPath string
}

func FileEntry(name string, parentPath string, size int64) Entry {
	return Entry{
		Name:       name,
		ParentPath: parentPath,
		Size:       size,
	}
}

func DirEntry(dirPath string) Entry {
	return Entry{
		Name:       ExtractBaseName(dirPath),
		ParentPath: ExtractDirPath(dirPath),
		IsDir:      true,
	}
}

func (entry *Entry) IsRoot() bool {
	return len(entry.ParentPath) == 0
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
