package main

import (
	"fmt"
	"github.com/urfave/cli"
	"github.com/willscott/go-nfs"
	nfsHelper "github.com/willscott/go-nfs/helpers"
	"gitreefs/core/common"
	"gitreefs/core/logger"
	"gitreefs/core/virtualfs/bfs"
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
	cli.AppHelpTemplate =
		`NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[Options]{{end}} clones-path [port] [cacheSize]

ARGS:
    clones-path{{ "\t" }}path to a directory containing git clones (with .git in them)
    port{{ "\t" }}(optional) to serve the server at, defaults to 2049
    cacheSize{{ "\t" }}(optional) size of file handlers cache, defaults to 1024

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`

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

	fileSystem, err := bfs.NewGitFileSystem(clonesPath)
	if err != nil {
		return fmt.Errorf("failed to create fuseserver on %v: %v", clonesPath, err)
	}

	handler := nfsHelper.NewNullAuthHandler(fileSystem)
	cacheHelper := nfsHelper.NewCachingHandler(handler, cacheSize)
	return nfs.Serve(listener, cacheHelper)
}
