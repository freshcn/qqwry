package main

import (
	"os"
	"fmt"
)

const (
	INDEX_NUM      = 7    // 索引长度
	COUNTRY_MODE_1 = 0x01 // 国家的类型, 指向另一个指向
	COUNTRY_MODE_2 = 0x02 // 国家的类型, 指向一个指向
)

type resultQQwry struct {
	Ip      string `json:"ip"`
	Country string `json:"country"`
	Area    string `json:"area"`
}

// IP 的数据信息
type qqwry struct {
//	Ip      uint32
	Country []byte
	Area    []byte
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

func (this qqwry) String() string {
	return fmt.Sprintf("ip:%d, country:%s, area:%s", this.Ip, this.Country, this.Area)
}

func (this index) String() string {
	return fmt.Sprintf("ip:%d, offset:%d", this.Ip, this.Offset)
}