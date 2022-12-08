package ipfscluster

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"

	"github.com/alanshaw/go-carbites"
	"github.com/jianbo-zh/go-errors"
	ipfsstorage "github.com/jianbo-zh/ipfs-storage"
)

// UploadCar upload car
func (cli *client) UploadCar(ctx context.Context, fcar ipfsstorage.UploadCarParam) (cid string, err error) {

	if fcar.Size <= MAX_REQUEST_BODY_SIZE {
		// small car file, upload directly
		ofile, err := os.Open(fcar.FilePath)
		if err != nil {
			return "", errors.New("open car file error").With(errors.Inner(err))
		}
		defer ofile.Close()

		return cli.uploadCar(ctx, ofile, fcar.Name)
	}

	// big car file
	spltr, err := carbites.NewTreewalkSplitterFromPath(fcar.FilePath, MAX_REQUEST_BODY_SIZE)
	if err != nil {
		err = errors.New("carbites split error").With(errors.Inner(err))
		return
	}

	for {
		chunkCar, err := spltr.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return cid, errors.New("carbites split next error").With(errors.Inner(err))
		}

		cid, err = cli.uploadCar(ctx, chunkCar, fcar.Name)
		if err != nil {
			return cid, errors.New("upload chunk car error").With(errors.Inner(err))
		}
	}

	return cid, nil
}

func (cli *client) uploadCar(ctx context.Context, fileReader io.Reader, fileName string) (cid string, err error) {

	url, _ := url.Parse(cli.conf.endpoint + "/add?cid-version=1&format=car")

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", fileName)
	if err != nil {
		return cid, errors.New("create form file error").With(errors.Inner(err))
	}

	_, err = io.Copy(fileWriter, fileReader)
	if err != nil {
		return cid, errors.New("io copy error").With(errors.Inner(err))
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req := http.Request{
		URL: url,
		Header: http.Header{
			// "Authorization": {"Bearer " + cli.conf.accesstoken},
			"Content-Type": {contentType},
			"Accept":       {"application/json"},
		},
		Method: http.MethodPost,
		Body:   io.NopCloser(bodyBuf),
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

	if response.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("[%d]%s", response.StatusCode, string(resBytes)))
		return
	}

	var res ResponseCar200
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		err = errors.New("json unmarshal response error").
			With(errors.Inner(err), errors.Playload(errors.MapData{"response": string(resBytes)}))
		return
	}

	return res.CID, nil
}
