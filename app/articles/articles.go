package articles

import (
	"context"
	"fmt"

	"go.uber.org/zap"
)

type Articles struct {
	logger     *zap.Logger
	parser     parser
	pageGetter pageGetter
	storage    Storage
}

func NewArticles(logger *zap.Logger, parser parser, pageGetter pageGetter, storage Storage) *Articles {
	return &Articles{logger: logger, parser: parser, pageGetter: pageGetter, storage: storage}
}

// CheckForUpdates Checks for updates on the Warmane website
func (articles *Articles) CheckForUpdates(ctx context.Context) error {
	logger := articles.logger.Named("CheckForUpdates")

	doc, err := articles.pageGetter.GetNews()
	if err != nil {
		logger.Error("CheckForUpdates: pageGetter.GetNews", zap.Error(err))
		return err
	}

	newArticles, err := articles.parser.ParseNews(doc)
	if err != nil {
		logger.Error("CheckForUpdates: parser.ParseNews", zap.Error(err))
		return err
	}

	if len(newArticles) == 0 {
		logger.Error("CheckForUpdates: No News found on web page")
		return fmt.Errorf("no News found on web page")
	}

	lastSavedArticle, err := articles.storage.GetLastArticle(ctx)
	if err != nil {
		logger.Error("CheckForUpdates: storage.GetLastArticle", zap.Error(err))
		return err
	}

	lastArticleOnPage := newArticles[0]
	if lastArticleOnPage.DateUTC == lastSavedArticle.DateUTC && lastArticleOnPage.Title == lastSavedArticle.Title && lastArticleOnPage.Body == lastSavedArticle.Body {
		logger.Info("CheckForUpdates: No updates found", zap.String("LastTitle", lastSavedArticle.Title), zap.String("LastPublishedAt", lastSavedArticle.DateUTC.String()))
		return nil
	}

	for _, article := range newArticles {
		alreadyExists, err := articles.storage.Exists(ctx, article)
		if err != nil {
			logger.Error("CheckForUpdates: Failed to check if article is already exists", zap.Error(err))
			return err
		}

		if alreadyExists {
			continue
		}

		err = articles.storage.SaveArticle(ctx, article)
		if err != nil {
			logger.Error("CheckForUpdates: Failed to save article", zap.Error(err), zap.String("Title", article.Title), zap.String("PublishedAt", article.DateUTC.String()))
			return err
		}

		logger.Info("CheckForUpdates: New article saved", zap.String("Title", article.Title), zap.String("PublishedAt", article.DateUTC.String()))
	}

	return nil
}
