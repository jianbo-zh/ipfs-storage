package nftstorage

import (
	"context"

	ipfsstorage "github.com/jianbo-zh/ipfs-storage"
)

func (cli *client) Valid(ctx context.Context, mimeType ipfsstorage.Mime) bool {
	return true
}
