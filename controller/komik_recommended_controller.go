package controller

import (
	"log"
	"mangap-api/entity"

	"github.com/gocolly/colly"
	"github.com/gofiber/fiber/v2"
)

func FetchRecommendedKomik(c *fiber.Ctx) error {
	baseUrl := "https://komikcast.lol"
	recommendedList := []entity.RecommendedKomik{}

	collector := colly.NewCollector(
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"),
	)

	collector.OnHTML("#content > .wrapper > .bixbox > .listupd > .swiper > .swiper-wrapper > .swiper-slide", func(e *colly.HTMLElement) {
		title := e.ChildText("a > .splide__slide-info > .title")
		rating := e.ChildText("a > .splide__slide-info > .other > .rate > .rating > .numscore")
		chapter := e.ChildText("a > .splide__slide-info > .other > .chapter")
		typ := e.ChildText("a > .splide__slide-image > .type")
		href := e.ChildAttr("a", "href")
		thumbnail := e.ChildAttr("a > .splide__slide-image > img", "src")

		recommended := entity.RecommendedKomik{
			Title:     title,
			Href:      href,
			Rating:    rating,
			Chapter:   chapter,
			Type:      typ,
			Thumbnail: thumbnail,
		}
		recommendedList = append(recommendedList, recommended)
	})

	collector.OnError(func(_ *colly.Response, err error) {
		log.Println("Request error:", err)
		c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
	})

	err := collector.Visit(baseUrl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data")
	}

	filteredRecommendedList := []entity.RecommendedKomik{}
	for _, v := range recommendedList {
		if v.Href != "" {
			filteredRecommendedList = append(filteredRecommendedList, v)
		}
	}

	return c.JSON(filteredRecommendedList)
}
