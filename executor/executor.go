// Package executor 是lua并发执行器
package executor

import (
	"github.com/bailu1901/goweblua/luajit"

	"errors"

	l4g "github.com/alecthomas/log4go"
)

type luaWorker struct {
	lua      *luajit.State
	funcName string
}

//printstack 打印lua调用栈
func printstack(s *luajit.State) {
	l4g.Error("--- stack:")
	n := s.Gettop()
	for i := 1; i <= n; i++ {
		switch s.Type(i) {
		case luajit.Tstring:
			l4g.Error("%s", s.Tostring(i))
		case luajit.Tnumber:
			l4g.Error("%f", s.Tonumber(i))
		case luajit.Tboolean:
			l4g.Error("%t", s.Toboolean(i))
		default:
			l4g.Error("(%s)", s.Typename(s.Type(i)))
		}
		l4g.Error("\n")
	}
	l4g.Error("---")
}

//NewLuaWorker 构造器
func NewLuaWorker(entryFile string, funcName string) (*luaWorker, error) {
	worker := &luaWorker{
		lua:      luajit.Newstate(),
		funcName: funcName,
	}

	worker.lua.Openlibs()
	worker.lua.Openpbclibs()
	worker.lua.Openlpeglibs()
	worker.lua.Opencjsonlibs()

	if "" != entryFile {
		if err := worker.lua.Loadfile(entryFile); nil != err {
			printstack(worker.lua)
			return nil, err
		}
	}

	return worker, nil
}

//Init 初始化
func (worker *luaWorker) Init() error {
	//执行已经加载的脚本
	if err := worker.lua.Pcall(0, 0, 0); err != nil {
		printstack(worker.lua)
		return err
	}

	//检查复盘函数是否为函数
	worker.lua.Getglobal(worker.funcName)
	if !worker.lua.Isfunction(-1) {
		printstack(worker.lua)
		return errors.New("worker.funcName is not a function")
	}

	worker.lua.Settop(0)

	return nil
}

//Execute 执行 groutine不安全
func (worker *luaWorker) Execute(data []byte) ([]byte, error) {

	//恢复lua栈
	worker.lua.Settop(0)

	worker.lua.Getglobal(worker.funcName)

	worker.lua.Pushlstring(string(data), len(data))
	if err := worker.lua.Pcall(1, 1, 0); err != nil {
		printstack(worker.lua)
		l4g.Error("[worker lua] Execute Pcall error=%s", err)
		return []byte(""), err
	}

	ret, _ := worker.lua.Tolstring(-1)

	return []byte(ret), nil
}

// Destroy 销毁lua对象
func (worker *luaWorker) Destroy() {

	// 不允许失败
	defer func() {
		if e := recover(); nil != e {
			l4g.Error("[worker lua] error=%v ", e)
		}
	}()
	worker.lua.Close()
}
