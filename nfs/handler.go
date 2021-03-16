package main

import (
	"github.com/dgraph-io/badger"
	"github.com/go-git/go-billy/v5"
	"github.com/google/uuid"
	"github.com/willscott/go-nfs"
	"github.com/willscott/go-nfs/helpers"
	"gitreefs/core/logger"
	"math"
	"path/filepath"
	"strings"
)

const (
	PathSeparator       = string(filepath.Separator)
	RootPathPlaceholder = "/"
)

type Handler struct {
	nfs.Handler
	fs billy.Filesystem
	db *badger.DB // stores both ways: path <---> fileHandleId , both as []byte
}

var _ nfs.Handler = &Handler{}

func NewHandler(fs billy.Filesystem, dataPath string) (*Handler, error) {
	opts := badger.DefaultOptions(dataPath)
	opts.Logger = &badgerLogger{}
	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return &Handler{
		Handler: helpers.NewNullAuthHandler(fs),
		db:      db,
		fs:      fs,
	}, nil
}

func newFileHandle() ([]byte, error) {
	return uuid.New().MarshalBinary()
}

func (handler *Handler) lookup(key []byte, addIfMissing bool, valueCreator func() ([]byte, error)) (value []byte, err error) {
	err = handler.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			switch err {
			case badger.ErrKeyNotFound:
				item = nil
			default:
				return err
			}
		}

		if item != nil {
			return item.Value(func(storedValue []byte) error {
				value = make([]byte, len(storedValue))
				copy(value, storedValue)
				return nil
			})
		} else if !addIfMissing {
			return nil
		}

		value, err = valueCreator()
		if err != nil {
			return err
		}

		err = txn.Set(key, value)
		if err != nil {
			return err
		}

		err = txn.Set(value, key)
		return err
	})
	return
}

func (handler *Handler) ToHandle(_ billy.Filesystem, path []string) []byte {
	fullPath := filepath.Join(path...)
	if len(fullPath) == 0 {
		fullPath = RootPathPlaceholder
	}
	handle, err := handler.lookup([]byte(fullPath), true, newFileHandle)
	if err != nil || handle == nil {
		logger.Error("handler.ToHandle: failed for '%v': %v", fullPath, err)
		return nil
	}
	return handle
}

func (handler *Handler) FromHandle(handle []byte) (fs billy.Filesystem, path []string, err error) {
	fs = handler.fs

	var fullPath []byte
	fullPath, err = handler.lookup(handle, false, nil)
	if err != nil || fullPath == nil {
		logger.Info("handler.ToHandle: could not resolve handle '%v': %v", handle, err)
		return nil, []string{}, &nfs.NFSStatusError{NFSStatus: nfs.NFSStatusStale}
	}

	fullPathStr := string(fullPath)
	if fullPathStr == RootPathPlaceholder {
		path = []string{""}
	} else {
		path = strings.Split(fullPathStr, PathSeparator)
	}
	return
}

func (handler Handler) HandleLimit() int {
	return math.MaxInt32
}

type badgerLogger struct {
}

func (bLogger *badgerLogger) Errorf(s string, i ...interface{}) {
	logger.Error(s, i...)
}

func (bLogger *badgerLogger) Warningf(s string, i ...interface{}) {
	logger.Error(s, i...)
}

func (bLogger *badgerLogger) Infof(s string, i ...interface{}) {
	logger.Info(s, i...)
}

func (bLogger *badgerLogger) Debugf(s string, i ...interface{}) {
	logger.Debug(s, i...)
}

var _ badger.Logger = &badgerLogger{}
