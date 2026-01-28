package ping

type EchoRequest struct {
	Msg string `json:"msg"`
}

type EchoResponse struct {
	Msg string `json:"msg"`
}

type HelloToAgentRequest struct {
	Msg string `json:"msg"`
}

type HelloToAgentResponse struct {
	Msg string `json:"msg"`
}
