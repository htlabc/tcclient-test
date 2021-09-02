package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	errors "githup.com/htl/tcclienttest/pkg/errors"
	log "githup.com/htl/tcclienttest/pkg/log"
	"os"
)

var (
	progressMessage = color.GreenString("==>")
)

type App struct {
	name     string
	basename string
	//name        string
	description string
	options     CliOptions
	runFunc     RunFunc
	silence     bool
	//noVersion   bool
	//noConfig    bool
	commands []*Command
	args     cobra.PositionalArgs
	cmd      *cobra.Command
}

type Option func(*App)

type RunFunc func(basename string) error

func NewApp(name string, basename string, opts ...Option) *App {
	a := &App{
		name:     name,
		basename: basename,
	}

	for _, o := range opts {
		o(a)
	}

	a.buildCommand()

	return a
}

func (a *App) buildCommand() {
	cmd := &cobra.Command{
		Use:   a.basename,
		Short: a.name,
		Long:  a.description,
		// stop printing usage when the command errors
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          a.args,
	}

	cmd.SetOut(os.Stdout)
	cmd.SetErr(os.Stderr)
	cmd.Flags().SortFlags = true
	//cliflag.InitFlags(cmd.Flags())

	if len(a.commands) > 0 {
		for _, command := range a.commands {
			cmd.AddCommand(command.cobraCommand())
		}
		//cmd.SetHelpCommand()
	}

	if a.runFunc != nil {
		cmd.RunE = a.runCommand
	}

	cmd.SetHelpFunc(func(command *cobra.Command, strings []string) {

	})
	cmd.SetUsageFunc(func(command *cobra.Command) error {
		return nil
	})

	a.cmd = cmd

}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {

	printWorkingDir()
	if a.options != nil {
		if err := a.applyOptionRules(); err != nil {
			return err
		}
	}
	// run application
	if a.runFunc != nil {
		return a.runFunc(a.basename)
	}

	return nil
}

func (a *App) Run() {
	if err := a.cmd.Execute(); err != nil {
		fmt.Printf("%v %v\n", "Error:", err)
		os.Exit(1)
	}
}

func WithOptions(opt CliOptions) Option {
	return func(a *App) {
		a.options = opt
	}
}

func WithRunFunc(run RunFunc) Option {
	return func(a *App) {
		a.runFunc = run
	}
}

func (a *App) applyOptionRules() error {
	if completeableOptions, ok := a.options.(CompleteableOptions); ok {
		if err := completeableOptions.Complete(); err != nil {
			return err
		}
	}

	if errs := a.options.Validate(); len(errs) != 0 {
		return errors.NewAggregate(errs)
	}

	if printableOptions, ok := a.options.(PrintableOptions); ok && !a.silence {
		log.Infof("", progressMessage, printableOptions.String())
	}

	return nil
}

func printWorkingDir() {
	wd, _ := os.Getwd()
	log.Infof("WorkingDir: %s", wd)
}
