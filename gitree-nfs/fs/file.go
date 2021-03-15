package fs

import (
	"github.com/go-git/go-billy/v5"
	"gitreefs/git"
	"gitreefs/logger"
	"io"
	"os"
)

type File struct {
	provider  *git.RepositoryProvider
	commitish *Commitish
	fullPath  string
	subPath   string
	position  int64
	size      int64
	isClosed  bool
	isFetched bool
	contents  string
}

var _ billy.File = &File{}

func NewFile(
	fullPath string,
	commitish *Commitish,
	subPath string,
	size int64,
) (file *File, err error) {
	return &File{
		commitish: commitish,
		fullPath:  fullPath,
		subPath:   subPath,
		size:      size,
		position:  0,
	}, nil
}

func (file *File) Name() string {
	return file.fullPath
}

func (file *File) Read(buff []byte) (n int, err error) {
	return file.ReadAt(buff, file.position)
}

func (file *File) ReadAt(buff []byte, offset int64) (int, error) {
	if file.isClosed {
		return 0, os.ErrClosed
	}

	if offset < 0 || offset >= file.size {
		return 0, io.EOF
	}

	targetCapacity := int64(len(buff))
	if offset+targetCapacity > file.size {
		targetCapacity = file.size - offset
	}

	if !file.isFetched {
		contents, err := file.commitish.FileContents(file.subPath)
		if err != nil {
			logger.Error("file.ReadAt: failed for %v: %v", file.fullPath, err)
			return 0, os.ErrNotExist
		}
		file.contents = contents
		file.isFetched = true
	}

	targetContents := file.contents[offset : offset+targetCapacity]
	bytesRead := copy(buff, targetContents)
	file.position += int64(bytesRead)

	return bytesRead, nil
}

func (file *File) Seek(offset int64, whence int) (int64, error) {
	if file.isClosed {
		return 0, os.ErrClosed
	}

	switch whence {
	case io.SeekCurrent:
		file.position += offset
	case io.SeekStart:
		file.position = offset
	case io.SeekEnd:
		file.position = file.size + offset
	}

	return file.position, nil
}

func (file *File) Close() error {
	if file.isClosed {
		return os.ErrClosed
	}

	file.isClosed = true
	return nil
}
