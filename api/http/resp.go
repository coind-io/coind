package httpapi

import (
	"encoding/json"
	"io"
)

type Resp struct {
	dict map[string]interface{}
}

func NewResp() *Resp {
	resp := new(Resp)
	resp.dict = map[string]interface{}{}
	resp.dict["retcode"] = "000000"
	resp.dict["retmsg"] = "success"
	return resp
}

func NewErrResp(err error) *Resp {
	e, ok := err.(ErrMsg)
	if ok == true {
		resp := NewResp()
		resp.dict["retcode"] = e.errcode
		resp.dict["retmsg"] = e.errmsg
		return resp
	} else {
		resp := NewResp()
		resp.dict["retcode"] = "000001"
		resp.dict["retmsg"] = err.Error()
		return resp
	}
}

func (resp *Resp) Put(key string, val interface{}) *Resp {
	resp.dict[key] = val
	return resp
}

func (resp *Resp) Encode(w io.Writer) error {
	return json.NewEncoder(w).Encode(resp.dict)
}
