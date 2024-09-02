package main
import (
 "github.com/Declan-Tokash/social-api/database"
 "github.com/Declan-Tokash/social-api/aws"
 "github.com/Declan-Tokash/social-api/router"
 "github.com/gofiber/fiber/v2"
 "github.com/gofiber/fiber/v2/middleware/cors"
 "github.com/gofiber/fiber/v2/middleware/logger"
 _ "github.com/lib/pq"
)
func main() {
 database.Connect()
 aws.SetUpS3Uploader()
 app := fiber.New()
 app.Use(logger.New())
 app.Use(cors.New())
 router.SetupRoutes(app)
 // handle unavailable route
 app.Use(func(c *fiber.Ctx) error {
  return c.SendStatus(404) // => 404 "Not Found"
 })
 app.Listen(":8080")
}