package web3storage

import (
	"context"

	"github.com/jianbo-zh/go-errors"
	ipfsstorage "github.com/jianbo-zh/ipfs-storage"
)

// UploadCar upload car
func (cli *client) UploadCid(ctx context.Context, fcid ipfsstorage.UploadCidParam) (cid string, err error) {
	err = errors.New("web3storage unsupport upload cid")
	return
}
