package main

import (
	"flag"
	"log"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	datFile := flag.String("qqwry", "./qqwry.dat", "纯真 IP 库的地址")
	//	port := flag.String("port", "2060", "HTTP 请求监听端口号")
	flag.Parse()

	IpData.FilePath = *datFile

	startTime := time.Now().UnixNano()
	res := IpData.InitIpData()

	if v, ok := res.(error); ok {
		log.Panic(v)
	}
	endTime := time.Now().UnixNano()

	// 下面开始加载 http 相关的服务

	log.Printf("IP 库加载完成 共加载:%d 条 IP 记录, 所花时间:%.1f秒", len(IpData.Index), float64(endTime-startTime)/1000000000)


	log.Println(IpData.Find("117.174.30.54"))

	time.Sleep(1000*time.Second)
}
