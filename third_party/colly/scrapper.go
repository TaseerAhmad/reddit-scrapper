package colly

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"log"
	"reddit-scrapper/models"
	scrapLogger "reddit-scrapper/util"
	"time"
)

var posts []models.Post //TODO TEMP

var collector *colly.Collector

func Init(domains ...string)  {
	collector = colly.NewCollector(
		colly.AllowedDomains(domains...),
		colly.AllowURLRevisit(),
		colly.Async(true))

	err := collector.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 1,
		RandomDelay: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("[FATAL] Limit Error: ", err)
		return
	}

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})
}

// Start keepAlive when true makes the crawler to fetch the data after every defined intervals in refreshRate (seconds)
func Start(keepAlive bool, refreshRate int, url ,fiName string) {
	if collector == nil {
		log.Fatal("[FATAL] Scrapper not initialized. Call to Init() required")
		return
	}

	ticker := time.NewTicker(time.Second)
	refreshInterval := time.Duration(refreshRate)

	collector.OnRequest(func(request *colly.Request) {
		//requestsMade++ //TODO Log request count
		if keepAlive {
			ticker.Reset(refreshInterval * time.Second)
			fmt.Println("[INFO] Crawler will restart in", refreshRate, "seconds")
		}
	})

	if keepAlive {
		for {
			select {
			case <-ticker.C:
				startAndLogCrawl(collector, url, fiName)
			}
		}
	} else {
		startAndLogCrawl(collector, url, fiName)
	}
}

func startAndLogCrawl(c *colly.Collector, url, fiName string) {
	startCrawl(collector, url)
	collector.Wait()
	go scrapLogger.LogToJson(posts, fiName)
}

func refineDomain(domain string) string {
	return domain[1 : len(domain)-1]
}

func startCrawl(collector *colly.Collector, url string) {
	err := collector.Visit(url) //TODO Pick the subreddits from user
	if err != nil {
		log.Fatal("[FATAL] Visit Error: ", err.Error())
		return
	}

	collector.OnHTML(".top-matter", func(e *colly.HTMLElement) {
		var post models.Post
		post.Title = e.ChildText("a[data-event-action=title]")
		post.Domain = refineDomain(e.ChildText("span[class=domain]")) //TODO Find a way to extract directly
		post.Url = e.ChildAttr("a[data-event-action=title]", "href")

		selector := e.DOM.Find("p")
		if selector.HasClass("tagline") {
			selector.Children().Each(func(i int, selection *goquery.Selection) {
				if dateTime, exists := e.DOM.Find("time").Attr("datetime"); exists {
					post.PostedOn = dateTime
				}

				if selection.HasClass("author may-blank") {
					post.Author = selection.Contents().Text()
				}
			})
		}
		posts = append(posts, post)
	})
}