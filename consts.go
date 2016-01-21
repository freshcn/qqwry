package main

import "os"

const (
	INDEX_NUM = 7 // 索引长度
)

// IP 的数据信息
type qqwry struct {
	Ip      uint32 `json:"ip"`
	Country []byte `json:"country"`
	Area    []byte `json:"area"`
}

// ip库的索引
type index struct {
	Ip     uint32
	Offset uint32
}

type ipData struct {
	Index    []index
	Data     map[uint32]qqwry
	FilePath string
	Path     *os.File
}
