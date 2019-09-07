package gif

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-resty/resty/v2"
)

const giphySearchURL = "https://api.giphy.com/v1/gifs/search"

type Rating string

const (
	RatingY    Rating = "Y"
	RatingG    Rating = "G"
	RatingPG   Rating = "PG"
	RatingPG13 Rating = "PG-13"
	RatingR    Rating = "R"
)

type Lang string

const (
	LangEN Lang = "en"
)

type GiphyClient struct {
	client   *resty.Client
	apiKey   string
	pageSize int
	lang     Lang
}

func NewGiphyClient(apiKey string) *GiphyClient {
	httpClient := resty.New()
	return &GiphyClient{
		client:   httpClient,
		apiKey:   apiKey,
		pageSize: 25,
		lang:     LangEN,
	}
}

func (c *GiphyClient) PageSize(pageSize int) *GiphyClient {
	c.pageSize = pageSize
	return c
}

func (c *GiphyClient) Lang(lang Lang) *GiphyClient {
	c.lang = lang
	return c
}

func (c *GiphyClient) search(query string, page int) (*resty.Response, error) {
	req := c.client.R().
		SetQueryParams(map[string]string{
			"q":       query,
			"api_key": c.apiKey,
			"lang":    string(c.lang),
			"limit":   strconv.Itoa(c.pageSize),
			"offset":  strconv.Itoa(page * c.pageSize),
		})
	return req.Get(giphySearchURL)
}

func (c *GiphyClient) Search(query string) *SearchHandle {
	return &SearchHandle{
		client: c,
		query:  query,
		page:   0,
	}
}

type SearchHandle struct {
	client *GiphyClient
	query  string
	page   int
}

func (h *SearchHandle) Next() (*SearchResult, error) {
	// Perform search
	resp, err := h.client.search(h.query, h.page)
	if err != nil {
		return nil, err
	}

	// Assert request succeeded
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("search failed: status=%s, body=%s", resp.Status(), resp.Body())
	}

	// Parse search results
	var searchResult SearchResult
	err = json.Unmarshal(resp.Body(), &searchResult)
	if err != nil {
		return nil, fmt.Errorf("error parsing search results: %s", resp.Body())
	}

	h.page += 1
	return &searchResult, nil
}

type SearchResult struct {
	Data       []GIFObject
	Pagination PaginationObject
}

type GIFObject struct {
	Title    string
	Url      string
	BitlyUrl string `json:"bitly_url"`
}

type PaginationObject struct {
	Offset     int
	TotalCount int `json:"total_count"`
	Count      int
}
