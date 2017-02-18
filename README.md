# qqwry 

![travis ci](https://travis-ci.org/freshcn/qqwry.svg?branch=master)

纯真 IP 库的一个服务。通过http提供一个ip地址归属地查询支持

## 软件介绍

我们大家做网站的时候，都会需要将用户的IP地址转换为归属地址功能，而之前的作法大都是从硬盘的数据文件中读取，这样不太高效。我这次将纯真IP库直接加载到内存中，并以http服务的方式提供接口出来。这样方便自己项目中所有的程序都能方便的接口IP查询功能了。

## 安装

### go安装

```
go get github.com/freshcn/qqwry
```
### 二进制包直接下载

https://github.com/freshcn/qqwry/releases

### 下载纯真IP库
请访问 http://www.cz88.net 下载纯真IP库，需要在windows中安装程序，然后在程序的安装目录可以找到qqwry.dat文件，复制出来放到和本程序同一个目录（当然也可是其他目录，只是需要在运行的时候指定IP库目录），

### 运行参数

运行 ./qqwry -h 可以看到本服务程序的可用运行参数

```
  -port string
    	HTTP 请求监听端口号 (default "2060")
  -qqwry string
    	纯真 IP 库的地址 (default "./qqwry.dat")
```

## 使用方法
```
http://127.0.0.1:2060?ip=8.8.8.8,114.114.114.114&callback=a
```

* ip - 要查询的IP地址，可使用半角逗号分隔查询多个IP地址。必填项
* callback - jsonp回调函数名，当提交了这个参数，将会按jsonp格式返回。非必填

** 返回结果 **

```json
{"114.114.114.114":{"ip":"114.114.114.114","country":"江苏省南京市","area":"南京信风网络科技有限公司GreatbitDNS服务器"},"8.8.8.8":{"ip":"8.8.8.8","country":"美国","area":"加利福尼亚州圣克拉拉县山景市谷歌公司DNS服务器"}}
```
* ip - 输入的ip地址
* country - 国家或地区
* area - 区域（我实际测试得到还有可能是运营商）


### 感谢

* 感谢[纯真IP库](http://www.cz88.net)一直以来坚持为大家提供免费的IP库资源
* 感谢[yinheli](https://github.com/yinheli)的[qqwry](https://github.com/yinheli/qqwry)项目，为我提供了纯真ip库文件格式算法
