package main

import (
	"encoding/json"
	"fmt"
	"github.com/timucingelici/news-gateway/pkg/config"
	"github.com/timucingelici/news-gateway/pkg/fetcher"
	"github.com/timucingelici/news-gateway/pkg/parser"
	"github.com/timucingelici/news-gateway/pkg/provider"
	"github.com/timucingelici/news-gateway/pkg/store"
	"log"
	"sync/atomic"
	"time"
)

// struct to carry parsed data over channels
type data struct {
	key   string
	value []byte
}

// Data needs to be stored will be sent to this channel
var collector = make(chan data)

// A signal will be sent to this channel when everything is done.
// The size of this buffered channel is 1 because I only need to get that done signal once
var done = make(chan bool, 1)

// This will the track the number of go routines spawned
var counter int32

func main() {

	// Watch for an interrupt signal so we can shutdown gracefully
	//c := make(chan os.Signal, 1)
	//signal.Notify(c, os.Interrupt)

	// Get the config
	conf, err := config.New()

	if err != nil {
		log.Fatalln("Failed to parse required env vars. Err : ", err)
	}

	log.Println("Current config is: ", conf)

	// Setup the data store connection
	s, err := store.New(conf.RedisProtocol, conf.RedisAddr, conf.RedisReadTimeout, conf.RedisWriteTimeout)

	if err != nil {
		log.Fatalf("Failed to connect to the data store: %s\n", err)
	}

	providers := provider.GetAll()

	// Fetcher is expected to run indefinitely and refresh the news cache in an interval define in the config
	// So this indefinite loop is to make that happen.
	//exit:
	for {

		// Number of go routines to be spawned
		counter = int32(provider.FeedCount())

		// Loop all the providers and feeds
		log.Println("Fetching news...")
		for _, p := range providers {
			for _, f := range p.Feeds {

				// Fetch the feed, parse the data and send it to the channel to get stored.
				go fetchAndParse(p, f, collector)
			}
		}

	loop:
		for {
			select {
			case d, ok := <-collector:

				// I'll break out of the loop if the channel is closed for some reason
				if !ok {
					break loop
				}

				// I can ignore the response here since it doesn't tell much.
				_, err := s.Set(d.key, d.value)

				if err != nil {
					log.Fatalf("Failed to write to the store. Err: %s", err)
				}
			//case <-c:
			//	break exit
			case <-done:
				log.Println("Fetching news done...")
				break loop
			}
		}

		time.Sleep(conf.FetcherInterval)
	}

	//log.Println("Shutting down fetcher...")
}

func fetchAndParse(p provider.Provider, f provider.Feed, c chan<- data) {

	// When finish, decrease the counter and see if this was the last one
	// If so, send a signal to the done channel so reading from the channel
	// can stop blocking
	defer func() {
		atomic.AddInt32(&counter, -1)

		// I am assigning this to a variable because if atomic.LoadInt32(&counter) == 0 returns false positive
		cc := atomic.LoadInt32(&counter)

		if cc == 0 {
			done <- true
		}
	}()

	// fetch the data
	source, err := fetcher.New(f.Url).Fetch()

	if err != nil {
		log.Printf("Failed to fetch the source for : %s. Err: %s\n", f.Url, err)
	}

	// parse the data
	items, err := parser.New().Parse(source, f.Layout)

	if err != nil {
		log.Printf("Failed to parse the source for : %s\n", f.Url)
	}

	// send the parsed data to the relevant channel to get stored

	log.Printf("Sending %d items to the store. Provider: %s, Feed: %s", len(items), p.Name, f.Category)

	for _, item := range items {

		// Create the key to hold the data
		k := fmt.Sprintf("%s.%s.%d", p.Name, f.Category, item.DateTime.Unix())

		item.Provider = p.Name
		item.Category = f.Category
		item.Thumbnail = "https://loremflickr.com/320/240"

		// Marshall it to JSON
		v, err := json.Marshal(item)

		if err != nil {
			log.Printf("Failed to encode an item from the source of : %s\nErr: %s\nItem: %v\n", f.Url, err, item)
			continue // Don't send it to the store if it's a bad egg
		}

		c <- data{k, v}
	}
}
