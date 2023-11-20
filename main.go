package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

// Article Represents a Warmane article
type Article struct {
	Title   string
	DateUTC time.Time
	Body    string
}

// GetISODate Returns the date in ISO format
func (article *Article) GetISODate() string {
	return article.DateUTC.Format("2006-01-02")
}

func main() {
	logConf := zap.NewDevelopmentConfig()
	logConf.DisableStacktrace = true
	logConf.Encoding = "console"

	logger, err := logConf.Build()
	if err != nil {
		panic(err)
	}

	url := "https://warmane.com"

	response, err := http.Get(url)
	if err != nil {
		logger.Error("main: http.Get warmane webpage", zap.Error(err))
		return
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		logger.Error("main: goquery.NewDocumentFromReader", zap.Error(err))
		return
	}

	articles, err := getArticles(logger, doc)
	if err != nil {
		logger.Error("main: getArticles", zap.Error(err))
		return
	}

	for _, article := range articles {
		fmt.Println(article.Title)
		fmt.Println(article.DateUTC)
		fmt.Println(article.Body)
		fmt.Println()
	}

}

func getArticles(logger *zap.Logger, doc *goquery.Document) ([]*Article, error) {
	articlesBlock := doc.Find(".page-articles-left")
	articles := make([]*Article, 0)

	titleList := make([]string, 0)
	dateList := make([]string, 0)
	bodyList := make([]string, 0)

	articlesBlock.Find(".wm-ui-article-title").Each(func(_ int, s *goquery.Selection) {
		paddings := s.Find("p")
		titleList = append(titleList, paddings.Eq(0).Text())
		dateList = append(dateList, paddings.Eq(1).Text())
	})

	articlesBlock.Find(".wm-ui-article-content").Each(func(_ int, s *goquery.Selection) {
		bodyList = append(bodyList, s.Find("p").Eq(0).Text())
	})

	if len(titleList) != len(dateList) || len(titleList) != len(bodyList) {
		logger.Error("getArticles: different length of titleList, dateList and bodyList", zap.Int("titleList", len(titleList)), zap.Int("dateList", len(dateList)), zap.Int("bodyList", len(bodyList)))
		return nil, fmt.Errorf("different length of titleList, dateList and bodyList")
	}

	for i := 0; i < len(titleList); i++ {
		title := titleList[i]
		if title == "" {
			logger.Error("getArticleElems: title is empty")
			continue
		}

		body := bodyList[i]
		if body == "" {
			logger.Error("getArticleElems: body is empty")
			continue
		}

		dateStr := dateList[i]
		date, err := time.Parse("January 2, 2006", dateStr)
		if err != nil {
			logger.Error("getArticleElems: Failed to parse date", zap.Error(err), zap.String("DateString", dateStr), zap.String("Title", title))
			continue
		}

		article := &Article{
			Title:   title,
			DateUTC: date.UTC(),
			Body:    body,
		}

		articles = append(articles, article)
	}

	return articles, nil
}
