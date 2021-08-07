package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"log"
	"time"
)

type post struct {
	Url      string `json:"url"`
	Title    string `json:"title"`
	Domain   string `json:"domain"`
	Author   string `json:"author"`
	PostedOn string `json:"postedOn"`
}

var posts []post
var requestsMade int
var refreshInterval int

func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("old.reddit.com"),
		colly.AllowURLRevisit(),
		colly.Async(true))

	err := c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 1,
		RandomDelay: 5 * time.Second,
	})
	if err != nil {
		log.Fatal("[FATAL] Limit Error: ", err)
		return
	}

	c.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL.String())
	})

	ticker := time.NewTicker(time.Second)

	c.OnRequest(func(request *colly.Request) {
		requestsMade++

		ticker.Reset(5 * time.Second)
	})

	for {
		select {
		case <- ticker.C:
			startCrawl(c)
			c.Wait()

			indent, err := json.MarshalIndent(posts, "", "  ")
			if err != nil {
				log.Fatal(err.Error())
				return
			}
			fmt.Println(string(indent))
			fmt.Println("TOTAL REQUESTS MADE: ", requestsMade)
		}
	}

}

func refineDomain(domain string) string {
	return domain[1 : len(domain)-1]
}

func startCrawl(collector *colly.Collector) {
	err := collector.Visit("https://old.reddit.com/r/worldnews/new")
	if err != nil {
		log.Fatal("[FATAL] Visit Error: ", err.Error())
		return
	}

	collector.OnHTML(".top-matter", func(e *colly.HTMLElement) {
		var post post
		post.Title = e.ChildText("a[data-event-action=title]")
		post.Domain = refineDomain(e.ChildText("span[class=domain]"))
		post.Url = e.ChildAttr("a[data-event-action=title]", "href")
		//post.comments = e.ChildAttr("a[data-event-action=comments]", "href")
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
