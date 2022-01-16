package web3storage

type Response200 struct {
	CID string `json:"cid"`
}

type Response404 struct {
	Message string `json:"message"`
}

type StatusResponse200 struct {
	CID     string `json:"cid"`
	DagSize int64  `json:"dagSize"`
	Created string `json:"created"`
	Pins    []Pin  `json:"pins"`
	Deals   []Deal `json:"deals"`
}

type Pin struct {
	PeerID   string `json:"peerId"`
	PeerName string `json:"peerName"`
	Region   string `json:"region"`
	Status   string `json:"status"` // ["PinQueued", "Pinning", "Pinned"]
	Updated  string `json:"updated"`
}

type Deal struct {
	DealID            int64  `json:"dealId"`
	StorageProvider   string `json:"storageProvider"`
	Status            string `json:"status"` // ["Queued" "Published" "Active"]
	PieceCid          string `json:"pieceCid"`
	DataCid           string `json:"dataCid"`
	DataModelSelector string `json:"dataModelSelector"`
	Activation        string `json:"activation"`
	Created           string `json:"created"`
	Updated           string `json:"updated"`
}
