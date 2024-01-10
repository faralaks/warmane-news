package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Faralaks/warmane-news/app/db"
	"github.com/Faralaks/warmane-news/app/models"
	_ "github.com/mattn/go-sqlite3"
)

// Storage Implements methods to work with sqlite3 storage
type Storage struct {
	db db.Ext
}

var getLastArticleQuery = `
SELECT title, date_utc, body
FROM articles
ORDER BY date_utc DESC
LIMIT 1
`

// GetLastArticle retrieves the last article from the storage
func (s *Storage) GetLastArticle(ctx context.Context) (*models.Article, error) {
	row := s.db.QueryRow(ctx, getLastArticleQuery)

	article, err := scanArticle(row)
	if err != nil {
		return nil, err
	}

	return article, nil
}

var saveArticleQuery = `
INSERT INTO articles (title, date_utc, body)
VALUES (?, ?, ?)
`

// SaveArticle Saves an article to the storage
func (s *Storage) SaveArticle(ctx context.Context, article *models.Article) error {
	strDate := article.DateUTC.Format("2006-01-02")
	_, err := s.db.Exec(ctx, saveArticleQuery,
		article.Title, strDate, article.Body,
	)
	return err
}

var getArticleQuery = `
SELECT title, date_utc, body
FROM articles
WHERE title = ? AND date_utc = ? AND body = ?
ORDER BY date_utc DESC
`

// Exists Checks if an article already exists in the storage
func (s *Storage) Exists(ctx context.Context, article *models.Article) (bool, error) {
	strDate := article.DateUTC.Format("2006-01-02")
	row := s.db.QueryRow(ctx, getArticleQuery,
		article.Title, strDate, article.Body,
	)

	_, err := scanArticle(row)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("failed to scan row: %w", err)
	}

	return true, nil
}

// scanArticle Scans an article from the database
func scanArticle(row *sql.Row) (*models.Article, error) {
	if row == nil {
		return nil, db.ErrNotFound
	}

	var stringDateUTC string
	var article models.Article

	err := row.Scan(
		&article.Title,
		&stringDateUTC,
		&article.Body)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, db.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	article.DateUTC, err = time.Parse("2006-01-02", stringDateUTC)
	if err != nil {
		return nil, fmt.Errorf("failed to parse date: %v: %w", err, db.ErrParsingFaled)
	}

	return &article, nil
}
