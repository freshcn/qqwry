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
func (f *fileData) InitIpData() (rs interface{}) {

	// 判断文件是否存在
	_, err := os.Stat(f.FilePath)
	if err != nil && os.IsNotExist(err) {
		rs = errors.New("文件不存在")
		return
	}

	// 打开文件句柄
	f.Path, err = os.OpenFile(f.FilePath, os.O_RDONLY, 0400)
	if err != nil {
		rs = err
		return
	}
	defer f.Path.Close()


	tmpData, err := ioutil.ReadAll(f.Path)
	if err != nil {
		log.Println(err)
		rs = err
		return
	}

	f.Data = tmpData

	buf := f.Data[0:8]
	start := binary.LittleEndian.Uint32(buf[:4])
	end := binary.LittleEndian.Uint32(buf[4:])

	f.IpNum = int64((end-start)/INDEX_LEN + 1)

	return true
}

// 新建 qqwry  类型
func NewQQwry() QQwry {
	return QQwry{
		Data: &IpData,
	}
}

// 从文件中读取数据
func (q *QQwry) ReadData(num int, offset ...int64) (rs []byte) {
	if len(offset) > 0 {
		q.SetOffset(offset[0])
	}
	nums := int64(num)
	end := q.Offset+nums
	dataNum := int64(len(q.Data.Data))
	if (q.Offset > dataNum) {
		return nil
	}

	if (end > dataNum) {
		end = dataNum
	}
	rs = q.Data.Data[q.Offset : end]
	q.Offset = end
	return
}

// 设置偏移量
func (q *QQwry) SetOffset(offset int64) {
	q.Offset = offset
}

func (q *QQwry) Find(ip string) (res resultQQwry) {

	res = resultQQwry{}

	res.Ip = ip
	if strings.Count(ip, ".") != 3 {
		return res
	}
	offset := q.searchIndex(binary.BigEndian.Uint32(net.ParseIP(ip).To4()))
	if offset <= 0 {
		return
	}

	var country []byte
	var area []byte

	mode := q.readMode(offset + 4)
	if mode == REDIRECT_MODE_1 {
		countryOffset := q.readUInt24()
		mode = q.readMode(countryOffset)
		if mode == REDIRECT_MODE_2 {
			c := q.readUInt24()
			country = q.readString(c)
			countryOffset += 4
		} else {
			country = q.readString(countryOffset)
			countryOffset += uint32(len(country) + 1)
		}
		area = q.readArea(countryOffset)
	} else if mode == REDIRECT_MODE_2 {
		countryOffset := q.readUInt24()
		country = q.readString(countryOffset)
		area = q.readArea(offset + 8)
	} else {
		country = q.readString(offset + 4)
		area = q.readArea(offset + uint32(5+len(country)))
	}


	enc := mahonia.NewDecoder("gbk")
	res.Country = enc.ConvertString(string(country))
	res.Area = enc.ConvertString(string(area))

	return
}

func (q *QQwry) readMode(offset uint32) byte {
	mode := q.ReadData(1, int64(offset))
	return mode[0]
}

func (q *QQwry) readArea(offset uint32) []byte {
	mode := q.readMode(offset)
	if mode == REDIRECT_MODE_1 || mode == REDIRECT_MODE_2 {
		areaOffset := q.readUInt24()
		if areaOffset == 0 {
			return []byte("")
		} else {
			return q.readString(areaOffset)
		}
	} else {
		return q.readString(offset)
	}
	return []byte("")
}

func (q *QQwry) readString(offset uint32) []byte {
	q.SetOffset(int64(offset))
	data := make([]byte, 0, 30)
	buf := make([]byte, 1)
	for {
		buf = q.ReadData(1)
		if buf[0] == 0 {
			break
		}
		data = append(data, buf[0])
	}
	return data
}

func (q *QQwry) searchIndex(ip uint32) uint32 {
	header := q.ReadData(8, 0)

	start := binary.LittleEndian.Uint32(header[:4])
	end := binary.LittleEndian.Uint32(header[4:])

	buf := make([]byte, INDEX_LEN)
	mid := uint32(0)
	_ip := uint32(0)

	for {
		mid = q.getMiddleOffset(start, end)
		buf = q.ReadData(INDEX_LEN, int64(mid))
		_ip = binary.LittleEndian.Uint32(buf[:4])

		if end-start == INDEX_LEN {
			offset := byteToUInt32(buf[4:])
			buf = q.ReadData(INDEX_LEN)
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

func (q *QQwry) readUInt24() uint32 {
	buf := q.ReadData(3)
	return byteToUInt32(buf)
}

func (q *QQwry) getMiddleOffset(start uint32, end uint32) uint32 {
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
