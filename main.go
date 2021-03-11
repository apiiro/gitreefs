package main

import (
	"context"
	"fmt"
	"github.com/jacobsa/fuse"
	"github.com/urfave/cli"
	"gitreefs/fs"
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
		logger.Error("main: %v", runErr)
		exitError = true
	}
	if internalErr != nil {
		logger.Error("main: %v", internalErr)
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

	err = logger.InitLoggers(opts.LogFile, opts.LogLevel)
	if err != nil {
		return fmt.Errorf("init log file: %w", err)
	}

	logger.Info("Using mountFs: %v --> %v", opts.ClonesPath, opts.MountPoint)

	var mountedFs *fuse.MountedFileSystem
	{
		mountedFs, err = mountFs(opts, false)

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

func unmount(mountPoint string) error {
	err := fuse.Unmount(mountPoint)
	if err != nil {
		logger.Error("Failed to unmount in response to SIGINT: %v", err)
		return err
	} else {
		logger.Info("Successfully unmounted in response to SIGINT.")
	}
	return nil
}

func mountFs(opts *Options, isRetry bool) (mountedFs *fuse.MountedFileSystem, err error) {

	fuseServer, err := fs.NewFsServer(opts.ClonesPath)
	if err != nil {
		return nil, fmt.Errorf("fs_server.NewFsServer: %w", err)
	}

	mountCfg := &fuse.MountConfig{
		FSName:      "gitree",
		VolumeName:  "gitreefs",
		ReadOnly:    true,
		DebugLogger: logger.DebugLogger,
		ErrorLogger: logger.ErrorLogger,
	}

	mountedFs, err = fuse.Mount(opts.MountPoint, fuseServer, mountCfg)
	if err == nil {
		return
	}

	if !isRetry {
		unmountErr := unmount(opts.MountPoint)
		if unmountErr == nil {
			return mountFs(opts, true)
		}
		logger.Error("Failed to unmount at %v after failing to mount: %v")
	}
	return nil, fmt.Errorf("fuse.Mount failed: %w", err)
}

func registerSignalHandler(mountPoint string) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		for {
			<-signalChan
			logger.Info("Received SIGINT, attempting to unmount...")
			err := unmount(mountPoint)
			if err != nil {
				logger.Error("Failed to unmount: %v", err)
			}
		}
	}()
}
