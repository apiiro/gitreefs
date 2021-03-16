package bfs

import (
	"os"
	"time"
)

type statInfo struct {
	name  string
	size  int64
	isDir bool
}

var _ os.FileInfo = &statInfo{}

func (info *statInfo) Name() string {
	return info.name
}

func (info *statInfo) Size() int64 {
	return info.size
}

func (info *statInfo) Mode() os.FileMode {
	if info.isDir {
		return os.ModeDir | os.ModePerm
	}
	return os.ModePerm
}

func (info *statInfo) ModTime() time.Time {
	return time.Now()
}

func (info *statInfo) IsDir() bool {
	return info.isDir
}

func (info *statInfo) Sys() interface{} {
	return nil
}

func statDir(name string) (os.FileInfo, error) {
	return &statInfo{
		name:  name,
		size:  0,
		isDir: true,
	}, nil
}

func statFile(name string, size int64) (os.FileInfo, error) {
	return &statInfo{
		name:  name,
		size:  size,
		isDir: false,
	}, nil
}
