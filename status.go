package ipfsstorage

import (
	"context"

	"github.com/jianbo-zh/go-errors"
)

type PinStatus = int

const (
	PIN_STATUS_PINNED   = 1 // pinned
	PIN_STATUS_UNPINNED = 0 // unpinned
)

func (c *client) Status(ctx context.Context, sname Name, cid Cid) (PinStatus, error) {
	for _, cli := range c.clis {
		if sname == cli.Name() {
			return cli.Status(ctx, sname, cid)
		}
	}

	return 0, errors.New("not found storage name")
}
