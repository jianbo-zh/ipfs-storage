package ipfsstorage

import (
	"context"
	"io"

	"github.com/ipfs/go-cid"
	"github.com/jianbo-zh/go-errors"
)

func (c *client) Put(ctx context.Context, r io.Reader) (cid.Cid, error) {
	mime := ""
	for _, cli := range c.clis {
		if ok := cli.Valid(ctx, mime); ok {
			return cli.Put(ctx, r)
		}
	}

	return cid.Cid{}, errors.New("no valid client")
}
