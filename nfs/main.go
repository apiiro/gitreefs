package main

import (
	"github.com/urfave/cli"
	"gitreefs/core/common"
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
   {{.Name}} {{if .Flags}}[Options]{{end}} clones-path storage-path [port]

ARGS:
    clones-path{{ "\t" }}path to a directory containing git clones (with .git in them)
    storage-path{{ "\t" }}path to a directory in which to keep persistent storage (file handler mapping)
    port{{ "\t" }}(optional) to serve the server at, defaults to 2049

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
	storagePath := opts.(*options).storagePath
	port := opts.(*options).port
	return Serve(clonesPath, "", port, storagePath)
}
