package toolbs

import (
	"context"

	"github.com/pushthat/bud/framework/web/webrt"
	"github.com/pushthat/bud/package/budserver"
	v8 "github.com/pushthat/bud/package/js/v8"
	"github.com/pushthat/bud/package/socket"

	"github.com/pushthat/bud/framework"
	"github.com/pushthat/bud/internal/cli/bud"
	"github.com/pushthat/bud/internal/pubsub"
)

func New(bud *bud.Command, in *bud.Input) *Command {
	return &Command{
		bud:  bud,
		in:   in,
		Flag: new(framework.Flag),
	}
}

type Command struct {
	bud  *bud.Command
	in   *bud.Input
	Flag *framework.Flag
}

func (c *Command) Run(ctx context.Context) error {
	log, err := bud.Log(c.in.Stdout, c.bud.Log)
	if err != nil {
		return err
	}
	module, err := bud.Module(c.bud.Dir)
	if err != nil {
		return err
	}
	vm, err := v8.Load()
	if err != nil {
		return err
	}
	// Load the file server
	servefs, err := bud.FileServer(log, module, vm, c.Flag)
	if err != nil {
		return err
	}
	bus := pubsub.New()
	server := budserver.New(servefs, bus, log, vm)
	budln, err := socket.Listen(":35729")
	if err != nil {
		return err
	}
	defer budln.Close()
	log.Info("Listening on http://127.0.0.1:35729")
	return webrt.Serve(ctx, budln, server)
}
