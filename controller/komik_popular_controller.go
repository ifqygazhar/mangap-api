package controller

import (
	"crypto/tls"
	"log"
	"mangap-api/entity"
	"mangap-api/util"
	"net/http"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
)

func FetchPopularKomik(c *fiber.Ctx) error {
	komikList := []entity.Komik{}

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
	)
	collector.WithTransport(&http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	collector.OnHTML(
		"#content > .wrapper > #sidebar .section > .widget-post > .serieslist.pop > ul > li",
		func(e *colly.HTMLElement) {
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
		},
	)

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
