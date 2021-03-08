package gitreefs

import (
	"context"
	"fmt"
	"github.com/jacobsa/daemonize"
	"github.com/jacobsa/fuse"
	"github.com/kardianos/osext"
	"github.com/urfave/cli"
	"gitreefs/fs"
	"gitreefs/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app := Init()

	var err error
	app.Action = func(context *cli.Context) {
		err = run(context)
	}

	err = app.Run(os.Args)
	if err != nil {
		return
	}

	if err != nil {
		logger.Error("main: %w", err)
		os.Exit(1)
	}
}

func dispatchDaemon() error {
	var err error
	var executablePath string
	executablePath, err = osext.Executable()
	if err != nil {
		return fmt.Errorf("osext.Executable: %w", err)
	}

	args := append([]string{"--foreground"}, os.Args[1:]...)

	env := []string{
		fmt.Sprintf("PATH=%s", os.Getenv("PATH")),
	}

	err = daemonize.Run(executablePath, args, env, os.Stdout)
	if err != nil {
		return fmt.Errorf("daemonize.Run: %w", err)
	}

	return nil
}

func run(ctx *cli.Context) error {
	var err error

	var opts = ParseOptions(ctx, err)
	if err != nil {
		return fmt.Errorf("parsing options: %w", err)
	}

	err = logger.InitLoggers(opts.LogFile, opts.LogLevel)
	if err != nil {
		return fmt.Errorf("init log file: %w", err)
	}

	logger.Info("Using mountFs: %v --> %v", opts.ClonesPath, opts.MountPoint)

	if !opts.Foreground {
		return dispatchDaemon()
	}

	var mountedFs *fuse.MountedFileSystem
	{
		mountedFs, err = mountFs(opts)

		if err == nil {
			logger.Info("File system has been successfully mounted.")
			daemonize.SignalOutcome(nil)
		} else {
			err = fmt.Errorf("mountFs: %w", err)
			daemonize.SignalOutcome(err)
			return err
		}
	}

	registerSignalHandler(mountedFs.Dir())

	err = mountedFs.Join(context.Background())
	if err != nil {
		return fmt.Errorf("MountedFileSystem.Join: %w", err)
	}

	logger.Info("File system has been successfully un-mounted.")

	return nil
}

func mountFs(opts *options) (mountedFs *fuse.MountedFileSystem, err error) {

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
	if err != nil {
		return nil, fmt.Errorf("fuse.Mount: %w", err)
	}

	return
}

func registerSignalHandler(mountPoint string) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		for {
			<-signalChan
			logger.Info("Received SIGINT, attempting to unmount...")

			err := fuse.Unmount(mountPoint)
			if err != nil {
				logger.Error("Failed to unmount in response to SIGINT: %w", err)
			} else {
				logger.Info("Successfully unmounted in response to SIGINT.")
				return
			}
		}
	}()
}
