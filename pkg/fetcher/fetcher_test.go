package fetcher_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/timucingelici/news-gateway/pkg/fetcher"
	"net/http"
	"net/http/httptest"
	"testing"
)

type xmlServer struct {
	t *testing.T
}

var data = `<?xml version="1.0" encoding="UTF-8"?>
<rss>
    <channel>
        <item>
            <title><![CDATA[Andrew Harper death: Thames Valley Police officer killed in Berkshire]]></title>
            <description><![CDATA[Newly-wed PC Andrew Harper died while investigating reports of a burglary, Thames Valley Police says.]]></description>
            <link>https://www.bbc.co.uk/news/uk-england-49368649</link>
            <guid isPermaLink="true">https://www.bbc.co.uk/news/uk-england-49368649</guid>
            <pubDate>Fri, 16 Aug 2019 18:51:26 GMT</pubDate>
        </item>
    </channel>
</rss>`

func TestFetcher_FetchShouldSuccess(t *testing.T) {
	handler := &xmlServer{t}

	server := httptest.NewServer(handler)
	defer server.Close()

	body, err := fetcher.New(server.URL + "/bbc.xml").Fetch()

	assert.Nil(t, err)
	assert.Equal(t, data, string(body))
}

func TestFetcher_FetchShouldFail(t *testing.T) {
	handler := &xmlServer{t}

	server := httptest.NewServer(handler)
	defer server.Close()

	body, err := fetcher.New(server.URL + "/404.xml").Fetch()

	assert.Error(t, err)
	assert.NotEqual(t, data, string(body))
}

func (h *xmlServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.String() == "/bbc.xml" {
		w.WriteHeader(200)
		_, err := w.Write([]byte(data))

		if err != nil {
			h.t.Error("Failed to respond on the test server! ", err)
		}
	}

	w.WriteHeader(404)
}
