package nftstorage

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/jianbo-zh/go-errors"
	ipfsstorage "github.com/jianbo-zh/ipfs-storage"
)

func (c *client) Status(ctx context.Context, sname ipfsstorage.Name, cid cid.Cid) (ipfsstorage.PinStatus, error) {

	if c.name != sname {
		return 0, errors.New("storage name error")
	}

	return 0, nil
}
