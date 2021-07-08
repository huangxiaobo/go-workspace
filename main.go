package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"money/core/config"
	"money/core/db"
	"money/core/fetch"
	"money/core/log"
	"money/core/task"
)

type Scheduler struct {
	tasks []*task.FetchTask
}

func (sd *Scheduler) crawler() {

	go func() {
		fetcher := &fetch.Fetcher{}
		fetcher.Start()
	}()

}

func (sd *Scheduler) start() {
	go sd.crawler()

}

func (sd *Scheduler) AddFetchTask(t *task.FetchTask) {
	sd.tasks = append(sd.tasks, t)
}

var ConfFilePath = flag.String("conf", "./config/money.yml", "config file path")

func init() {
	// Init log
	log.InitLog("./output/", "money", "utf-8")
	// Parse command line arguments
	flag.Parse()
	// Load config
	config.LoadConfig(*ConfFilePath)
	// Database SetUp
	db.SetUp(&config.Conf)
}

func main() {

	log.Info("Crawler")

	sd := Scheduler{}

	for _, item := range config.Conf.Crawler.Tasks {
		fetchTask := &task.FetchTask{
			Url: item.Url,
		}
		fetchTask.OnSuccess("#js-initialData", func(selection *goquery.Selection) {
			content := selection.Text()
			data := map[string]interface{}{}
			if err := json.Unmarshal([]byte(content), &data); err != nil {
				log.Error(err)
				return
			}

			items := strings.Split("initialState/entities/users", "/")
			m := data
			for _, item := range items {
				m = m[item].(map[string]interface{})

			}

			for userId, userData := range m {
				log.InfoWithFields(nil, log.Fields{"UserId": userId})

				userDataStr, err := json.MarshalIndent(userData, "", "    ")
				if err != nil {
					log.Error(err)
					continue
				}

				pp := map[string]interface{}{}
				json.Unmarshal(userDataStr, &pp)

				log.Info(fmt.Sprintf("%+v", pp))
			}
		})
		sd.AddFetchTask(fetchTask)

	}

	sd.start()

	for {
		time.Sleep(time.Microsecond)
	}

}
