package zhihu

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"money/core/log"
)

type People struct {
	AccountStatus     []interface{} `json:"accountStatus"`
	AllowMessage      bool          `json:"allowMessage"`
	AnswerCount       int           `json:"answerCount"`
	ArticlesCount     int           `json:"articlesCount"`
	AvatarHue         string        `json:"avatarHue"`
	AvatarUrl         string        `json:"avatarUrl"`
	AvatarUrlTemplate string        `json:"avatarUrlTemplate"`
	Badge             []interface{} `json:"badge"`

	Business struct {
		AvatarUrl string `json:"avatarUrl"`
		Id        string `json:"id"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		Url       string `json:"url"`
	} `json:"business"`
	ColumnsCount            int    `json:"columnsCount"`
	CommercialQuestionCount int    `json:"commercialQuestionCount"`
	CoverUrl                string `json:"coverUrl"`
	Description             string `json:"description"`
	Educations              []struct {
		Diploma int `json:"diploma"`
		School  struct {
			AvatarUrl string `json:"avatarUrl"`
			Id        string `json:"id"`
			Name      string `json:"name"`
			Type      string `json:"type"`
			Url       string `json:"url"`
		} `json:"school"`
	} `json:"educations"`
	Employments            []interface{} `json:"employments"`
	FavoriteCount          int           `json:"favoriteCount"`
	FavoritedCount         int           `json:"favoritedCount"`
	FollowerCount          int           `json:"followerCount"`
	FollowingColumnsCount  int           `json:"followingColumnsCount"`
	FollowingCount         int           `json:"followingCount"`
	FollowingFavlistsCount int           `json:"followingFavlistsCount"`
	FollowingQuestionCount int           `json:"followingQuestionCount"`
	FollowingTopicCount    int           `json:"followingTopicCount"`
	Gender                 int           `json:"gender"`
	HasApplyingColumn      bool          `json:"hasApplyingColumn"`
	Headline               string        `json:"headline"`
	HostedLiveCount        int           `json:"hostedLiveCount"`
	Id                     string        `json:"id"`
	IncludedAnswersCount   int           `json:"includedAnswersCount"`
	IncludedArticlesCount  int           `json:"includedArticlesCount"`
	IncludedText           string        `json:"includedText"`
	IsActive               int           `json:"isActive"`
	IsAdvertiser           bool          `json:"isAdvertiser"`
	IsBlocked              bool          `json:"isBlocked"`
	IsBlocking             bool          `json:"isBlocking"`
	IsFollowed             bool          `json:"isFollowed"`
	IsFollowing            bool          `json:"isFollowing"`
	IsForceRenamed         bool          `json:"isForceRenamed"`
	IsOrg                  bool          `json:"isOrg"`
	IsPrivacyProtected     bool          `json:"isPrivacyProtected"`
	IsRealname             bool          `json:"isRealname"`

	Locations []struct {
		AvatarUrl string `json:"avatarUrl"`
		Id        string `json:"id"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		Url       string `json:"url"`
	} `json:"locations"`
	LogsCount             int    `json:"logsCount"`
	MessageThreadToken    string `json:"messageThreadToken"`
	MutualFolloweesCount  int    `json:"mutualFolloweesCount"`
	Name                  string `json:"name"`
	ParticipatedLiveCount int    `json:"participatedLiveCount"`
	PinsCount             int    `json:"pinsCount"`
	QuestionCount         int    `json:"questionCount"`
	RecognizedCount       int    `json:"recognizedCount"`
	ThankFromCount        int    `json:"thankFromCount"`
	ThankToCount          int    `json:"thankToCount"`
	ThankedCount          int    `json:"thankedCount"`
	Type                  string `json:"type"`
	Url                   string `json:"url"`
	UrlToken              string `json:"urlToken"`
	UseDefaultAvatar      bool   `json:"useDefaultAvatar"`
	UserType              string `json:"userType"`

	VoteFromCount int `json:"voteFromCount"`
	VoteToCount   int `json:"voteToCount"`
	VoteupCount   int `json:"voteupCount"`
	ZvideoCount   int `json:"zvideoCount"`
}

type ParserZhihuPeople struct {
}

func (p *ParserZhihuPeople) Parse(html string) error {

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

		pp := People{}
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
