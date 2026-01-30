package cli

import "github.com/alecthomas/kong"

type Params struct {
	Mode DialogMode `help:"dialog mode 1.without context 2.with context" default:"1" enum:"1,2"`
}

func parseCliParams() *Params {
	params := &Params{}
	_ = kong.Parse(params)
	return params
}

func (p *Params) WithContext() bool {
	return p.Mode == WithContext
}
