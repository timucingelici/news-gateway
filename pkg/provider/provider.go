package provider

import "time"

var providers []Provider

type Feed struct {
	Category string `json:"category"`
	Url      string `json:"url"`
	Layout   string `json:"layout"`
}

type feedList []Feed

type Provider struct {
	Name  string   `json:"name"`
	Feeds feedList `json:"feeds"`
}

func init() {
	providers = []Provider{
		{
			"BBC",
			feedList{
				Feed{
					Category: "UK",
					Url:      "http://feeds.bbci.co.uk/news/uk/rss.xml",
					Layout:   time.RFC1123,
				},
				Feed{
					Category: "Technology",
					Url:      "http://feeds.bbci.co.uk/news/technology/rss.xml",
					Layout:   time.RFC1123,
				},
			},
		},
		{
			"Reuters",
			feedList{
				Feed{
					Category: "UK",
					Url:      "http://feeds.reuters.com/reuters/UKdomesticNews?format=xml",
					Layout:   time.RFC1123Z,
				},
				Feed{
					Category: "Technology",
					Url:      "http://feeds.reuters.com/reuters/technologyNews?format=xml",
					Layout:   time.RFC1123Z,
				},
			},
		},
	}
}

func GetAll() []Provider {
	return providers
}

func FeedCount() int {
	var result int

	for _, p := range providers {
		result += len(p.Feeds)
	}

	return result
}
