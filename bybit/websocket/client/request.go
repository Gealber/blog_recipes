package client

type Request struct {
	ReqID string        `json:"req_id,omitempty"`
	Op    string        `json:"op,omitempty"`
	Args  []interface{} `json:"args,omitempty"`
}