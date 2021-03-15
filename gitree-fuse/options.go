package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gitreefs/common"
	"os"
	"path"
	"path/filepath"
)

type options struct {
	logFile    string
	logLevel   string
	clonesPath string
	mountPoint string
}

var _ common.Options = &options{}

func (opts *options) LogFile() string {
	return opts.logFile
}

func (opts *options) LogLevel() string {
	return opts.logLevel
}

func (app *FuseApp) ParseOptions(ctx *cli.Context) (opts common.Options, err error) {
	var clonesPath, mountPoint string
	switch len(ctx.Args()) {

	case 2:
		clonesPath = ctx.Args()[0]
		mountPoint = ctx.Args()[1]

	default:
		exeName := path.Base(os.Args[0])
		err = fmt.Errorf("%s takes exactly two arguments. Run `%s --help` for more info", exeName, exeName)

		return
	}

	err = common.ValidateDirectory(clonesPath, false)
	if err != nil {
		return
	}

	mountPoint, err = filepath.Abs(mountPoint)
	if err != nil {
		err = fmt.Errorf("canonicalizing mountFs point: %v", err)
		return
	}

	err = common.ValidateDirectory(mountPoint, true)
	if err != nil {
		return
	}

	opts = &options{
		logFile:    ctx.String("log-file"),
		logLevel:   ctx.String("log-level"),
		clonesPath: clonesPath,
		mountPoint: mountPoint,
	}
	return
}
