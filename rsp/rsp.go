package rsp

import (
	"fmt"

	"github.com/virzz/ginx/code"
)

type Rsp struct {
	code.APICode
	Data any `json:"data"`
}

func (r *Rsp) WithCode(c code.APICode) *Rsp {
	r.APICode = c
	return r
}
func (r *Rsp) WithData(v any) *Rsp {
	r.Data = v
	return r
}
func (r *Rsp) WithMsg(v string) *Rsp {
	r.Msg = v
	return r
}
func (r *Rsp) WithItem(total int64, items any) *Rsp {
	r.Data = &Items{Items: items, Total: total}
	return r
}
func (r *Rsp) WithItemExt(total int64, items any, ext any) *Rsp {
	r.Data = &Items{Items: items, Total: total, Ext: ext}
	return r
}

func New() *Rsp                    { return &Rsp{} }
func C(c code.APICode) *Rsp        { return &Rsp{APICode: c} }
func S(data any) *Rsp              { return &Rsp{APICode: code.Success, Data: data} }
func M(msg string) *Rsp            { return (&Rsp{APICode: code.Success}).WithMsg(msg) }
func SM(data any, msg string) *Rsp { return (&Rsp{APICode: code.Success, Data: data}).WithMsg(msg) }

func OK() *Rsp            { return &Rsp{APICode: code.Success} }
func UnImplemented() *Rsp { return &Rsp{APICode: code.UnImplemented} }

func E(c code.APICode, msg any) *Rsp {
	switch m := msg.(type) {
	case string:
		return (&Rsp{APICode: c}).WithMsg(m)
	case error:
		return (&Rsp{APICode: c}).WithMsg(m.Error())
	default:
		return (&Rsp{APICode: c}).WithMsg(fmt.Sprintf("%v", msg))
	}
}
