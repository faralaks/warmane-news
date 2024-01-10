package parser

import (
	"fmt"
	"time"

	"github.com/Faralaks/warmane-news/app/models"
	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

// Parser Is the parser for the Warmane website
type Parser struct {
	logger *zap.Logger
}

func NewParser(logger *zap.Logger) *Parser {
	return &Parser{logger: logger}
}

// ParseNews Returns the news from the Warmane news web page
func (parser *Parser) ParseNews(doc *goquery.Document) ([]*models.Article, error) {
	logger := parser.logger.Named("ParseNews")

	articlesBlock := doc.Find(".page-articles-left")
	articles := make([]*models.Article, 0)

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
		logger.Error("ParseNews: different length of titleList, dateList and bodyList", zap.Int("titleList", len(titleList)), zap.Int("dateList", len(dateList)), zap.Int("bodyList", len(bodyList)))
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

		article := &models.Article{
			Title:   title,
			DateUTC: date.UTC(),
			Body:    body,
		}

		articles = append(articles, article)
	}

	return articles, nil
}
