package controller

import (
	"log"
	"mangap-api/entity"
	"mangap-api/util"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
)

func FetchKomikList(c *fiber.Ctx) error {
	page := c.Params("page")
	var url string
	if page == "1" {
		url = util.BASE_URL + "/daftar-komik/"
	} else {
		url = util.BASE_URL + "/daftar-komik/page/" + page
	}

	komikList := []entity.ListOfKomik{}

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
	)

	collector.OnHTML("#content > .wrapper > .komiklist > .komiklist_filter > .list-update > .list-update_items > .list-update_items-wrapper > .list-update_item", func(e *colly.HTMLElement) {
		title := e.ChildText("a > .list-update_item-info > h3")
		chapter := e.ChildText("a > .list-update_item-info > .other > .chapter")
		typ := e.ChildText("a > .list-update_item-image > .type")
		thumbnail := e.ChildAttr("a > .list-update_item-image > img", "src")
		rating := e.ChildText("a > .list-update_item-info > .other > .rate > .rating > .numscore")
		href := e.ChildAttr("a", "href")

		komik := entity.ListOfKomik{
			Title:     title,
			Chapter:   chapter,
			Type:      typ,
			Thumbnail: thumbnail,
			Rating:    rating,
			Href:      href,
		}
		komikList = append(komikList, komik)
	})

	collector.OnError(func(_ *colly.Response, err error) {
		log.Println("Request error:", err)
		c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
	})

	err := collector.Visit(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
	}

	filteredKomikList := []entity.ListOfKomik{}
	for _, komik := range komikList {
		if komik.Href != "" {
			komik.TotalData = len(komikList)
			filteredKomikList = append(filteredKomikList, komik)
		}
	}

	return c.JSON(filteredKomikList)
}
