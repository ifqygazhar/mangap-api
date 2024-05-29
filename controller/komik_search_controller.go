package controller

import (
	"log"
	"mangap-api/entity"
	"mangap-api/util"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
)

func SearchKomik(c *fiber.Ctx) error {
	keyword := c.Query("keyword")
	if keyword == "" {
		return c.Status(fiber.StatusBadRequest).SendString("keyword is required")
	}

	komikList := []entity.LatestKomik{}

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
	)

	collector.OnHTML("#content > .wrapper > .postbody > .dev > #main > .list-update", func(e *colly.HTMLElement) {
		e.ForEach(".list-update_items > .list-update_items-wrapper > .list-update_item", func(_ int, el *colly.HTMLElement) {
			title := el.ChildText("a > .list-update_item-info > h3")
			href := el.ChildAttr("a", "href")
			typ := el.ChildText("a > .list-update_item-image > .type")
			rating := el.ChildText("a > .list-update_item-info > .other > .rate > .rating > .numscore")
			chapter := el.ChildText("a > .list-update_item-info > .other > .chapter")
			thumbnail := el.ChildAttr("a > .list-update_item-image > img", "src")

			komik := entity.LatestKomik{
				Title:     title,
				Href:      href,
				Thumbnail: thumbnail,
				Type:      typ,
				Chapter:   chapter,
				Rating:    rating,
			}
			komikList = append(komikList, komik)
		})
	})

	collector.OnError(func(_ *colly.Response, err error) {
		log.Println("Request error:", err)
		c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
	})

	err := collector.Visit(util.BASE_URL + "/?s=" + keyword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
	}

	return c.JSON(komikList)
}
