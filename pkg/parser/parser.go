package parser

import (
	"encoding/xml"
	"time"
)

func New() Parser {
	return &parser{}
}

type Parser interface {
	Parse([]byte, string) ([]Item, error)
}

type parser struct{}

func (p *parser) Parse(data []byte, layout string) ([]Item, error) {
	resp := News{}
	err := xml.Unmarshal(data, &resp)

	if err != nil {
		return resp.Items, err
	}

	for i, item := range resp.Items {
		t, err := time.Parse(layout, item.PubDate)

		if err != nil {
			return resp.Items, err
		}

		resp.Items[i].DateTime = t.Local()
	}

	return resp.Items, err
}

type News struct {
	Items []Item `xml:"channel>item"`
}

type Item struct {
	Title       string    `xml:"title" json:"title"`
	Link        string    `xml:"link" json:"link"`
	Description string    `xml:"description" json:"description"`
	PubDate     string    `xml:"pubDate" json:"-"`
	DateTime    time.Time `json:"datetime"`
	Provider    string    `json:"provider"`
	Category    string    `json:"category"`
	Thumbnail   string    `json:"thumbnail"`
}
