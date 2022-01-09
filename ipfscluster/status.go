package ipfscluster

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jianbo-zh/go-errors"
	ipfsstorage "github.com/jianbo-zh/ipfs-storage"
)

const (
	TrackerStatusUndefined            = "undefined"
	TrackerStatusClusterError         = "cluster_error"
	TrackerStatusPinError             = "pin_error"
	TrackerStatusUnpinError           = "unpin_error"
	TrackerStatusError                = "error"
	TrackerStatusPinned               = "pinned"
	TrackerStatusPinning              = "pinning"
	TrackerStatusUnpinning            = "unpinning"
	TrackerStatusUnpinned             = "unpinned"
	TrackerStatusRemote               = "remote"
	TrackerStatusPinQueued            = "pin_queued"
	TrackerStatusUnpinQueued          = "unpin_queued"
	TrackerStatusQueued               = "queued"
	TrackerStatusSharded              = "sharded"
	TrackerStatusUnexpectedlyUnpinned = "unexpectedly_unpinned"
)

func (cli *client) Status(ctx context.Context, cid string) (pinStatus string, err error) {

	url, _ := url.Parse(cli.conf.endpoint + "/pins/" + cid)

	req := http.Request{
		URL: url,
		Header: http.Header{
			"Authorization": {"Bearer " + cli.conf.accesstoken},
			"Accept":        {"application/json"},
		},
		Method: http.MethodGet,
	}

	httpCli := http.Client{}

	response, err := httpCli.Do(&req)
	if err != nil {
		err = errors.New("http request error", errors.WithError(err))
		return
	}
	defer response.Body.Close()

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New("ioutile read body error", errors.WithError(err))
		return
	}

	if response.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("[%d]%s", response.StatusCode, string(resBytes)))
		return
	}

	var res StatusResponse200
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		err = errors.New(
			"json unmarshal response error",
			errors.WithError(err),
			errors.WithContext(errors.Context{
				"response": string(resBytes),
			}),
		)
		return
	}

	iPinStatus := ipfsstorage.PIN_STATUS_UNKOWN
	for _, pinTrack := range res.PeerMap {
		if pinTrack.Status == TrackerStatusPinned {
			iPinStatus = ipfsStorageStatus(pinTrack.Status)
			break
		}
		iPinStatus = ipfsStorageStatus(pinTrack.Status)
	}

	fmt.Printf("%s", string(resBytes))

	return iPinStatus, nil
}

// ipfsStorageStatus
func ipfsStorageStatus(icPinstatus string) string {
	switch icPinstatus {
	case TrackerStatusPinQueued:
		return ipfsstorage.PIN_STATUS_QUEUED
	case TrackerStatusPinning:
		return ipfsstorage.PIN_STATUS_PINNING
	case TrackerStatusPinned:
		return ipfsstorage.PIN_STATUS_PINNED
	default:
		return ipfsstorage.PIN_STATUS_UNKOWN
	}
}
