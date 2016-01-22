package main

import (
	"os"
	"net/http"
)

const (
	INDEX_LEN       = 7    // 索引长度
	REDIRECT_MODE_1 = 0x01 // 国家的类型, 指向另一个指向
	REDIRECT_MODE_2 = 0x02 // 国家的类型, 指向一个指向
)

type resultQQwry struct {
	Ip      string `json:"ip"`
	Country string `json:"country"`
	Area    string `json:"area"`
}

type fileData struct {
	Data []byte
	FilePath string
	Path     *os.File
	IpNum    int64
}

type QQwry struct {
	Data     *fileData
	Offset   int64
}

// 向客户端返回数据的
type response struct {
	r *http.Request
	w http.ResponseWriter
}
