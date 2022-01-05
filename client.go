package ipfsstorage

type client struct {
	clis []Client
}

func (c *client) Name() string {
	return "merge-client"
}

func NewClient(clis ...Client) (Client, error) {
	cli := client{clis: clis}
	return &cli, nil
}

var _ Client = (*client)(nil)
