package web3storage

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
	PIN_STATUS_QUEUED  = "PinQueued"
	PIN_STATUS_PINNING = "Pinning"
	PIN_STATUS_PINNED  = "Pinned"
)

func (cli *client) Status(ctx context.Context, cid string) (pinStatus ipfsstorage.PinStatus, err error) {

	url, _ := url.Parse(cli.conf.endpoint + "/status/" + cid)

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
		err = errors.New("http request error", errors.WithError(err))
		return
	}
	defer response.Body.Close()

	resBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = errors.New("ioutile read body error", errors.WithError(err))
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
	for _, pin := range res.Pins {
		if pin.Status == PIN_STATUS_PINNED {
			iPinStatus = ipfsStorageStatus(pin.Status)
			break
		}
		iPinStatus = ipfsStorageStatus(pin.Status)
	}

	return iPinStatus, nil
}

// ipfsStorageStatus
func ipfsStorageStatus(web3Pinstatus string) string {
	switch web3Pinstatus {
	case PIN_STATUS_QUEUED:
		return ipfsstorage.PIN_STATUS_QUEUED
	case PIN_STATUS_PINNING:
		return ipfsstorage.PIN_STATUS_PINNING
	case PIN_STATUS_PINNED:
		return ipfsstorage.PIN_STATUS_PINNED
	default:
		return ipfsstorage.PIN_STATUS_UNKOWN
	}
}
