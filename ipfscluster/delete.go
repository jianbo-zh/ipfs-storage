package ipfscluster

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/jianbo-zh/go-errors"
)

func (cli *client) Delete(ctx context.Context, cid string) (ok bool, err error) {

	url, _ := url.Parse(cli.conf.endpoint + "/pins/" + cid)

	req := http.Request{
		URL: url,
		Header: http.Header{
			"Authorization": {"Bearer " + cli.conf.accesstoken},
			"Accept":        {"application/json"},
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

	if response.StatusCode != 200 {
		err = errors.New(fmt.Sprintf("[%d]%s", response.StatusCode, string(resBytes)))
		return
	}

	var res DeleteResponse200
	err = json.Unmarshal(resBytes, &res)
	if err != nil {
		err = errors.New("json unmarshal response error").
			With(errors.Inner(err), errors.Playload(errors.MapData{"response": string(resBytes)}))
		return
	}

	return true, nil
}
