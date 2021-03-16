package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gitreefs/core/common"
	"os"
	"path"
)

type options struct {
	logFile    string
	logLevel   string
	clonesPath string
	storagePath string
	port       string
}

var _ common.Options = &options{}

func (opts *options) LogFile() string {
	return opts.logFile
}

func (opts *options) LogLevel() string {
	return opts.logLevel
}

func (app *NfsApp) ParseOptions(ctx *cli.Context) (common.Options, error) {

	var err error
	opts := &options{
		logFile:   ctx.String("log-file"),
		logLevel:  ctx.String("log-level"),
		port:      "2049",
	}

	args := ctx.Args()
	argsLen := len(args)
	if argsLen < 2 {
		exeName := path.Base(os.Args[0])
		return nil, fmt.Errorf("%s takes two to three arguments. Run `%s --help` for more info", exeName, exeName)
	}

	opts.clonesPath = args[0]
	opts.storagePath = args[1]

	if argsLen >= 3 {
		opts.port = args[2]
	}

	err = common.ValidateDirectory(opts.clonesPath, false)
	if err != nil {
		return nil, err
	}

	err = common.ValidateDirectory(opts.storagePath, true)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
