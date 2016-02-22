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
func (r *response) ReturnSuccess(data interface{}) {
	r.Return(data, 200)
}

// 返回错误信息
func (r *response) ReturnError(statuscode, code int, errMsg string) {
	r.Return(map[string]interface{}{"errcode":code, "errmsg":errMsg}, statuscode)
}

// 向客户返回回数据
func (r *response) Return(data interface{}, code int) {
	jsonp := r.IsJSONP()

	rs, err := ffjson.Marshal(data)
	if  err != nil {
		code = 500
		rs = []byte(fmt.Sprintf(`{"errcode":500, "errmsg":"%s"}`, err.Error()))
	}

	r.w.WriteHeader(code)
	if jsonp == "" {
		r.w.Header().Add("Content-Type", "application/json")
		r.w.Write(rs)
	} else {
		r.w.Header().Add("Content-Type", "application/javascript")
		r.w.Write([]byte(fmt.Sprintf(`%s(%s)`, jsonp, rs)))
	}
}

// 是否为jsonp 请求
func (r *response) IsJSONP() string {
	if r.r.Form.Get("callback") != "" {
		return r.r.Form.Get("callback")
	}
	return ""
}
