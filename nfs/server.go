package main

import (
	"fmt"
	"github.com/willscott/go-nfs"
	"gitreefs/core/logger"
	"gitreefs/core/virtualfs/bfs"
	"net"
)

func Serve(clonesPath string, host string, port string, storagePath string) error {
	listener, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on %v: %v", port, err)
	}

	logger.Info("nfs server running at %s and mirroring git clones at %v", listener.Addr(), clonesPath)

	fileSystem, err := bfs.NewGitFileSystem(clonesPath)
	if err != nil {
		return fmt.Errorf("failed to create fuseserver on %v: %v", clonesPath, err)
	}

	handler, err := NewHandler(fileSystem, storagePath)
	if err != nil {
		return fmt.Errorf("failed to create handler on %v: %v", clonesPath, err)
	}

	return nfs.Serve(listener, handler, logger.DebugLogger(), logger.InfoLogger())
}
