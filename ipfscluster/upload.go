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

	"github.com/jianbo-zh/go-errors"
	ipfsstorage "github.com/jianbo-zh/ipfs-storage"
)

func (cli *client) Upload(ctx context.Context, file ipfsstorage.UploadParam) (cid string, err error) {

	if file.Size > MAX_REQUEST_BODY_SIZE {
		err = ErrRequestBodyLimit
		return
	}

	url, _ := url.Parse(cli.conf.endpoint + "/add?cid-version=1")

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", file.Name)
	if err != nil {
		return cid, errors.New("create form file error").With(errors.Inner(err))
	}

	_, err = io.Copy(fileWriter, file.IOReader)
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

	var res Response200
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		err = errors.New("json unmarshal response error").
			With(errors.Inner(err), errors.Playload(errors.MapData{"response": string(resBytes)}))
		return
	}

	fmt.Printf("%s", string(resBytes))

	return res.CID, nil
}
