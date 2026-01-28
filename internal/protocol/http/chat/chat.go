package chat

type ChatRequest struct {
	Msg string `json:"msg"`
}

type ChatResponse struct {
	Msg string `json:"msg"`
}
