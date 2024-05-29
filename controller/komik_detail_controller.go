package controller

import (
	"log"
	"mangap-api/entity"
	"mangap-api/util"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
)

func FetchKomikDetails(c *fiber.Ctx) error {
	url := c.Params("url")
	readList := []entity.ReadKomik{}

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
	)

	collector.OnHTML("#content > .wrapper", func(e *colly.HTMLElement) {
		title := e.ChildText(".chapter_headpost > h1")
		panels := []string{}

		e.ForEach(".chapter_ > #chapter_body > .main-reading-area > img", func(_ int, el *colly.HTMLElement) {
			panels = append(panels, el.Attr("src"))
		})

		readKomik := entity.ReadKomik{
			Title: title,
			Panel: panels,
		}
		readList = append(readList, readKomik)
	})

	collector.OnError(func(res *colly.Response, err error) {
		if res != nil && res.StatusCode == 404 {
			c.Status(fiber.StatusNotFound).SendString("comic not found")
		} else {
			log.Println("Request error:", err)
			c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
		}
	})

	err := collector.Visit(util.BASE_URL + "/" + url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
	}

	return c.JSON(readList)
}

func FetchKomikDetail(c *fiber.Ctx) error {
	url := c.Params("url")
	komikList := []entity.DetailedKomik{}

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0"),
	)

	collector.OnHTML("#content > .wrapper > .komik_info", func(e *colly.HTMLElement) {
		var komik entity.DetailedKomik
		komik.Title = e.ChildText(".komik_info-body > .komik_info-content > .komik_info-content-body > h1")
		komik.AltTitle = e.ChildText(".komik_info-body > .komik_info-content > .komik_info-content-body > .komik_info-content-native")
		komik.Released = e.ChildText(".komik_info-body > .komik_info-content > .komik_info-content-body > .komik_info-content-meta > span:nth-child(1)")
		komik.Author = e.ChildText(".komik_info-body > .komik_info-content > .komik_info-content-body > .komik_info-content-meta > span:nth-child(2)")
		komik.Status = e.ChildText(".komik_info-body > .komik_info-content > .komik_info-content-body > .komik_info-content-meta > span:nth-child(3)")
		komik.Type = e.ChildText(".komik_info-body > .komik_info-content > .komik_info-content-body > .komik_info-content-meta > span:nth-child(4)")
		komik.Description = e.ChildText(".komik_info-description > .komik_info-description-sinopsis > p")
		komik.Thumbnail = e.ChildAttr(".komik_info-cover-box > .komik_info-cover-image > img", "src")
		komik.UpdatedOn = e.ChildText(".komik_info-body > .komik_info-content > .komik_info-content-body > .komik_info-content-meta > .komik_info-content-update")

		e.ForEach(".komik_info-body > .komik_info-chapters > ul > li", func(_ int, el *colly.HTMLElement) {
			var chap entity.Chapter
			chap.Title = el.ChildText("a")
			chap.Href = el.ChildAttr("a", "href")
			chap.Date = el.ChildText(".chapter-link-time")
			komik.Chapter = append(komik.Chapter, chap)
		})

		e.ForEach(".komik_info-body > .komik_info-content > .komik_info-content-body > .komik_info-content-genre > a", func(_ int, el *colly.HTMLElement) {
			var gen entity.Genre
			gen.Title = el.Text
			gen.Href = el.Attr("href")
			komik.Genre = append(komik.Genre, gen)
		})

		komikList = append(komikList, komik)
	})

	collector.OnError(func(res *colly.Response, err error) {
		if res != nil && res.StatusCode == 404 {
			c.Status(fiber.StatusNotFound).SendString("comic not found")
		} else {
			log.Println("Request error:", err)
			c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
		}
	})

	err := collector.Visit(util.BASE_URL + "/komik/" + url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
	}

	return c.JSON(komikList)
}
