package web3storage

import (
	"context"

	"github.com/jianbo-zh/go-errors"
)

func (cli *client) Delete(ctx context.Context, cid string) (ok bool, err error) {
	err = errors.New("web3storage unsupport delete api, only login admin")
	return
}
