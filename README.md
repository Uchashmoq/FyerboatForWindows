# FyerboatForWindows
一个基于socks5的加密代理,GUI采用fyne框架

[服务端](https://github.com/Uchashmoq/TunProxyServer)采用java的netty框架开发

[客户端](https://github.com/Uchashmoq/TunProxyClient)命令行版，可编译成Windows，Linux和Mac可执行文件(这个也可以，但是安装fyne要折腾一下)

## 预览

![主菜单](https://github.com/Uchashmoq/FyerboatForWindows/blob/main/img/menu.png)

![我的节点](https://github.com/Uchashmoq/FyerboatForWindows/blob/main/img/nodes.png)

![日志](https://github.com/Uchashmoq/FyerboatForWindows/blob/main/img/log.png)

![添加节点](https://github.com/Uchashmoq/FyerboatForWindows/blob/main/img/add.png)

```
添加节点：
名称，如 MyNode
IP:端口，如 192.168.1.88:14445
密钥，16位英文数字以及特殊符号，与服务端staticKey一致,如 0123456789abcdef
```



## 编译

```go
go build
```

## 运行

点击可执行文件，用支持socks5代理的浏览器插件(如switchyomega)或其他代理软件将流量转发到本软件即可，详情看[教程.doc](https://github.com/Uchashmoq/FyerboatForWindows/blob/main/教程.doc)

## 注意事项

1.msyh.ttc是字体文件，删除后会导致中文乱码

2.config.json保存有节点信息，监听端口，请勿删除
