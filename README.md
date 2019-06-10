# gox
[![go report card](https://goreportcard.com/badge/github.com/hetianyi/gox)](https://goreportcard.com/report/github.com/hetianyi/gox)

gox is a go library integrated with many useful tools.



## Packages and functions

- cache

提供对象/字节数组缓存功能

- conn

提供连接池管理功能，以及连接自动回收。

- ~~consulx~~

- convert

提供丰富的数据类型转换功能

- file

提供常用的文件操作功能

- gpip

提供一种基于tcp协议的server/client数据通讯方法。

- httpx

mock实现了简单的链式http请求构造器，可以用一行代码搞定http请求。

- img

基于```github.com/disintegration/imaging```，封装了常用的图片处理接口。

- gifx

基于img，能够处理GIF生成，GIF加水印等操作。

- logger

基于```github.com/sirupsen/logrus, github.com/logrusorgru/aurora```，能够自定义日志格式，以及输出彩色日志。

- ~~pluginx~~

- pool

提供任务池功能，能够设置并行，等待任务数。

- timer

简单的日期工具包

- uuid

基于```github.com/satori/go.uuid```，能够生成UUID。

- ws

基于```github.com/gorilla/websocket```，实现处理websocket。

