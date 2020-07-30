package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	l4g "github.com/alecthomas/log4go"
	"github.com/byebyebruce/goweblua/executor"
	"github.com/byebyebruce/goweblua/web"
)

var (
	gWeb       = flag.String("web", "10001", "web listen address")
	searchPath = flag.String("search", "assets/lua", "lua search path")
	funcName   = flag.String("func", "MyFunc", "entry function name")
	file       = flag.String("file", "assets/lua/main.lua", "entry function name")
)

func main() {

	rand.Seed(time.Now().Unix())
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	defer func() {
		l4g.Info("[main] quit...")
		l4g.Global.Close()
	}()

	cfg := executor.Config{
		WorkerNum:     8,
		MaxTask:       32,
		TimeoutMS:     1000,
		RecycleTimes:  512,
		LuaEntryFile:  *file,
		LuaEntryFunc:  *funcName,
		LuaSearchPath: *searchPath,
	}
	mgr := executor.NewManager(cfg)
	mgr.Run()

	web.Mgr = mgr

	http.HandleFunc("/", web.HTTPHandleFunc)

	go func() {
		e := http.ListenAndServe(fmt.Sprintf(":%s", *gWeb), nil)
		if nil != e {
			panic(e)
		}
	}()
	l4g.Info("[main] http.ListenAndServe port=[%s]", *gWeb)

	l4g.Info("[main] executor start...")

	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// console input
	/*
		go func() {

			for {
				input := bufio.NewScanner(os.Stdin)
				input.Scan()
				t := input.Text()
				l4g.Warn("[service] is quitting by press key[%s]", t)
				if "q" == t {
					l4g.Warn("[service] is quitting by press key[q]")
					sigs <- syscall.SIGINT
					return
				}
			}

		}()
	*/
	l4g.Info("[main] run...")
QUIT:
	for {
		select {
		case sig := <-sigs:
			l4g.Info("[main] Signal: %s=[%d]", sig.String(), sig)
			break QUIT
		case <-ticker.C:
			l4g.Info("[main] i am running...")
		}
	}

	l4g.Warn("[main] is quiting...")

	mgr.Close()
	l4g.Warn("[main] executor is stopped")
}
