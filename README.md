# 本项目展示了一个用go搭建的能执行lua脚本的web。

* 可以扩展成战斗逻辑用lua编写的战斗验证服务器。事实上我就是这么用的，但鉴于是公司项目，不能把所有的都开源出来，这里只提供golang运行多个lua虚拟机，充分利用多核性能的这么一个实现。

## build
* 在static目录下执行 make battle

* **windows下编译battle**因为用了c代码，编译需要gcc和make，windows下需要安装mingw-64x或者tdm64-gcc(需要把tdm64-gcc/bin/mingw32-make.exe改成make.exe)。

* 因为用了cgo所以不支持交叉编译


## lua库
* luajit_c:luajit的c源码，目前版本是2.0beta1

## 6. web:HTTP入口（仅供调试用）
* 如果运行参数web不为0，才开启http，可以浏览器打开http://127.0.0.1:web




