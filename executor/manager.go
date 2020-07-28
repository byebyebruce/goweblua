package executor

import (
	"context"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	l4g "github.com/alecthomas/log4go"
)

const (
	k_MAX_MSG                   = 256                    //最大消息数
	k_PROCESS_TIMEOUT           = 200 * time.Millisecond //超时时间
	k_AfterTimesCreateNewWorker = 512                    //计算这么多次之后强制创建新的worker
)

//param 战斗参数
type param struct {
	Data    []byte
	RetChan chan<- []byte
}

var (
	msgChan          chan *param           //待计算队列
	closeChan        chan struct{}         //关闭
	wg               sync.WaitGroup        //wg
	recreateChanMap  map[int]chan struct{} //重建worker通道map
	recreateChanLock sync.Mutex            //重建锁
)

func init() {
	msgChan = make(chan *param, k_MAX_MSG)
	closeChan = make(chan struct{})
	wg = sync.WaitGroup{}
	recreateChanMap = make(map[int]chan struct{})
}

//Run 主逻辑
func Run(workerNum int, file string, funcName string, searchPath ...string) {
	recreateChanLock.Lock()
	defer recreateChanLock.Unlock()

	if workerNum <= 0 {
		//如果没有指定个数就用CPU核心数
		workerNum = runtime.NumCPU()
	}

	l4g.Info("[executor] worker num=[%d]", workerNum)

	// 加载锁
	mtxNew := sync.Mutex{}

	//创建worker
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		recreateChanMap[i] = make(chan struct{}, 1)

		//创建worker
		worker, err := NewLuaWorker(file, funcName, searchPath...)
		if nil != err {
			panic(err)
		}
		if err := worker.Init(); nil != err {
			panic(err)
		}

		l4g.Info("[executor] create worker[%d]", i)

		go func(idx int, cmdChan chan struct{}) {
			defer wg.Done()

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
					newWorker, err := NewLuaWorker(file, funcName, searchPath...)
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

				case p := <-msgChan: // 消息队列
					start := time.Now()
					l4g.Info("[executor] worker[%d] begin", idx)

					//计算
					ret, err := worker.Execute(p.Data)
					if nil != err {
						l4g.Error("[executor] worker[%d].Execute error=%s", idx, err.Error())
					}

					elapsed := time.Since(start)

					p.RetChan <- ret

					if elapsed > k_PROCESS_TIMEOUT {
						l4g.Warn("[executor] worker[%d] end !!too slow!! time=[%f]s, ", idx, elapsed.Seconds())
					} else {
						l4g.Info("[executor] worker[%d] end time=[%f]s, ", idx, elapsed.Seconds())
					}

					// 累计次数加1
					execTimes++
					// 需要new新的worker
					if 0 == execTimes%k_AfterTimesCreateNewWorker {
						select {
						case cmdChan <- struct{}{}:
						default:
						}
					}
				case <-closeChan: // 外面关闭了
					l4g.Info("[executor] worker[%d] is exiting...", idx)
					break LOOP

				}

			}

			// 外面关闭了，但是还是要把msgChan里的任务执行完毕
			for hasTask := true; hasTask; {
				select {
				case p := <-msgChan:
					start := time.Now()
					l4g.Info("[executor] worker[%d] begin", idx)

					ret, err := worker.Execute(p.Data)
					if nil != err {
						l4g.Error("[executor] worker.Execute error=%s", err.Error())
					}

					elapsed := time.Since(start)

					p.RetChan <- ret

					if elapsed > k_PROCESS_TIMEOUT {
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

		}(i, recreateChanMap[i])

	}

}

//Process 执行参数
func Process(ctx context.Context, data []byte) <-chan []byte {

	retChan := make(chan []byte, 1)

	p := &param{
		Data:    data,
		RetChan: retChan,
	}

	select {
	case <-ctx.Done():
		l4g.Error("[executor] Process msgChan is full error:[%s]", ctx.Err().Error())
	case msgChan <- p:
	}

	return retChan

}

// RecreateNewWorker 创建新的worker
func RecreateNewWorker() {
	recreateChanLock.Lock()
	defer recreateChanLock.Unlock()
	for _, v := range recreateChanMap {
		select {
		case v <- struct{}{}:
		default:
		}
	}
	l4g.Error("[executor] CreateNewWorker")
}

// Close 关闭
func Close() {
	close(closeChan)
	wg.Wait()
	remain := len(msgChan)
	if remain > 0 {
		l4g.Critical("[executor] remain = [%d]", remain)
	}
	l4g.Warn("[executor] Close")
}
