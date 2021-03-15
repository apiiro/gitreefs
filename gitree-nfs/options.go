package main

import (
	"fmt"
	"github.com/urfave/cli"
	"gitreefs/common"
	"os"
	"path"
	"strconv"
)

type options struct {
	logFile    string
	logLevel   string
	clonesPath string
	port       string
	cacheSize  int
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
		cacheSize: 1024,
	}

	args := ctx.Args()
	argsLen := len(args)
	if argsLen < 1 {
		exeName := path.Base(os.Args[0])
		err = fmt.Errorf("%s takes one to three arguments. Run `%s --help` for more info", exeName, exeName)
	}

	opts.clonesPath = args[0]

	if argsLen >= 2 {
		opts.port = args[1]
	}
	if argsLen >= 3 {
		opts.cacheSize, err = strconv.Atoi(args[2])
		if err != nil {
			return nil, fmt.Errorf(
				"cache size has invalid value: %v", args[2])
		}
	}

	err = common.ValidateDirectory(opts.clonesPath, false)
	if err != nil {
		return nil, err
	}

	return opts, nil
}
