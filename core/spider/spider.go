package spider

import (
	"context"
	"fmt"
	"time"

	"github.com/huangxiaobo/gospider/core/log"
)

type Spider struct {
	tasksCh       chan *FetchTask
	taskReponseCh chan *TaskResponse

	fetchManger *FetcherManager

	// ctx, 传递结束信号
	ctx context.Context
}

func NewSpider(ctx context.Context) *Spider {

	s := &Spider{
		tasksCh:       make(chan *FetchTask, 100),
		taskReponseCh: make(chan *TaskResponse, 100),
		ctx:           ctx,
	}

	s.fetchManger = NewFetchManager(ctx)
	s.fetchManger.Start()

	return s
}

func (s *Spider) Run() {

	for {
		select {
		case t, ok := <-s.tasksCh:
			if !ok {
				log.Info("spider task channel is closed")
			}
			s.fetchManger.AddTask(t)
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Spider) AddUrl(url string, parser Parser) {
	s.tasksCh <- &FetchTask{Url: url, Parser: parser}
}

func (s *Spider) Shutdown(ctx context.Context) error {
	log.Info("spider shutdown start")
	close(s.tasksCh)
	close(s.taskReponseCh)
	s.fetchManger.Stop(ctx)

	log.Info("spider shutdown success")

	return nil
}

func (s *Spider) GracefullyShutdown() error {
	log.Info("gracefully shutdown")

	errCh := make(chan error, 1)
	defer close(errCh)

	// 执行spider退出，设置超时时间5s
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	go func() {
		if err := s.Shutdown(timeoutCtx); err != nil {
			log.Fatal("shutdown failed: ", err)
			errCh <- err
		} else {
			errCh <- nil
		}
	}()

	select {
	case <-timeoutCtx.Done():
		if timeoutCtx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("shutdown timeout")
		}
	case err := <-errCh:
		return err
	}

	return nil
}
