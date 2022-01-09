package web3storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jianbo-zh/go-errors"
	ipfsstorage "github.com/jianbo-zh/ipfs-storage"
)

func (cli *client) Upload(ctx context.Context, file ipfsstorage.UploadParam) (cid string, err error) {

	if file.Size > MAX_REQUEST_BODY_SIZE {
		err = ErrRequestBodyLimit
		return
	}

	rurl, _ := url.Parse(cli.conf.endpoint + "/upload")

	req := http.Request{
		URL: rurl,
		Header: http.Header{
			"Authorization": {"Bearer " + cli.conf.accesstoken},
			// "Content-Type":  {"application/car"},
			"Accept": {"application/json"},
			"X-NAME": {url.QueryEscape(file.Name)},
		},
		Method: http.MethodPost,
		Body:   io.NopCloser(file.IOReader),
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

	var res Response200
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

	fmt.Printf("%s", string(resBytes))

	return res.CID, nil
}
