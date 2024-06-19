package spider

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/huangxiaobo/gospider/core/log"
)

var agents = []string{
	"Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:50.0) Gecko/20100101 Firefox/50.0",
	"Opera/9.80 (Macintosh; Intel Mac OS X 10.6.8; U; en) Presto/2.8.131 Version/11.11",
	"Opera/9.80 (Windows NT 6.1; U; en) Presto/2.8.131 Version/11.11",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; 360SE)",
	"Mozilla/5.0 (Windows NT 6.1; rv:2.0.1) Gecko/20100101 Firefox/4.0.1",
	"Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; The World)",
	"User-Agent,Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_6_8; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
	"User-Agent, Mozilla/4.0 (compatible; MSIE 7.0; Windows NT 5.1; Maxthon 2.0)",
	"User-Agent,Mozilla/5.0 (Windows; U; Windows NT 6.1; en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
}

/**
* 随机返回一个User-Agent
 */
func getAgent() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	size := len(agents)
	return agents[r.Intn(size)]
}

func Fetch(urlString string) (bool, string) {
	log.Info(fmt.Sprintf("download >>> url: %s", urlString))

	client := &http.Client{
		// Transport: transport,
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", urlString, nil)
	if err != nil {
		log.Fatal("new request failed,", err.Error())
		return false, ""
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("User-Agent", getAgent())
	// req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)

	if err != nil || resp == nil {
		log.InfoWithFields("do request fail>>>: ", log.Fields{"err": err})
		return false, ""
	}

	log.InfoWithFields("", log.Fields{"url": urlString, "status_code": resp.StatusCode})
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read body failed:", err.Error())
	}

	return true, string(body)
}

// FetcherManager 管理多个fetcher
type FetcherManager struct {
	// 待执行的fetch task管道
	fetchTasksCh chan *FetchTask

	// fetch worker 的数量
	fetchWorkerNum    int
	fetchWorkerWg     *sync.WaitGroup
	fetchWorkerCancel context.CancelFunc

	ctx context.Context
}

func NewFetchManager(ctx context.Context) *FetcherManager {
	f := &FetcherManager{
		fetchTasksCh:   make(chan *FetchTask, 100),
		fetchWorkerNum: 10,
		fetchWorkerWg:  nil,
		ctx:            ctx,
	}

	return f
}

func (m *FetcherManager) AddTask(task *FetchTask) error {
	m.fetchTasksCh <- task
	return nil
}

func (m *FetcherManager) Start() {
	m.fetchWorkerWg = &sync.WaitGroup{}

	workerCtx, workerCancel := context.WithCancel(context.Background())

	m.fetchWorkerCancel = workerCancel

	for i := 0; i < m.fetchWorkerNum; i++ {
		m.fetchWorkerWg.Add(1)
		fw := &FetchWorker{
			Id: fmt.Sprintf("fetch-worker-%d", i+1),
		}
		go func() {
			defer m.fetchWorkerWg.Done()
			fw.Run(workerCtx, m.fetchTasksCh)
		}()
	}

}

func (m *FetcherManager) Stop(ctx context.Context) {
	log.Info("fetch manager stop begin")
	doneCh := make(chan struct{})
	go func() {
		// 关闭任务管道
		close(m.fetchTasksCh)
		// 取消context
		m.fetchWorkerCancel()
		// 等待所有worker协程结束
		m.fetchWorkerWg.Wait()
		// 关闭完成管道
		close(doneCh)
	}()

	select {
	case <-doneCh:
		log.Info("fetch manager stop success")
	case <-ctx.Done():
		log.Info("fetch manager stop timeout")
	}

}

// Fetcher worker

type FetchWorker struct {
	Id string
}

// 通过关闭taskCh来结束worker
func (fw *FetchWorker) Run(ctx context.Context, taskCh <-chan *FetchTask) {
	for {
		select {
		case t, ok := <-taskCh:
			if !ok {
				log.Info("fetch task channel is closed: ", fw)
				return
			}
			reqUrl := t.Url
			log.Info("download url: ", reqUrl)

			ok, html := Fetch(t.Url)
			log.Info(fmt.Sprintf("download %s, status: %t ", reqUrl, ok))

			t.OnSuccess(html)
		case <-ctx.Done():
			log.Info("fetch task is closed: ", fw)
			return
		}
	}
}

func (fw FetchWorker) String() string {
	return fmt.Sprintf("FetchWorker{Id=%s}", fw.Id)
}
