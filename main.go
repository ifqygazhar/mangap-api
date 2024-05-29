package main

import (
	"log"
	"mangap-api/controller"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("view/index.html")
	})

	app.Get("/popular", controller.FetchPopularKomik)
	app.Get("/recommended", controller.FetchRecommendedKomik)
	app.Get("/terbaru", controller.FetchLatestKomik)
	app.Get("/genre/:url", controller.FetchKomikByGenre)
	app.Get("/genre", controller.FetchAllGenres)
	app.Get("/read/:url", controller.FetchKomikDetails)
	app.Get("/search", controller.SearchKomik)
	app.Get("/detail/:url", controller.FetchKomikDetail)

	log.Println("Starting server on :8080")
	log.Fatal(app.Listen(":8080"))
}
