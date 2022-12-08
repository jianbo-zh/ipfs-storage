package nftstorage

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
	PIN_STATUS_QUEUED  = "queued"  // queued
	PIN_STATUS_PINNING = "pinning" // pinning
	PIN_STATUS_PINNED  = "pinned"  // pinned
	PIN_STATUS_FAILED  = "failed"  // failed
)

func (cli *client) Status(ctx context.Context, cid string) (pinStatus ipfsstorage.PinStatus, err error) {

	url, _ := url.Parse(cli.conf.endpoint + "/check/" + cid)

	req := http.Request{
		URL: url,
		Header: http.Header{
			"Authorization": {"Bearer " + cli.conf.accesstoken},
			// "Content-Type":  {"application/car"},
			"Accept": {"application/json"},
		},
		Method: http.MethodGet,
	}

	httpCli := http.Client{}

	response, err := httpCli.Do(&req)
	if err != nil {
		err = errors.New("http request error").With(errors.Inner(err))
		return
	}
	defer response.Body.Close()

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New("ioutile read body error").With(errors.Inner(err))
		return
	}

	if response.StatusCode == http.StatusNotFound {
		// not found
		return ipfsstorage.PIN_STATUS_UNPINNED, nil
	}

	if response.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("[%d]%s", response.StatusCode, string(resBytes)))
		return
	}

	var res Response200
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		err = errors.New("json unmarshal response error").
			With(errors.Inner(err), errors.Playload(errors.MapData{"response": string(resBytes)}))
		return
	}

	return ipfsStorageStatus(res.Value.Pin.Status), nil
}

// ipfsStorageStatus
func ipfsStorageStatus(nftPinstatus string) string {
	switch nftPinstatus {
	case PIN_STATUS_QUEUED:
		return ipfsstorage.PIN_STATUS_QUEUED
	case PIN_STATUS_PINNING:
		return ipfsstorage.PIN_STATUS_PINNING
	case PIN_STATUS_PINNED:
		return ipfsstorage.PIN_STATUS_PINNED
	case PIN_STATUS_FAILED:
		return ipfsstorage.PIN_STATUS_FAILED
	default:
		return ipfsstorage.PIN_STATUS_UNKOWN
	}
}
