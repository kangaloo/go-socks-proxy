# go-socks-proxy
go-socks-proxy是一个使用socks5协议的代理服务器

## 安装
1. go get github.com/kangaloo/go-socks-proxy
2. cd go-socks-proxy
3. go build
4. ./go-socks-proxy

## 流量监控
go-socks-proxy集成了prometheus监控，实现了以下多个视角的流量监控
* 针对不同客户端的流量监控
* 针对不同的访问站点的流量监控
* 针对上传、下载总流量的监控

## 根据不同站点区分的上传流量
![统计运行错误](https://github.com/kangaloo/go-socks-proxy/blob/f0fd303049b28bf87680aa28f4ac8923efb627ce/resource/site-upload.jpg)

## 根据不同站点区分的下载流量
![统计运行错误](https://github.com/kangaloo/go-socks-proxy/blob/f0fd303049b28bf87680aa28f4ac8923efb627ce/resource/site-download.jpg)

## 上传总流量和下载总流量
![统计运行错误](https://github.com/kangaloo/go-socks-proxy/blob/f0fd303049b28bf87680aa28f4ac8923efb627ce/resource/upload-download-total.jpg)

## 转发过程中出现的错误
![统计运行错误](https://github.com/kangaloo/go-socks-proxy/blob/ae01c36eb6e34c288606de87e9dfda5b199937fa/resource/failed_total.jpg)