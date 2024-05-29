package controller

import (
	"log"
	"mangap-api/entity"
	"mangap-api/util"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
)

func FetchPopularKomik(c *fiber.Ctx) error {
	komikList := []entity.Komik{}

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	collector.OnHTML("#content > .wrapper > #sidebar .section > .widget-post > .serieslist.pop > ul > li", func(e *colly.HTMLElement) {
		title := e.ChildText(".leftseries > h2 > a")
		year := e.ChildText(".leftseries > span:nth-child(3)")
		genre := e.ChildText(".leftseries > span:nth-child(2)")
		thumbnail := e.ChildAttr(".imgseries > a > img", "src")
		href := e.ChildAttr(".imgseries > a", "href")

		komik := entity.Komik{
			Title:     title,
			Href:      href,
			Genre:     genre,
			Year:      year,
			Thumbnail: thumbnail,
		}
		komikList = append(komikList, komik)
	})

	collector.OnError(func(_ *colly.Response, err error) {
		log.Println("Request error:", err)
		c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
	})

	err := collector.Visit(util.BASE_URL)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
	}

	return c.JSON(komikList)
}
