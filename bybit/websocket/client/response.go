package client

type PublicResponse struct {
	Topic string      `json:"topic"`
	Type  string      `json:"type"`
	TS    int64       `json:"ts"`
	Data  interface{} `json:"data"`
}

type TickersResponse struct {
	Topic string       `json:"topic"`
	Type  string       `json:"type"`
	TS    int64        `json:"ts"`
	Data  *TickersData `json:"data"`
}

type TickersData struct {
	Symbol            string `json:"symbol"`
	TickDirection     string `json:"tickDirection"`
	Price24HPcnt      string `json:"price24hPcnt"`
	LastPrice         string `json:"lastPrice"`
	PrevPrice24H      string `json:"prevPrice24h"`
	HighPrice24H      string `json:"highPrice24h"`
	LowPrice24H       string `json:"lowPrice24h"`
	PrevPrice1H       string `json:"prevPrice1h"`
	MarkPrice         string `json:"markPrice"`
	IndexPrice        string `json:"indexPrice"`
	OpenInterest      string `json:"openInterest"`
	OpenInterestValue string `json:"openInterestValue"`
	Turnover24H       string `json:"turnover24h"`
	Volume24H         string `json:"volume24h"`
	NextFundingTime   string `json:"nextFundingTime"`
	FundingRate       string `json:"fundingRate"`
	Bid1Price         string `json:"bid1Price"`
	Bid1Size          string `json:"bid1Size"`
	Ask1Price         string `json:"ask1Price"`
	Ask1Size          string `json:"ask1Size"`
}
