# simple-file-server

简单的 HTTP 文件服务器，界面类似于 Nginx 的 autoindex

人总是有一些奇奇怪怪的需求，比如偶尔要在两个云环境中传输东西。一直想找一个简单又方便的工具总是找不到。不是功能太杂，就是前戏配置麻烦，再要不就属于长久性的，不符合我们临时用用的精神理念

于是就有了这个东西，即下即开即用

## 快速上手

下载编译好的程序，直接运行即可

```sh
./sfs
```

默认会监听 `8080` 端口，并将程序当前目录作为文件服务器根目录

## 使用说明

命令行可以接收的参数

参数|默认值|描述
-|-|-
-port|8080|文件服务器监听的端口
-path|./|文件服务器服务根目录对应的本地路径

## 许可证

MIT © ZShab Niba