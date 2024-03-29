package main

import (
	"fmt"
	"github.com/jacobsa/fuse"
	"github.com/urfave/cli"
	"gitreefs/core/common"
	"gitreefs/core/logger"
	"gitreefs/fuse/fuseserver"
	"golang.org/x/net/context"
	"os"
	"os/signal"
	"syscall"
)

type FuseApp struct {
}

func main() {
	var app common.App = &FuseApp{}
	common.RunApp(app)
}

func (app *FuseApp) DeclareCli() *cli.App {
	cli.AppHelpTemplate =
		`NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[Options]{{end}} clones-path mount-point

ARGS:
    clones-path{{ "\t" }}path to a directory containing git clones (with .git in them)
    mount-point{{ "\t" }}path to target location to mount the virtual fs at

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`
	return &cli.App{
		Name:    "gitreefs-fuse",
		Version: Version,
		Usage:   "Mount a forest of git trees as a virtual file system backed by FUSE",
		Writer:  os.Stdout,
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

func (app *FuseApp) RunUntilStopped(opts common.Options) (err error) {

	clonesPath := opts.(*options).clonesPath
	mountPoint := opts.(*options).mountPoint
	logger.Info("Mounting: %v --> %v", clonesPath, mountPoint)

	var mountedFs *fuse.MountedFileSystem
	{
		mountedFs, err = fuseserver.Mount(clonesPath, mountPoint, false)

		if err == nil {
			logger.Info("fileHandler system has been successfully mounted.")
		} else {
			err = fmt.Errorf("mountFs: %w", err)
			return err
		}
	}

	registerSignalHandler(mountedFs.Dir())

	err = mountedFs.Join(context.Background())
	if err != nil {
		return fmt.Errorf("MountedFileSystem.Join: %w", err)
	}

	logger.Info("fileHandler system has been successfully un-mounted.")

	return nil
}

func registerSignalHandler(mountPoint string) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		for {
			<-signalChan
			logger.Info("Received SIGINT, attempting to unmount...")
			err := fuseserver.Unmount(mountPoint)
			if err != nil {
				logger.Error("Failed to unmount: %v", err)
			}
		}
	}()
}
