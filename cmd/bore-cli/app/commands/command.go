package commands

import (
	"errors"

	"github.com/urfave/cli/v2"
	"go.trulyao.dev/bore/v2"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/config"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/handler"
	"go.trulyao.dev/bore/v2/cmd/bore-cli/app/tui"
)

var (
	ErrHandlerNotProvided       = errors.New("handler is not provided")
	ErrBoreNotProvided          = errors.New("bore is not provided")
	ErrTUIManagerNotProvided    = errors.New("tui manager is not provided")
	ErrConfigManagerNotProvided = errors.New("config manager is not provided")
)

type NewCommandOptions struct {
	Handler *handler.Handler
	Bore    *bore.Bore
	TUI     *tui.Manager
	Config  *config.Manager
	Version string
}

type Manager struct {
	handler *handler.Handler
	bore    *bore.Bore
	tui     *tui.Manager
	config  *config.Manager
	version string
}

type SubCommand interface {
	Name() string
	Usage() string
	Execute(ctx *cli.Context) error
}

func NewManager(options *NewCommandOptions) (*Manager, error) {
	if options.Handler == nil {
		return nil, ErrHandlerNotProvided
	}

	if options.Bore == nil {
		return nil, ErrBoreNotProvided
	}

	if options.TUI == nil {
		return nil, ErrTUIManagerNotProvided
	}

	if options.Config == nil {
		return nil, ErrConfigManagerNotProvided
	}

	if options.Version == "" {
		options.Version = "dev"
	}

	return &Manager{
		handler: options.Handler,
		bore:    options.Bore,
		tui:     options.TUI,
		config:  options.Config,
		version: options.Version,
	}, nil
}

func (c *Manager) Execute(ctx *cli.Context, sub SubCommand) error {
	if c == nil {
		return cli.Exit("command is not initialized", 1)
	}

	return sub.Execute(ctx)
}
