package nftstorage

import (
	"context"
	"io"

	"github.com/ipfs/go-cid"
)

func (cli *client) Put(ctx context.Context, r io.Reader) (cid.Cid, error) {
	return cid.Cid{}, nil
}
