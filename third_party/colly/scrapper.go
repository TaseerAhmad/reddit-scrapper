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

var posts []models.Post
var collector *colly.Collector
var pageScrapCounter int

func Init() {
	collector = colly.NewCollector(
		colly.AllowedDomains("old.reddit.com"),
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
}

func Start(pages int, url, fiName string) {
	if collector == nil {
		log.Fatal("[FATAL] Scrapper not initialized. Call to Init() required")
		return
	}

	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Scrapping ", request.URL)
	})

	collector.OnScraped(func(r *colly.Response) {
		fmt.Println("[SUCCESS] Scrapping done!")
	})

	collector.OnHTML("span.next-button", func(h *colly.HTMLElement) {
		if pageScrapCounter != pages {
			t := h.ChildAttr("a", "href")
			err := collector.Visit(t)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
		pageScrapCounter++
	})

	startCrawl(collector, url)
	collector.Wait()

	scrapLogger.LogToJson(posts, fiName)
}

func refineDomain(domain string) string {
	return domain[1 : len(domain)-1]
}

func startCrawl(collector *colly.Collector, url string) {
	err := collector.Visit(url)
	if err != nil {
		log.Fatal("[FATAL] Visit Error: ", err.Error())
		return
	}

	collector.OnHTML(".top-matter", func(e *colly.HTMLElement) {
		var post models.Post
		post.Title = e.ChildText("a[data-event-action=title]")
		post.Domain = refineDomain(e.ChildText("span[class=domain]"))
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
