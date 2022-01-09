package ipfsstorage

import (
	"context"
	"io"
)

const (
	PIN_STATUS_QUEUED  = "queued"  // queued
	PIN_STATUS_PINNING = "pinning" // pinning
	PIN_STATUS_PINNED  = "pinned"  // pinned
	PIN_STATUS_FAILED  = "failed"  // failed
	PIN_STATUS_UNKOWN  = "unkown"  // unkown
)

type Cid = string
type StorageName = string
type MimeType = string
type PinStatus = string

type Client interface {
	Name() StorageName
	Valid(context.Context, MimeType) bool
	Upload(context.Context, UploadParam) (Cid, error)
	UploadCar(context.Context, UploadCarParam) (Cid, error)
	UploadCid(context.Context, UploadCidParam) (Cid, error)
	Delete(context.Context, Cid) (bool, error)
	Status(context.Context, Cid) (PinStatus, error)
}

type UploadParam struct {
	IOReader io.Reader `json:"IOReader"`
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	Mime     string    `json:"mime"`
}

type UploadCarParam struct {
	IOReader io.Reader `json:"IOReader"`
	Name     string    `json:"name"`
	Size     int64     `json:"size"`
	Mime     string    `json:"mime"`
}

type UploadCidParam struct {
	CID  string `json:"cid"`
	Name string `json:"name"`
	Size int64  `json:"size"`
	Mime string `json:"mime"`
}

type UploadResult struct {
	CID     string      `json:"cid"`
	Storage StorageName `json:"storage"`
}

type StatusResult struct {
	Status  PinStatus   `json:"status"`
	Storage StorageName `json:"storage"`
}

type DeleteResult struct {
	OK      bool        `json:"ok"`
	Storage StorageName `json:"storage"`
}
