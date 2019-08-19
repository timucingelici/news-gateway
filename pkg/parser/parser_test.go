package parser_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/timucingelici/news-gateway/pkg/parser"
	"testing"
	"time"
)

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

func TestParser_Parse(t *testing.T) {
	items, err := parser.New().Parse([]byte(data), time.RFC1123)

	assert.Nil(t, err)
	assert.Equal(t, 1, len(items))

	assert.Equal(t, "Andrew Harper death: Thames Valley Police officer killed in Berkshire", items[0].Title)
	assert.Equal(t, "https://www.bbc.co.uk/news/uk-england-49368649", items[0].Link)
	assert.Equal(t, "Fri, 16 Aug 2019 18:51:26 GMT", items[0].PubDate)
	assert.Equal(t, "2019-08-16 19:51:26 +0100 BST", items[0].DateTime.String())
	assert.Equal(t, "Newly-wed PC Andrew Harper died while investigating reports of a burglary, Thames Valley Police says.", items[0].Description)
}
