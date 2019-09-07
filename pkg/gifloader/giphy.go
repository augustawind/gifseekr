package gifloader

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
		return nil, fmt.Errorf("error parsing search results: %s", err)
	}

	h.page += 1
	return &searchResult, nil
}

type SearchResult struct {
	Pagination PaginationObject `json:"pagination"`
	Data       []GIFObject      `json:"data"`
}

type PaginationObject struct {
	Offset     int `json:"offset"`
	TotalCount int `json:"total_count"`
	Count      int `json:"count"`
}

type GIFObject struct {
	ID       string       `json:"id"`
	Title    string       `json:"title"`
	URL      string       `json:"url"`
	BitlyURL string       `json:"bitly_url"`
	Images   ImagesObject `json:"images"`
}

type ImagesObject struct {
	Original      ImageVariant `json:"original"`
	FixedWidth200 ImageVariant `json:"fixed_width"`
	Downsized2MB  ImageVariant `json:"downsized"`
	Downsized5MB  ImageVariant `json:"downsized_medium"`
	Downsized8MB  ImageVariant `json:"downsized_large"`
}

type ImageVariant struct {
	URL    string `json:"url"`
	Width  int64  `json:"width"`
	Height int64  `json:"height"`
	Size   int64  `json:"size"`
}

func (img *ImageVariant) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	aux := &struct {
		URL    string `json:"url"`
		Width  string `json:"width"`
		Height string `json:"height"`
		Size   string `json:"size"`
	}{}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	img.URL = aux.URL
	if img.Width, err = strconv.ParseInt(aux.Width, 10, 64); err != nil {
		return err
	}
	if img.Height, err = strconv.ParseInt(aux.Height, 10, 64); err != nil {
		return err
	}
	if img.Size, err = strconv.ParseInt(aux.Size, 10, 64); err != nil {
		return err
	}
	return nil
}
