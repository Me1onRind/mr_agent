package tools

type CallToolReqest struct {
	Tool   string         `json:"tool"`
	Params map[string]any `json:"params"`
}

type CallToolResponse struct {
	Result any `json:"result"`
}
