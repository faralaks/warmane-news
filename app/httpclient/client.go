package httpclient

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Urls Represents the urls of the Warmane website
type Urls struct {
	News string
}

type Client struct {
	urls       *Urls
	httpClient *http.Client
}

func NewClient(urls *Urls, httpClient *http.Client) *Client {
	return &Client{urls: urls, httpClient: httpClient}
}

// GetNews Returns the news page from the Warmane website
func (client *Client) GetNews() (*goquery.Document, error) {
	response, err := client.httpClient.Get(client.urls.News)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}
