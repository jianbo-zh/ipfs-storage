package nftstorage

import ipfsstorage "github.com/jianbo-zh/ipfs-storage"

type client struct {
	name string
	conf clientConfig
}

func (c *client) Name() string {
	return c.name
}

func NewClient(name string, opts ...Option) (*client, error) {
	conf := clientConfig{
		endpoint: endpoint,
	}

	for _, opt := range opts {
		if err := opt(&conf); err != nil {
			return nil, err
		}
	}

	cli := client{name: name, conf: conf}

	return &cli, nil
}

var _ ipfsstorage.Client = (*client)(nil)
