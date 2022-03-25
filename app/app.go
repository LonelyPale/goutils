package app

import (
	log "github.com/sirupsen/logrus"

	"github.com/lonelypale/goutils"
)

type App struct {
	Runners []Runner
}

func NewApp(runs ...Runner) *App {
	return &App{
		Runners: append([]Runner{}, runs...),
	}
}

// Start 启动 Boot 应用
func (app *App) Start() {
	for _, run := range app.Runners {
		go run.Start()
	}
}

// Stop 停止 Boot 应用
func (app *App) Stop() {
	// OnStopApplication 是否需要有 Timeout 的 Context？
	// 仔细想想没有必要，程序想要优雅退出就得一直等，等到所有工作
	// 做完，用户如果等不急了可以使用 kill -9 进行硬杀，也就是
	// 是否优雅退出取决于用户。这样的话，OnStopApplication 不
	// 依赖 appCtx 的 Context，就只需要考虑 SafeGoroutine
	// 的退出了，而这只需要 Context 一 cancel 也就完事了。

	log.Info("app boot exiting")
	var wg goutils.WaitGroup
	for _, run := range app.Runners {
		r := run                    // 避免延迟绑定
		wg.Add(func() { r.Stop() }) // 异步执行
	}
	wg.Wait()
	log.Info("app boot exited")
}
