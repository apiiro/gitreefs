package main

import (
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/urfave/cli"
	"github.com/willscott/go-nfs"
	nfsHelper "github.com/willscott/go-nfs/helpers"
	"gitreefs/common"
	"gitreefs/logger"
	"net"
	"os"
)

type NfsApp struct {
}

func main() {
	var app common.App = &NfsApp{}
	common.RunApp(app)
}

func (app *NfsApp) DeclareCli() *cli.App {
	return &cli.App{
		Name:    "gitreefs-nfs",
		Version: Version,
		Usage:   "NFS server providing access to a forest of git trees as a virtual file system",
		Writer:  os.Stderr,
		Flags: []cli.Flag{

			cli.StringFlag{
				Name:  "log-file",
				Value: "logs/gitreefs-%v-%v.log",
				Usage: "Output logs file path format.",
			},

			cli.StringFlag{
				Name:  "log-level",
				Value: "DEBUG",
				Usage: "Set log level.",
			},
		},
	}
}

func (app *NfsApp) RunUntilStopped(opts common.Options) error {
	clonesPath := opts.(*options).clonesPath
	port := opts.(*options).port
	cacheSize := opts.(*options).cacheSize

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen on %v: %v", port, err)
	}

	logger.Info("nfs server running at %s and mirroring git clones at %v", listener.Addr(), clonesPath)

	var fileSystem billy.Filesystem = nil // TODO

	handler := nfsHelper.NewNullAuthHandler(fileSystem)
	cacheHelper := nfsHelper.NewCachingHandler(handler, cacheSize)
	return nfs.Serve(listener, cacheHelper)
}
