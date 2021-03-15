package main

import (
	"context"
	"fmt"
	"github.com/jacobsa/fuse"
	"github.com/urfave/cli"
	"gitreefs/gitree-fuse/fs"
	"gitreefs/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := Init()

	var internalErr error
	app.Action = func(context *cli.Context) {
		internalErr = run(context)
	}

	runErr := app.Run(os.Args)
	exitError := false
	if runErr != nil {
		logger.Error("gitree-fuse: %v", runErr)
		exitError = true
	}
	if internalErr != nil {
		logger.Error("gitree-fuse: %v", internalErr)
		exitError = true
	}
	logger.CloseLoggers()
	if exitError {
		os.Exit(1)
	}
	return
}

func run(ctx *cli.Context) error {

	opts, err := ParseOptions(ctx)
	if err != nil {
		return fmt.Errorf("parsing options: %w", err)
	}

	err = logger.InitLoggers(opts.LogFile, opts.LogLevel, Version)
	if err != nil {
		return fmt.Errorf("init log file: %w", err)
	}

	logger.Info("Mounting: %v --> %v", opts.ClonesPath, opts.MountPoint)

	var mountedFs *fuse.MountedFileSystem
	{
		mountedFs, err = fs.Mount(opts.ClonesPath, opts.MountPoint, false)

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
