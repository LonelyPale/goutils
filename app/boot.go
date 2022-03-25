package app

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// Runner 应用执行器
type Runner interface {
	Start() // 启动执行器
	Stop()  // 关闭执行器
}

var exitChan chan struct{}

func Boot(runs ...Runner) {
	Run(NewApp(runs...))
}

// Run 启动执行器
func Run(runner Runner) {
	exitChan = make(chan struct{})

	// 响应控制台的 Ctrl+C 及 kill 命令。
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		log.Info("got signal, program will exit")
		Exit()
	}()

	runner.Start() // 注意: runner.Start() 一定不能阻塞主线程
	<-exitChan
	runner.Stop()
}

// Exit 关闭执行器
func Exit() {
	closeChan(exitChan)
}

func closeChan(ch chan struct{}) {
	select {
	case <-ch:
		// chan 已关闭，无需再次关闭。
	default:
		close(ch)
	}
}
