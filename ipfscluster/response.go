package ipfscluster

import "github.com/ipfs/go-cid"

type Response200 struct {
	Name string  `json:"name"`
	CID  cid.Cid `json:"cid"`
	Size int64   `json:"size"`
}

type ResponseCar200 struct {
	Name  string  `json:"name"`
	CID   cid.Cid `json:"cid"`
	Size  int64   `json:"size"`
	Bytes int64   `json:"bytes"` // chunk size
}

type ResponseCid200 struct {
	Name string  `json:"name"`
	CID  cid.Cid `json:"cid"`
	// many other fields ignore
}

type DeleteResponse200 struct {
	Name string  `json:"name"`
	CID  cid.Cid `json:"cid"`
	// many other fields ignore
}

type StatusResponse200 struct {
	CID     cid.Cid                  `json:"cid"`
	Name    string                   `json:"name"`
	PeerMap map[string]PeerPinStatus `json:"peer_map"`
}

type PeerPinStatus struct {
	PeerName     string `json:"peername"`
	Status       string `json:"status"`
	Timestamp    string `json:"timestamp"`
	Error        string `json:"error"`
	AttemptCount int    `json:"attempt_count"`
	PriorityPin  bool   `json:"priority_pin"`
}
