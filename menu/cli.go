package cli

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/devonboyer/airbot/bot"
)

type cliSource struct {
	msgs    chan bot.Message
	stopped chan struct{}
	wg      sync.WaitGroup
	io.Writer
}

func newCLISource() *cliSource {
	return &cliSource{
		msgs:    make(chan bot.Message),
		stopped: make(chan struct{}),
		wg:      sync.WaitGroup{},
		Writer:  os.Stdout,
	}
}

func (c *cliSource) Messages() <-chan bot.Message {
	return c.msgs
}

func (c *cliSource) Send(reply bot.Reply) {
	fmt.Fprintf(c, reply.Text)
}

func (c *cliSource) stop() {
	close(c.stopped)
	c.wg.Wait()
}
