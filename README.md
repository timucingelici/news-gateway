# # News Gateway API
This service (or the collection of services) will fetch, parse, and store news from different RSS feeds and then expose the collected data to its consumers over a RESTful API

## Data Sources
The project uses the following 4 RSS feeds to fetch data.

* BBC
	* http://feeds.bbci.co.uk/news/uk/rss.xml (BBC News UK)
	* http://feeds.bbci.co.uk/news/technology/rss.xml (BBC News Technology)

* Reuters
	* http://feeds.reuters.com/reuters/UKdomesticNews?format=xml (Reuters UK)
	* http://feeds.reuters.com/reuters/technologyNews?format=xml (Reuters Technology)


## Architecture
The following infrastructure diagram is to show how I would design this if I had weeks and a team.

Still, I was able to up come up with a working `fetcher` service and a working `api` solution.

![](https:i.imgur.com/DZENGr7.png)


## Running the project
Easiest way to see the demo is spinning up the docker containers after cloning the repo.
* `git clone https://github.com/timucingelici/news-gateway`
* `docker-compose up` 

## API Design and examples.
The API has 3 endpoints.

* `GET /news` - will return everything from the data store (`GET /news?provider=BBC&category=UK` will return the news for the given provider and/or category)

* `GET /news/{newsID}` will return a single news from the DB.

* `POST /news/{newsID}/share` will act like as if it is sending an email to the provided address.

The `share` endpoint only expects an `email` field with a valid email address.

Following `curl` commands can be used to test the endpoints.

* To list all the news
`curl http://localhost:8090/news`

* While listing all the news, filtering by provider and category is possible;

`curl http://localhost:8090/news?provider=BBC`  

`curl http://localhost:8090/news?provider=BBC&category=Technology`  

* To get the details of a single news (you may need to change the news ID)
`curl http://localhost:8090/news/BBC.UK.1566194051` 

This endpoint will return 404 if the news doesn’t exist and 500 if it can’t fetch the news from the data store.

* To share a news with someone via email
`curl -X POST http://localhost:8090/news/BBC.UK.1566194051/share -data "email=shanon@example.com"` 

This endpoint will return 404 if the news doesn’t exist and 400 if you fail to provide the email field or a valid email.

## Notes ( more like confessions)
* My plan was creating 2 services, the API and the fetcher. I started with fetcher and picked Redis as my data store but then I realise Redis was the worst choice for this task, not only because the data transformation I would need but I also discovered that the Redis client I use was not thread safe.

So my `fetcher` fetches the URLs concurrently but writing into Redis through a blocking unbuffered channel.

It doesn’t create any problems for this PoC but it would not be ideal for a real world application.

* `provider` package was the first thing I created, only to hold data for my fetcher and move the data into a JSON file but it is still there. I simply run out of time while working on other things.

* I should have used a package for reading environment variables into my `config` solution.

* Almost all the comments and examples are missing, again, simply because spending too much time on other things and running out of time.

* I am not super happy with my http handlers. They are kind of repetitive it’s mostly because of my choice of store.

* I needed to create unique news ID’s but the feeds were not providing anything unique, so I converted the date times to unix timestamps and used them to create unique keys. It should be OK for this PoC since neither of the source has posted two news at the same second.

* There are lots of improvements I can see on my code and how I structure things but I promised to send it on Monday morning, so I am leaving it as it is.

## Possible problems
Never happened during my tests but `fetcher` can try to write to Redis if Redis takes too long to spin up. 

That would kill the `fetcher` since it uses `log.Fatal` to log that kind of problems.