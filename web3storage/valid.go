package web3storage

import (
	"context"

	ipfsstorage "github.com/jianbo-zh/ipfs-storage"
)

func (cli *client) Valid(ctx context.Context, mimeType ipfsstorage.MimeType) bool {
	return true
}
