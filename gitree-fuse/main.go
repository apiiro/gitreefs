package main

import (
	"fmt"
	"github.com/jacobsa/fuse"
	"github.com/urfave/cli"
	"gitreefs/common"
	"gitreefs/gitree-fuse/fs"
	"gitreefs/logger"
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

func (app *FuseApp) Initialize() *cli.App {
	return &cli.App{
		Name:    "gitreefs-fuse",
		Version: Version,
		Usage:   "Mount a forest of git trees as a virtual file system backed by FUSE",
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

func (app *FuseApp) RunUntilStopped(opts common.Options) (err error) {

	clonesPath := opts.(*options).clonesPath
	mountPoint := opts.(*options).mountPoint
	logger.Info("Mounting: %v --> %v", clonesPath, mountPoint)

	var mountedFs *fuse.MountedFileSystem
	{
		mountedFs, err = fs.Mount(clonesPath, mountPoint, false)

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
			err := fs.Unmount(mountPoint)
			if err != nil {
				logger.Error("Failed to unmount: %v", err)
			}
		}
	}()
}
