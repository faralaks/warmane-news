package models

import "time"

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
