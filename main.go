package main

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"go.uber.org/zap"
)

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
