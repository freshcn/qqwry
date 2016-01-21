package main

import (
	"encoding/binary"
	"errors"
//	"github.com/axgle/mahonia"
	"io"
	"log"
	"net"
	"os"
)

var IpData ipData

// 初始化ip库数据到内存中
func (this *ipData) InitIpData() (rs interface{}) {

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

	this.Path.Seek(0, 0)
	indexPos := make([]byte, 8)
	if _, err := this.Path.Read(indexPos); err != nil {
		rs = err
		return
	}

	start := binary.LittleEndian.Uint32(indexPos[:4])
	end := binary.LittleEndian.Uint32(indexPos[4:])

	// 索引数量
	indexNum := (end-start)/INDEX_NUM + 1

	// 临时索引,当文件加载完成后将索引给到 ipData 提供给查询使用
	tmpIndex := make([]index, indexNum)
	// 临时文件数据,当文件加载完成后,将数据给到 ipData 提供给查询使用
	tmpData := make(map[uint32]qqwry)

	// index 的缓存数据
	indexBuf := make([]byte, INDEX_NUM)
	// data 的临时缓冲数据
	dataTmpBuf := make([]byte, 1)
	// 从数据区获取 IP 地址
	ipBuf := make([]byte, 4)
	// 数据的偏移量
	dataOffset := uint32(0)

	j := 0
	zeroNum := 0

	// gbk转utf8
//	enc := mahonia.NewEncoder("gbk")
	// 开始加载索引
	for i := uint32(0); i < indexNum; i++ {

		// 将文件的指针跳转到索引开始的位置
		this.Path.Seek(int64(start+i*INDEX_NUM), 0)
		if _, err := this.Path.Read(indexBuf); err != nil {
			if err == io.EOF {
				break
			}
			continue
		}

		dataOffset = byteToUInt32(indexBuf[4:])

		tmpIndex[j] = index{
			Ip:     binary.LittleEndian.Uint32(indexBuf[:4]),
			Offset: dataOffset,
		}

		j++

		// 开始获取 IP 的地址数据
//		this.Path.Seek(int64(dataOffset), 0)
		this.Path.Seek(int64(dataOffset + 4), 0)

//		if _, err = this.Path.Read(ipBuf); err != nil {
//			continue
//		}

		tmpQQwry := qqwry{
//			Ip:      binary.LittleEndian.Uint32(ipBuf),
			Country: make([]byte, 0, 50),
			Area:    make([]byte, 0, 50),
		}

		zeroNum = 0

		// 先判断是国家是什么类型
		mode := this.getMode()

		switch mode {
		case COUNTRY_MODE_1:// 模式1 国家和区域都走了

		case COUNTRY_MODE_2: // 模板2 国家走了
		default:

		}

		for i := 0; i < 1024; i++ {

			if zeroNum > 1 {
				zeroNum = 0
				break
			}

			if _, err = this.Path.Read(dataTmpBuf); err != nil {
				continue
			}

			if dataTmpBuf[0] == 0 {
				zeroNum++
				continue
			}

			if zeroNum == 0 {
				tmpQQwry.Country = append(tmpQQwry.Country, dataTmpBuf[0])
			} else if zeroNum == 1 {
				tmpQQwry.Area = append(tmpQQwry.Area, dataTmpBuf[0])
			}

//			// 当国家为字符串时，需要对其进行转码为utf-8
//			if tmpQQwry.Country[0] != COUNTRY_MODE_1 && tmpQQwry.Country[0] != COUNTRY_MODE_2 {
//				tmpQQwry.Country = []byte(enc.ConvertString(string(tmpQQwry.Country)))
//			}
//
//			// 当区域的数据不为空时
//			if len(tmpQQwry.Area) > 0 {
//				tmpQQwry.Area = []byte(enc.ConvertString(string(tmpQQwry.Area)))
//			}

		}

//		log.Printf("ip:%d, country:%s, area:%s \n", tmpQQwry.Ip, tmpQQwry.Country, tmpQQwry.Area);
		tmpData[dataOffset] = tmpQQwry

//		if i == 0 {
//			os.Exit(0)
//		}
	}

	this.Index = tmpIndex
	this.Data = tmpData

	return
}

// 从文件中判断国家的类型
func (this *ipData) getMode(offset ...int) byte {
	buf := make([]byte, 1)
	if len(offset) > 0 {
		this.Path.Seek(offset[0], 0)
	}

	if _, err := this.Path.Read(buf); err == nil {
		return buf[0]
	}

}

// 查询数据
func (this *ipData) Find(ip string) interface{} {
	userIp := binary.BigEndian.Uint32(net.ParseIP(ip).To4())

	log.Println("用户输入的ip地址转：", userIp)

	start := 0
	end := len(this.Index)

	for {
		log.Println("start:", start, "end:", end)
		mid := this.FindMiddle(start, end)
		log.Println("中间点：", mid)

		offset := this.Index[mid].Offset
		log.Println("中间点的数据", this.Index[mid])

		if _, e := this.Data[offset]; e {

			if end-start == 1 {
				res := resultQQwry{
					Ip: ip,
				}
				res.Country, res.Area = this.ReadCountryAndArea(this.Data[offset])
				return res
			}

			log.Printf("ip记录：%#v\n", this.Data[offset])
			if this.Data[offset].Ip == userIp {
				log.Println("ip地址正好对上===============================")
				res := resultQQwry{
					Ip: ip,
				}

				res.Country, res.Area = this.ReadCountryAndArea(this.Data[offset])
				return res
			} else if this.Data[offset].Ip > userIp {
				end = mid
			} else if this.Data[offset].Ip < userIp {
				start = mid
			} else {
				log.Println("2个ip之间完全没有关系", this.Data[offset].Ip, " - ", userIp)
			}
		} else {
			log.Println("没有找到中间点的IP记录")
			return false
		}
	}
}

// 获取国家和地区数据
func (this *ipData) ReadCountryAndArea(data qqwry) (country, area string) {
	log.Printf("显示这个数据%s\n", data)
	switch data.Country[0] {
	case COUNTRY_MODE_1: // 模式1,地址和国家都走了
		countryOffset := byteToUInt32(data.Country[1:4])

		tmpData := this.Data[countryOffset]
		log.Println("国家的偏移量", tmpData)
		if tmpData.Country[0] == COUNTRY_MODE_2 {
			country = string(this.Data[byteToUInt32(tmpData.Country[1:4])].Country)
		} else {
			country = string(tmpData.Country)
		}
		area = string(tmpData.Area)
	case COUNTRY_MODE_2: // 模式2,国家走了
		area = string(data.Area)
		countryOffset := byteToUInt32(data.Country[1:4])
		country = string(this.Data[countryOffset].Country)
	default:
		area = string(data.Area)
		country = string(data.Country)
	}

	//	enc := mahonia.NewEncoder("gbk")
	//	area = enc.ConvertString(area)
	//	country = enc.ConvertString(country)
	return
}

// 查找中间位置
func (this *ipData) FindMiddle(start, end int) int {
	return start + ((end - start) >> 1)
}

// 将 byte 转换为uint32
func byteToUInt32(data []byte) uint32 {
	i := uint32(data[0]) & 0xff
	i |= (uint32(data[1]) << 8) & 0xff00
	i |= (uint32(data[2]) << 16) & 0xff0000
	return i
}
