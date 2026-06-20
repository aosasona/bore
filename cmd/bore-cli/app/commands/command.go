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

type Command struct {
	handler *handler.Handler
	bore    *bore.Bore
	tui     *tui.Manager
	config  *config.Manager
	version string
}

type SubCommand interface {
	*Command
	Name() string
	Usage() string
	Execute(ctx *cli.Context) error
}

func New(options *NewCommandOptions) (*Command, error) {
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

	return &Command{
		handler: options.Handler,
		bore:    options.Bore,
		tui:     options.TUI,
		config:  options.Config,
	}, nil
}
