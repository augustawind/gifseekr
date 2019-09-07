package gifloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/suite"
)

var mockSearchResultPath = path.Join(".", "test_fixtures", "mock-search-result.json")

type GiphyTestSuite struct {
	suite.Suite
	searchURLRegexp  *regexp.Regexp
	mockSearchResult []byte
	client           *GiphyClient
}

func (s *GiphyTestSuite) SetupSuite() {
	s.searchURLRegexp = regexp.MustCompile(fmt.Sprintf("%s.*", giphySearchURL))
	var err error
	s.mockSearchResult, err = ioutil.ReadFile(mockSearchResultPath)
	if err != nil {
		panic(err)
	}
}

func (s *GiphyTestSuite) SetupTest() {
	s.client = NewGiphyClient("foo-api-key-123")
	httpmock.ActivateNonDefault(s.client.client.GetClient())
	responder := httpmock.NewStringResponder(http.StatusOK, string(s.mockSearchResult))
	httpmock.RegisterRegexpResponder("GET", s.searchURLRegexp, responder)
}

func (s *GiphyTestSuite) TearDownTest() {
	httpmock.DeactivateAndReset()
}

func (s *GiphyTestSuite) TestBasicSearch() {
	handle := s.client.Search("some query")
	result, err := handle.Next()
	s.Require().Nil(err)

	var expected *SearchResult
	err = json.Unmarshal(s.mockSearchResult, &expected)
	s.Require().Nil(err)

	s.Require().Equal(expected, result)
}

func TestGiphyTestSuite(t *testing.T) {
	suite.Run(t, new(GiphyTestSuite))
}
