package main

import (
	"encoding/binary"
	"errors"
	"github.com/axgle/mahonia"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

var IpData fileData

// 初始化ip库数据到内存中
func (this *fileData) InitIpData() (rs interface{}) {

	// 判断文件是否存在
	_, err := os.Stat(this.FilePath)
	if err != nil && os.IsNotExist(err) {
		rs = errors.New("文件不存在")
		return
	}

	// 打开文件句柄
	this.Path, err = os.OpenFile(this.FilePath, os.O_RDONLY, 0400)
	if err != nil {
		rs = err
		return
	}
	defer this.Path.Close()


	tmpData, err := ioutil.ReadAll(this.Path)
	if err != nil {
		log.Println(err)
		rs = err
		return
	}

	this.Data = tmpData

	buf := this.Data[0:8]
	start := binary.LittleEndian.Uint32(buf[:4])
	end := binary.LittleEndian.Uint32(buf[4:])

	this.IpNum = int64((end-start)/INDEX_LEN + 1)

	return true
}

// 新建 qqwry  类型
func NewQQwry() QQwry {
	return QQwry{
		Data: &IpData,
	}
}

// 从文件中读取数据
func (this *QQwry) ReadData(num int, offset ...int64) (rs []byte) {
	if len(offset) > 0 {
		this.SetOffset(offset[0])
	}
	nums := int64(num)
	end := this.Offset+nums
	dataNum := int64(len(this.Data.Data))
	if (this.Offset > dataNum) {
		return nil
	}

	if (end > dataNum) {
		end = dataNum
	}
	rs = this.Data.Data[this.Offset : end]
	this.Offset = end
	return
}

// 设置偏移量
func (this *QQwry) SetOffset(offset int64) {
	this.Offset = offset
}

func (this *QQwry) Find(ip string) (res resultQQwry) {

	res = resultQQwry{}

	res.Ip = ip
	if strings.Count(ip, ".") != 3 {
		return res
	}
	offset := this.searchIndex(binary.BigEndian.Uint32(net.ParseIP(ip).To4()))
	if offset <= 0 {
		return
	}

	var country []byte
	var area []byte

	mode := this.readMode(offset + 4)
	if mode == REDIRECT_MODE_1 {
		countryOffset := this.readUInt24()
		mode = this.readMode(countryOffset)
		if mode == REDIRECT_MODE_2 {
			c := this.readUInt24()
			country = this.readString(c)
			countryOffset += 4
		} else {
			country = this.readString(countryOffset)
			countryOffset += uint32(len(country) + 1)
		}
		area = this.readArea(countryOffset)
	} else if mode == REDIRECT_MODE_2 {
		countryOffset := this.readUInt24()
		country = this.readString(countryOffset)
		area = this.readArea(offset + 8)
	} else {
		country = this.readString(offset + 4)
		area = this.readArea(offset + uint32(5+len(country)))
	}


	enc := mahonia.NewDecoder("gbk")
	res.Country = enc.ConvertString(string(country))
	res.Area = enc.ConvertString(string(area))

	return
}

func (this *QQwry) readMode(offset uint32) byte {
	mode := this.ReadData(1, int64(offset))
	return mode[0]
}

func (this *QQwry) readArea(offset uint32) []byte {
	mode := this.readMode(offset)
	if mode == REDIRECT_MODE_1 || mode == REDIRECT_MODE_2 {
		areaOffset := this.readUInt24()
		if areaOffset == 0 {
			return []byte("")
		} else {
			return this.readString(areaOffset)
		}
	} else {
		return this.readString(offset)
	}
	return []byte("")
}

func (this *QQwry) readString(offset uint32) []byte {
	this.SetOffset(int64(offset))
	data := make([]byte, 0, 30)
	buf := make([]byte, 1)
	for {
		buf = this.ReadData(1)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

func (this *QQwry) searchIndex(ip uint32) uint32 {
	header := this.ReadData(8, 0)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, INDEX_LEN)
	mid := uint32(0)
	_ip := uint32(0)

	for {
		mid = this.getMiddleOffset(start, end)
		buf = this.ReadData(INDEX_LEN, int64(mid))
		_ip = binary.LittleEndian.Uint32(buf[:4])

		if end-start == INDEX_LEN {
			offset := byteToUInt32(buf[4:])
			buf = this.ReadData(INDEX_LEN)
			if ip < binary.LittleEndian.Uint32(buf[:4]) {
				return offset
			} else {
				return 0
			}
		}

		// 找到的比较大，向前移
		if _ip > ip {
			end = mid
		} else if _ip < ip { // 找到的比较小，向后移
			start = mid
		} else if _ip == ip {
			return byteToUInt32(buf[4:])
		}

	}
	return 0
}

func (this *QQwry) readUInt24() uint32 {
	buf := this.ReadData(3)
	return byteToUInt32(buf)
}

func (this *QQwry) getMiddleOffset(start uint32, end uint32) uint32 {
	records := ((end - start) / INDEX_LEN) >> 1
	return start + records*INDEX_LEN
}

// 将 byte 转换为uint32
func byteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
