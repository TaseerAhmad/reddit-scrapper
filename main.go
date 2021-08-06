package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"log"
	"time"
)

type post struct {
	comments 		 string
	url				 string
	title 			 string
	source 			 string
	domain           string
	author           string
	timeOfSubmission string
}

func main() {
	var posts []post

	collector := colly.NewCollector(
			colly.AllowedDomains("old.reddit.com"),
			colly.Async(true))

	collector.OnHTML(".top-matter", func(e *colly.HTMLElement) {
		var post post
		post.title = e.ChildText("a[data-event-action=title]")
		post.domain = refineDomain(e.ChildText("span[class=domain]"))
		post.url = e.ChildAttr("a[data-event-action=title]", "href")
		//post.comments = e.ChildAttr("a[data-event-action=comments]", "href")
		selector := e.DOM.Find("p")
		if selector.HasClass("tagline") {
			selector.Children().Each(func(i int, selection *goquery.Selection) {
				if dateTime, exists := e.DOM.Find("time").Attr("datetime"); exists {
					post.timeOfSubmission = dateTime
				}

				if selection.HasClass("author may-blank") {
					post.author = selection.Contents().Text()
				}
			})
		}

		posts = append(posts, post)
	})

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

	err = collector.Visit("https://old.reddit.com/r/worldnews/new")
	if err != nil {
		log.Fatal("[FATAL] Visit Error: ", err.Error())
		return
	}

	collector.Wait()

	fmt.Println(posts)
}

func refineDomain(domain string) string {
	return domain[1: len(domain) - 1]
}
