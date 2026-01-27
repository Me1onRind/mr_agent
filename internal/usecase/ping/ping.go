package ping

import (
	"context"

	"github.com/Me1onRind/mr_agent/protocol/http/ping"
)

type PingUsecase struct {
}

func NewPingUsecase() *PingUsecase {
	return &PingUsecase{}
}

func (p *PingUsecase) Echo(ctx context.Context, request *ping.EchoRequest) (*ping.EchoResponse, error) {
	resp := &ping.EchoResponse{
		Msg: request.Msg,
	}
	return resp, nil
}
