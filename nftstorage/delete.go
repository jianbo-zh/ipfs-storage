package nftstorage

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jianbo-zh/go-errors"
)

func (cli *client) Delete(ctx context.Context, cid string) (ok bool, err error) {

	url, _ := url.Parse(cli.conf.endpoint + "/" + cid)

	req := http.Request{
		URL: url,
		Header: http.Header{
			"Authorization": {"Bearer " + cli.conf.accesstoken},
			// "Content-Type":  {"application/car"},
			"Accept": {"application/json"},
		},
		Method: http.MethodDelete,
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

	var resp200 DeleteResponse200
	err = json.Unmarshal(resBytes, &resp200)
	if err != nil {
		err = errors.New("json unmarshal error").
			With(errors.Inner(err), errors.Playload(errors.MapData{"response": string(resBytes)}))
		return
	}

	if !resp200.OK {
		logging.Infow(
			"delete cid",
			"cid", cid,
			"response", string(resBytes),
		)
	}

	return resp200.OK, nil
}
