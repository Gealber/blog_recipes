package utils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func ScrapeFragmentMarket() ([]AnonNumbers, error) {
	numbers := make([]AnonNumbers, 0)
	c := colly.NewCollector()

	c.OnHTML("tr td.thin-last-col a", func(e *colly.HTMLElement) {
		name := strings.Trim(e.Attr("href"), "/number/")
		valueStr := e.ChildText("div.table-cell-value.tm-value.icon-before.icon-ton")
		value, err := strconv.Atoi(valueStr)
		if err != nil {
			return
		}		

		numbers = append(numbers, AnonNumbers{Name: name, Value: value})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit("https://fragment.com/numbers?sort=price_asc&filter=sale")
	if err != nil {
		return nil, err
	}

	return numbers, nil
}
