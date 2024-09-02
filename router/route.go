package router
import (
 "github.com/Declan-Tokash/social-api/handler"
 "github.com/Declan-Tokash/social-api/middleware"
 "github.com/gofiber/fiber/v2"
)
// SetupRoutes func
func SetupRoutes(app *fiber.App) {
 // grouping
 api := app.Group("/api")

 //auth routes
 auth := api.Group("/auth")
 auth.Post("/login", handler.Login)
 
 // user routes
 v1 := api.Group("/user")
 v1.Get("/", handler.GetAllUsers)
 v1.Get("/:id", middleware.Protected(), handler.GetSingleUser)
 v1.Post("/", handler.CreateUser)
 v1.Put("/:id", middleware.Protected(), handler.UpdateUser)
 v1.Delete("/:id", middleware.Protected(), handler.DeleteUserByID)

 //post routes
 post := api.Group("/post")
 post.Post("/", middleware.Protected(), middleware.ExtractUserID(), handler.CreatePost)
 post.Get("/:id", middleware.Protected(), handler.GetUserPosts)
 post.Get("/", middleware.Protected(), middleware.ExtractUserID(), handler.GetPostsInArea)

 //location
 location := api.Group("/location")
 location.Post("/", middleware.Protected(), middleware.ExtractUserID(), handler.UpdateUserLocation)
 location.Get("/:id", middleware.Protected(), middleware.ExtractUserID(), handler.GetUserLocation)

}