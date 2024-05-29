package controller

import (
	"log"
	"mangap-api/entity"
	"mangap-api/util"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
)

func FetchAllGenres(c *fiber.Ctx) error {
	komikList := []entity.Komik{}

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
	)

	collector.OnHTML("#content > .wrapper > #sidebar > .section > ul.genre > li", func(e *colly.HTMLElement) {
		title := e.ChildText("a")
		href := e.ChildAttr("a", "href")

		komik := entity.Komik{
			Title: title,
			Href:  href,
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
