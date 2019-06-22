# go-socks-proxy
go-socks-proxy是一个使用socks5协议的代理服务器

## 流量监控
go-socks-proxy集成了prometheus监控，实现了以下多个视角的流量监控
* 针对不同客户端的流量监控
* 针对不同的访问站点的流量监控
* 针对上传、下载总流量的监控



## 转发过程中出现的错误
![统计运行错误](https://github.com/kangaloo/go-socks-proxy/blob/ae01c36eb6e34c288606de87e9dfda5b199937fa/resource/failed_total.jpg)