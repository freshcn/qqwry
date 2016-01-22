package main

import (
	"net/http"
	"github.com/pquerna/ffjson/ffjson"
	"fmt"
)

func NewResponse(w http.ResponseWriter, r *http.Request) response {
	r.ParseForm()
	return response{
		w: w,
		r: r,
	}
}

// 返回正确的信息
func (this *response) ReturnSuccess(data interface{}) {
	this.Return(data, 200)
}

// 返回错误信息
func (this *response) ReturnError(statuscode, code int, errMsg string) {
	this.Return(map[string]interface{}{"errcode":code, "errmsg":errMsg}, statuscode)
}

// 向客户返回回数据
func (this *response) Return(data interface{}, code int) {
	jsonp := this.IsJSONP()

	rs, err := ffjson.Marshal(data)
	if  err != nil {
		code = 500
		rs = []byte(fmt.Sprintf(`{"errcode":500, "errmsg":"%s"}`, err.Error()))
	}

	this.w.WriteHeader(code)
	if jsonp == "" {
		this.w.Header().Add("Content-Type", "application/json")
		this.w.Write(rs)
	} else {
		this.w.Header().Add("Content-Type", "application/javascript")
		this.w.Write([]byte(fmt.Sprintf(`%s(%s)`, jsonp, rs)))
	}
}

// 是否为jsonp 请求
func (this *response) IsJSONP() string {
	if this.r.Form.Get("callback") != "" {
		return this.r.Form.Get("callback")
	}
	return ""
}
