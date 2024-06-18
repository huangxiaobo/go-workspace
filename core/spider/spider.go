package spider

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Spider struct {
	tasksCh chan *FetchTask
	doneCh  chan struct{}

	ctx context.Context
}

func NewSpider() *Spider {
	s := &Spider{
		tasksCh: make(chan *FetchTask, 100),
		doneCh:  make(chan struct{}),
		ctx:     context.Background(),
	}
	return s
}

func (sd *Spider) Start() {
	go func() {
		fetcher := &Fetcher{}
		fetcher.Start(sd.tasksCh)
	}()

	// 监听退出信号
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func(doneCh chan struct{}) {
		sig := <-sigs
		fmt.Println(sig)
		doneCh <- struct{}{}
	}(sd.doneCh)

	// 阻塞，直到进程结束
	<-sd.doneCh
}

func (sd *Spider) AddFetchTask(t *FetchTask) {
	sd.tasksCh <- t
}
