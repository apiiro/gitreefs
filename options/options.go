package options

import (
	"fmt"
	"github.com/urfave/cli"
	"gitreefs/fs"
	"os"
	"path"
	"path/filepath"
)

func init() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[global Options]{{end}} clones mountpoint
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`
}

func Init() (app *cli.App) {
	app = &cli.App{
		Name:    "gitreefs",
		Version: "0.1",
		Usage:   "Mount a forest of git trees as a virtual file system",
		Writer:  os.Stderr,
		Flags: []cli.Flag{

			cli.BoolFlag{
				Name:  "foreground",
				Usage: "Stay in the foreground after mounting.",
			},

			cli.StringFlag{
				Name:  "log-file",
				Value: "logs/gitreefs-%v.log",
				Usage: "Output logs file path format.",
			},

			cli.StringFlag{
				Name:  "log-level",
				Value: "DEBUG",
				Usage: "Set log level.",
			},
		},
	}

	return
}

type Options struct {
	Foreground bool
	LogFile    string
	LogLevel   string
	ClonesPath string
	MountPoint string
}

func ParseOptions(c *cli.Context, err error) (opts *Options) {
	var clonesPath, mountPoint string
	switch len(c.Args()) {

	case 2:
		clonesPath = c.Args()[0]
		mountPoint = c.Args()[1]

	default:
		err = fmt.Errorf(
			"%s takes exactly two arguments. Run `%s --help` for more info",
			path.Base(os.Args[0]),
			path.Base(os.Args[0]),
			)

		return
	}

	err = fs.ValidateDirectory(clonesPath)
	if err != nil {
		return
	}

	mountPoint, err = filepath.Abs(mountPoint)
	if err != nil {
		err = fmt.Errorf("canonicalizing mountFs point: %w", err)
		return
	}

	opts = &Options{
		Foreground: c.Bool("foreground"),
		LogFile:    c.String("log-file"),
		LogLevel:   c.String("log-level"),
		ClonesPath: clonesPath,
		MountPoint: mountPoint,
	}
	return
}

