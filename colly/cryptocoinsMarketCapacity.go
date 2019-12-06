package main

import (
	"encoding/csv"
	"github.com/gocolly/colly"
	"log"
	"os"
)

func main() {
	fName := "cryptocoinmarketcap.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("Cannot create file %q: %s\n", fName, err)
	}

	defer file.Close()

	write := csv.NewWriter(file)
	defer write.Flush()

	write.Write([]string{"Name", "Symbol", "Price(USD)", "Volume(USD)", "Market capacity (USD)", "Change (1h)", "Change (24h)", "Change (7d)"})

	c := colly.NewCollector()

	c.OnHTML("tbody tr", func(e *colly.HTMLElement) {
		write.Write([]string{
			e.ChildText(".cmc-table__cell--sort-by__name a"),
			e.ChildText(".cmc-table__cell--sort-by__symbol div"),
			e.ChildText(".cmc-table__cell--sort-by__price a"),
			e.ChildText(".cmc-table__cell--sort-by__volume-24-h a"),
			e.ChildText(".cmc-table__cell--sort-by__market-cap div"),
			e.ChildText(".cmc-table__cell--sort-by__percent-change-1-h div"),
			e.ChildText(".cmc-table__cell--sort-by__percent-change-24-h div"),
			e.ChildText(".cmc-table__cell--sort-by__percent-change-7-d div"),
		})
	})

	c.Visit("https://coinmarketcap.com/all/views/all/")

	log.Printf("Scraping finished,check file %q fro results\n", fName)

}
