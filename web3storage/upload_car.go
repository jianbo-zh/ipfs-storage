package web3storage

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/alanshaw/go-carbites"
	"github.com/jianbo-zh/go-errors"
	ipfsstorage "github.com/jianbo-zh/ipfs-storage"
)

// UploadCar upload car
func (cli *client) UploadCar(ctx context.Context, fcar ipfsstorage.UploadCarParam) (cid string, err error) {

	if fcar.Size <= MAX_REQUEST_BODY_SIZE {
		// small car file
		return cli.uploadCar(ctx, fcar.IOReader, fcar.Name)
	}

	// big car file
	strategy := carbites.Treewalk
	spltr, err := carbites.Split(fcar.IOReader, MAX_REQUEST_BODY_SIZE, strategy)
	if err != nil {
		err = errors.New("carbites split error", errors.WithError(err))
		return
	}

	for {
		chunkCar, err := spltr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return cid, errors.New("carbites split next error", errors.WithError(err))
		}

		cid, err = cli.uploadCar(ctx, chunkCar, fcar.Name)
		if err != nil {
			return cid, errors.New("upload chunk car error", errors.WithError(err))
		}
	}

	return cid, nil
}

// uploadCar upload small car
func (cli *client) uploadCar(ctx context.Context, r io.Reader, fileName string) (cid string, err error) {
	rurl, _ := url.Parse(cli.conf.endpoint + "/car")

	req := http.Request{
		URL: rurl,
		Header: http.Header{
			"Authorization": {"Bearer " + cli.conf.accesstoken},
			"Content-Type":  {"application/car"},
			"Accept":        {"application/json"},
			"X-NAME":        {url.QueryEscape(fileName)},
		},
		Method: http.MethodPost,
		Body:   io.NopCloser(r),
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
