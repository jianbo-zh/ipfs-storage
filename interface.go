package ipfsstorage

import (
	"context"
	"io"

	"github.com/ipfs/go-cid"
)

type Cid = cid.Cid
type Name = string
type Mime = string

type Client interface {
	Name() Name
	Valid(context.Context, Mime) bool
	Put(context.Context, io.Reader) (Cid, error)
	Status(context.Context, Name, Cid) (PinStatus, error)
}
