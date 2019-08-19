package handlers_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/timucingelici/news-gateway/internal/app/handlers"
	"github.com/timucingelici/news-gateway/pkg/store/mocks"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock the data store
var storeMock = &mocks.DataStore{}

// Register the handlers
var handler = handlers.SetupRoutes(storeMock)

func TestRestHandler_GetAllNews(t *testing.T) {

	// Data from the DB
	var mockNews = map[string]string{
		"BBC.Technology.1565105946": `{"title":"8chan: Where are users going now?","link":"https://www.bbc.co.uk/news/technology-49249574","description":"The site's owner says a million people are now \"looking for a home\" after 8chan was driven offline.","datetime":"2019-08-06T16:39:06+01:00","provider":"BBC","category":"Technology","thumbnail":"https://loremflickr.com/320/240"}`,
		"BBC.UK.1566030230":         `{"title":"The Troubles: Dealing with past is priority for future","link":"https://www.bbc.co.uk/news/uk-northern-ireland-49375286","description":"As Northern Ireland looked back at 1969, it showed there is no agreed narrative on the Troubles.","datetime":"2019-08-17T09:23:50+01:00","provider":"BBC","category":"UK","thumbnail":"https://loremflickr.com/320/240"}`,
	}

	// Data expected from the endpoint
	var expected = `[{"id":"BBC.UK.1566030230","title":"The Troubles: Dealing with past is priority for future","link":"https://www.bbc.co.uk/news/uk-northern-ireland-49375286","description":"As Northern Ireland looked back at 1969, it showed there is no agreed narrative on the Troubles.","datetime":"2019-08-17T09:23:50+01:00","provider":"BBC","category":"UK","thumbnail":"https://loremflickr.com/320/240"},{"id":"BBC.Technology.1565105946","title":"8chan: Where are users going now?","link":"https://www.bbc.co.uk/news/technology-49249574","description":"The site's owner says a million people are now \"looking for a home\" after 8chan was driven offline.","datetime":"2019-08-06T16:39:06+01:00","provider":"BBC","category":"Technology","thumbnail":"https://loremflickr.com/320/240"}]`

	// Mock the data store call to get mocked news
	storeMock.On("GetAllWithValues", "*").Return(mockNews, nil)

	r, err := http.NewRequest("GET", "/news", nil)

	assert.Nil(t, err)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	resp := w.Result()

	defer func(t *testing.T) {
		err := resp.Body.Close()
		assert.Nil(t, err)
	}(t)

	body, err := ioutil.ReadAll(resp.Body)

	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, expected, string(body))
}

func TestRestHandler_GetNews(t *testing.T) {
	// Data from DB
	var mockNews = `{"title":"8chan: Where are users going now?","link":"https://www.bbc.co.uk/news/technology-49249574","description":"The site's owner says a million people are now \"looking for a home\" after 8chan was driven offline.","datetime":"2019-08-06T16:39:06+01:00","provider":"BBC","category":"Technology","thumbnail":"https://loremflickr.com/320/240"}`

	// Data expected from the endpoint
	var expected = `{"id":"BBC.Technology.1565105946","title":"8chan: Where are users going now?","link":"https://www.bbc.co.uk/news/technology-49249574","description":"The site's owner says a million people are now \"looking for a home\" after 8chan was driven offline.","datetime":"2019-08-06T16:39:06+01:00","provider":"BBC","category":"Technology","thumbnail":"https://loremflickr.com/320/240"}`

	// Mock the data store call to get mocked news
	storeMock.On("Get", "BBC.Technology.1565105946").Return(mockNews, nil)

	r, err := http.NewRequest("GET", "/news/BBC.Technology.1565105946", nil)

	assert.Nil(t, err)

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	resp := w.Result()

	defer func(t *testing.T) {
		err := resp.Body.Close()
		assert.Nil(t, err)
	}(t)

	body, err := ioutil.ReadAll(resp.Body)

	assert.Nil(t, err)
	assert.Equal(t, resp.StatusCode, http.StatusOK)
	assert.Equal(t, expected, string(body))
}
