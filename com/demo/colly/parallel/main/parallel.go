package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main(){
	c :=colly.NewCollector(
		colly.MaxDepth(2),
		colly.Async(true),
		)
	c.Limit(&colly.LimitRule{DomainGlob:"*",Parallelism:2})

	c.OnHTML("a[href]", func(element *colly.HTMLElement) {
		link:=element.Attr("href")
		fmt.Println(link)
		element.Request.Visit(link)
	})

	c.Visit("https://en.wikipedia.org")
	c.Wait()
}
