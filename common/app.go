package common

import (
	"fmt"
	"github.com/urfave/cli"
	"gitreefs/logger"
	"os"
)

type App interface {
	Initialize() *cli.App
	ParseOptions(*cli.Context) (opts Options, err error)
	RunUntilStopped(opts Options) error
}

func RunApp(app App) {
	cliApp := app.Initialize()

	var internalErr error
	cliApp.Action = func(context *cli.Context) {
		internalErr = runWithContext(context, app)
	}

	runErr := cliApp.Run(os.Args)
	exitError := false
	if runErr != nil {
		logger.Error("%v: %v", cliApp.Name, runErr)
		exitError = true
	}
	if internalErr != nil {
		logger.Error("%v: %v", cliApp.Name, internalErr)
		exitError = true
	}
	logger.CloseLoggers()
	if exitError {
		os.Exit(1)
	}
	return
}

func runWithContext(ctx *cli.Context, app App) error {

	opts, err := app.ParseOptions(ctx)
	if err != nil {
		return fmt.Errorf("parsing options: %w", err)
	}

	err = logger.InitLoggers(opts.LogFile(), opts.LogLevel(), ctx.App.Version)
	if err != nil {
		return fmt.Errorf("init loggers: %w", err)
	}

	return app.RunUntilStopped(opts)
}
