package fs

import (
	"fmt"
	"github.com/jacobsa/fuse"
	"gitreefs/logger"
)

func Unmount(mountPoint string) error {
	err := fuse.Unmount(mountPoint)
	if err != nil {
		logger.Error("Failed to unmount at %v: %v", mountPoint, err)
		return err
	}
	logger.Info("Successfully unmounted at %v", mountPoint)
	return nil
}

func Mount(clonesPath string, mountPoint string, isRetry bool) (mountedFs *fuse.MountedFileSystem, err error) {

	fuseServer, err := NewFsServer(clonesPath)
	if err != nil {
		return nil, fmt.Errorf("fs_server.NewFsServer: %w", err)
	}

	mountCfg := &fuse.MountConfig{
		FSName:      "gitree",
		VolumeName:  "gitreefs",
		ReadOnly:    true,
		DebugLogger: logger.DebugLogger(),
		ErrorLogger: logger.ErrorLogger(),
	}

	mountedFs, err = fuse.Mount(mountPoint, fuseServer, mountCfg)
	if err == nil {
		return
	}

	if !isRetry {
		unmountErr := Unmount(mountPoint)
		if unmountErr == nil {
			return Mount(clonesPath, mountPoint, true)
		}
		logger.Error("Failed to unmount at %v after failing to mount: %v", mountPoint, err)
	}
	return nil, fmt.Errorf("gitree-fuse.Mount failed: %w", err)
}
