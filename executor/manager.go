package executor

import (
	"context"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	l4g "github.com/alecthomas/log4go"
)

// Config 配置
type Config struct {
	WorkerNum     int
	MaxTask       int    // 最大消息数
	TimeoutMS     int64  // 超时时间
	RecycleTimes  int    // 计算这么多次之后强制创建新的worker
	LuaEntryFile  string // lua入口脚本
	LuaEntryFunc  string // lua入口函数
	LuaSearchPath string // lua搜索目录
}

//param 战斗参数
type param struct {
	Data    []byte
	RetChan chan<- []byte
}

// Manager 管理器
type Manager struct {
	cfg              Config        // 配置
	msgChan          chan *param   //计算队列
	closeChan        chan struct{} //关闭
	wg               sync.WaitGroup
	recreateChanMap  map[int]chan struct{} //重建worker通道map
	recreateChanLock sync.Mutex            //重建锁
}

// NewManager 构造
func NewManager(cfg Config) *Manager {
	m := &Manager{
		cfg:             cfg,
		msgChan:         make(chan *param, cfg.MaxTask),
		closeChan:       make(chan struct{}),
		recreateChanMap: make(map[int]chan struct{}),
	}
	return m
}

// Run run
func (m *Manager) Run() {
	m.recreateChanLock.Lock()
	defer m.recreateChanLock.Unlock()

	if m.cfg.WorkerNum <= 0 {
		//如果没有指定个数就用CPU核心数
		m.cfg.WorkerNum = runtime.NumCPU()
	}

	l4g.Info("[executor] worker num=[%d]", m.cfg.WorkerNum)

	// 加载锁
	mtxNew := sync.Mutex{}

	//创建worker
	for i := 0; i < m.cfg.WorkerNum; i++ {
		m.wg.Add(1)
		m.recreateChanMap[i] = make(chan struct{}, 1)

		//创建worker
		worker, err := NewLuaWorker(m.cfg.LuaEntryFile, m.cfg.LuaEntryFunc, m.cfg.LuaSearchPath)
		if nil != err {
			panic(err)
		}
		if err := worker.Init(); nil != err {
			panic(err)
		}

		l4g.Info("[executor] create worker[%d]", i)

		go func(idx int, cmdChan chan struct{}) {
			defer m.wg.Done()

			defer func() {
				e := recover()
				if nil != e {
					backtrace := debug.Stack()
					l4g.Critical("[executor] worker[%d] panic e:[%v]\nstack[%s]", idx, e, string(backtrace))
				}
			}()

			// 计算次数
			execTimes := 0
		LOOP:
			for {
				select {
				case <-cmdChan: // 创建一个新的worker
					mtxNew.Lock() // 加载脚本文件需要加锁，不然会多次打开文件报错
					newWorker, err := NewLuaWorker(m.cfg.LuaEntryFile, m.cfg.LuaEntryFunc, m.cfg.LuaSearchPath)
					if nil != err {
						l4g.Error("[executor] create new worker[%d] NewLuaWorker error=[%s]", idx, err.Error())
						mtxNew.Unlock()
						break
					}

					if err := newWorker.Init(); nil != err {
						l4g.Error("[executor] create new worker[%d] worker.Init() error=[%s]", idx, err.Error())
						mtxNew.Unlock()
						break
					}
					mtxNew.Unlock()

					// 旧的worker销毁
					worker.Destroy()
					// 用新的worker替换
					worker = newWorker
					// 重置计算次数
					execTimes = 0

					runtime.GC()
					l4g.Warn("[executor] create new worker[%d] ok", idx)

				case p := <-m.msgChan: // 消息队列
					start := time.Now()
					l4g.Info("[executor] worker[%d] begin", idx)

					//计算
					ret, err := worker.Execute(p.Data)
					if nil != err {
						l4g.Error("[executor] worker[%d].Execute error=%s", idx, err.Error())
					}

					elapsed := time.Since(start)

					p.RetChan <- ret

					if elapsed.Milliseconds() > m.cfg.TimeoutMS {
						l4g.Warn("[executor] worker[%d] end !!too slow!! time=[%f]s, ", idx, elapsed.Seconds())
					} else {
						l4g.Info("[executor] worker[%d] end time=[%f]s, ", idx, elapsed.Seconds())
					}

					// 累计次数加1
					execTimes++
					// 需要new新的worker
					if 0 == execTimes%m.cfg.RecycleTimes {
						select {
						case cmdChan <- struct{}{}:
						default:
						}
					}
				case <-m.closeChan: // 外面关闭了
					l4g.Info("[executor] worker[%d] is exiting...", idx)
					break LOOP

				}

			}

			// 外面关闭了，但是还是要把msgChan里的任务执行完毕
			for hasTask := true; hasTask; {
				select {
				case p := <-m.msgChan:
					start := time.Now()
					l4g.Info("[executor] worker[%d] begin", idx)

					ret, err := worker.Execute(p.Data)
					if nil != err {
						l4g.Error("[executor] worker.Execute error=%s", err.Error())
					}

					elapsed := time.Since(start)

					p.RetChan <- ret

					if elapsed.Milliseconds() > m.cfg.TimeoutMS {
						l4g.Warn("[executor] worker[%d] end !!too slow!! time=[%f]s, ", idx, elapsed.Seconds())
					} else {
						l4g.Info("[executor] worker[%d] end time=[%f]s, ", idx, elapsed.Seconds())
					}
				default:
					hasTask = false
				}
			}

			worker.Destroy()
			l4g.Info("[executor] worker[%d] destroy", idx)

		}(i, m.recreateChanMap[i])

	}

}

//Process 执行参数
func (m *Manager) Process(ctx context.Context, data []byte) <-chan []byte {

	retChan := make(chan []byte, 1)

	p := &param{
		Data:    data,
		RetChan: retChan,
	}

	select {
	case <-ctx.Done():
		l4g.Error("[executor] Process msgChan is full error:[%s]", ctx.Err().Error())
	case m.msgChan <- p:
	}

	return retChan

}

// RecreateNewWorker 创建新的worker
func (m *Manager) RecreateNewWorker() {
	m.recreateChanLock.Lock()
	defer m.recreateChanLock.Unlock()
	for _, v := range m.recreateChanMap {
		select {
		case v <- struct{}{}:
		default:
		}
	}
	l4g.Error("[executor] CreateNewWorker")
}

// Close 关闭
func (m *Manager) Close() {
	close(m.closeChan)
	m.wg.Wait()
	remain := len(m.msgChan)
	if remain > 0 {
		l4g.Critical("[executor] remain = [%d]", remain)
	}
	l4g.Warn("[executor] Close")
}
