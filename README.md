[![GoDoc](https://godoc.org/github.com/bailu1901/goweblua?status.png)](https://godoc.org/github.com/bailu1901/goweblua)
[![Build Status](https://travis-ci.org/bailu1901/goweblua.svg?branch=master)](https://travis-ci.org/bailu1901/goweblua)
[![Go Report](https://goreportcard.com/badge/github.com/bailu1901/goweblua)](https://goreportcard.com/report/github.com/bailu1901/goweblua)


# 简介
	本项目是用go语言编写，结合cgo功能，支持高并发执行lua脚本的程序。

## 扩展
* 可以扩展成逻辑用lua编写的战斗验证服务器。这里展示了golang运行多个lua虚拟机，充分利用多核性能的个实现。
* 可以扩展用grpc做外部的可负载均衡的接口，我这里只简单的实现了用http做外部接口

## 编译
* 输入命令 make
* **windows下编译battle**因为用了c代码，编译需要gcc和make，windows下需要安装mingw-64x或者tdm64-gcc(需要把tdm64-gcc/bin/mingw32-make.exe改成make.exe)。
* 因为用了cgo所以不支持交叉编译

## 运行
```
CGO_ENABLED=1 go run ./cmd
```

## lua库
* 我用的是LuaJIT-2.1.0-beta1，可自行替换lua的c代码，但要修改c/Makefile

## Web接口
* 浏览器提供了一个简单web页面 http://127.0.0.1
* 热加载lua脚本 http://127.0.0.1?cmd=reload



