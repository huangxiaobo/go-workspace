package parser

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"money/core/model"
	"strings"

	"money/core/log"
)

type ZhihuPeople struct {
}

func (p *ZhihuPeople) Parse(html string) error {

	// Load the HTML document
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Error(err)
		return err
	}

	content := dom.Find("#js-initialData").First().Text()

	data := map[string]interface{}{}
	err = json.Unmarshal([]byte(content), &data)
	if err != nil {
		log.Error(err)
		return nil
	}

	for userId, userData := range getUserData(data, "initialState/entities/users") {
		log.InfoWithFields(nil, log.Fields{"UserId": userId})

		userDataStr, err := json.MarshalIndent(userData, "", "    ")
		if err != nil {
			log.Error(err)
			continue
		}

		pp := model.People{}
		json.Unmarshal(userDataStr, &pp)

		log.Info(fmt.Sprintf("%+v", pp))
	}

	return nil
}

func getUserData(data map[string]interface{}, path string) map[string]interface{} {
	items := strings.Split(path, "/")
	m := data
	for _, item := range items {
		m = m[item].(map[string]interface{})

	}
	return m
}
