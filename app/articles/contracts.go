package articles

import (
	"context"

	"github.com/Faralaks/warmane-news/app/models"
	"github.com/PuerkitoBio/goquery"
)

type pageGetter interface {
	GetNews() (*goquery.Document, error)
}

type parser interface {
	ParseNews(doc *goquery.Document) ([]*models.Article, error)
}

type Storage interface {
	GetLastArticle(ctx context.Context) (*models.Article, error)
	SaveArticle(ctx context.Context, article *models.Article) error
	Exists(ctx context.Context, article *models.Article) (bool, error)
}
