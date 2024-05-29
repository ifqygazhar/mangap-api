package controller

import (
	"log"
	"mangap-api/entity"
	"mangap-api/util"
	"strconv"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
)

func FetchKomikByGenre(c *fiber.Ctx) error {
	page := c.Query("page")
	if page == "" {
		return c.Status(fiber.StatusBadRequest).SendString("page is required")
	}

	url := c.Params("url")
	komikList := []entity.LatestKomik{}

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	var currentPage, lengthPage string

	collector.OnHTML("#content > .wrapper > .postbody > .bixbox", func(e *colly.HTMLElement) {
		currentPage = e.ChildText(".listupd > .list-update_items > .pagination > .current")

		for i := 5; i <= 11; i++ {
			classAttr := e.ChildAttr(".pagination > .page-numbers:nth-child("+strconv.Itoa(i)+")", "class")
			if classAttr == "next page-numbers" {
				lengthPage = e.ChildText(".pagination > .page-numbers:nth-child(" + strconv.Itoa(i-1) + ")")
				break
			} else if classAttr == "page-numbers current" {
				lengthPage = e.ChildText(".pagination > .page-numbers:nth-child(" + strconv.Itoa(i) + ")")
				break
			}
		}

		e.ForEach(".listupd > .list-update_items > .list-update_items-wrapper > .list-update_item", func(_ int, el *colly.HTMLElement) {
			title := el.ChildText("a > .list-update_item-info > h3")
			chapter := el.ChildText("a > .list-update_item-info > .other > .chapter")
			typ := el.ChildText("a > .list-update_item-image > .type")
			thumbnail := el.ChildAttr("a > .list-update_item-image > img", "src")
			rating := el.ChildText("a > .list-update_item-info > .other > .rate > .rating > .numscore")
			href := el.ChildAttr("a", "href")

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

	err := collector.Visit(util.BASE_URL + "/genres/" + url + "/page/" + page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
	}

	currentPageInt, _ := strconv.Atoi(currentPage)
	lengthPageFloat, _ := strconv.ParseFloat(lengthPage, 64)

	return c.JSON(fiber.Map{
		"status":       "success",
		"current_page": currentPageInt,
		"length_page":  lengthPageFloat,
		"data":         komikList,
	})
}
