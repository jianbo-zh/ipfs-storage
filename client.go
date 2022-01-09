package ipfsstorage

import (
	"context"
)

type ClientHub struct {
	clis []Client
}

func NewClientHub(clis ...Client) *ClientHub {
	cli := ClientHub{clis: clis}
	return &cli
}

func (c *ClientHub) Upload(ctx context.Context, storages []StorageName, file UploadParam) (ress []UploadResult, errs []error) {
	for _, cli := range c.clis {
		for _, sname := range storages {
			if cli.Name() == sname && cli.Valid(ctx, file.Mime) {
				res, err := cli.Upload(ctx, file)
				if err != nil {
					errs = append(errs, err)
				}
				ress = append(ress, UploadResult{Storage: sname, CID: res})
			}
		}
	}

	return ress, errs
}

func (c *ClientHub) UploadCar(ctx context.Context, storages []StorageName, fcar UploadCarParam) (ress []UploadResult, errs []error) {
	for _, cli := range c.clis {
		for _, sname := range storages {
			if cli.Name() == sname && cli.Valid(ctx, fcar.Mime) {
				res, err := cli.UploadCar(ctx, fcar)
				if err != nil {
					errs = append(errs, err)
				}
				ress = append(ress, UploadResult{Storage: sname, CID: res})
			}
		}
	}

	return ress, errs
}

func (c *ClientHub) UploadCid(ctx context.Context, storages []StorageName, fcid UploadCidParam) (ress []UploadResult, errs []error) {
	for _, cli := range c.clis {
		for _, sname := range storages {
			if cli.Name() == sname && cli.Valid(ctx, fcid.Mime) {
				res, err := cli.UploadCid(ctx, fcid)
				if err != nil {
					errs = append(errs, err)
				}
				ress = append(ress, UploadResult{Storage: sname, CID: res})
			}
		}
	}

	return
}

func (c *ClientHub) Delete(ctx context.Context, storages []StorageName, cid Cid) (ress []DeleteResult, errs []error) {
	for _, cli := range c.clis {
		for _, sname := range storages {
			if sname == cli.Name() {
				res, err := cli.Delete(ctx, cid)
				if err != nil {
					errs = append(errs, err)
				}
				ress = append(ress, DeleteResult{Storage: sname, OK: res})
			}
		}
	}

	return
}

func (c *ClientHub) Status(ctx context.Context, storages []StorageName, cid Cid) (ress []StatusResult, errs []error) {
	for _, cli := range c.clis {
		for _, sname := range storages {
			if sname == cli.Name() {
				res, err := cli.Status(ctx, cid)
				if err != nil {
					errs = append(errs, err)
				}
				ress = append(ress, StatusResult{Storage: sname, Status: res})
			}
		}
	}

	return
}
