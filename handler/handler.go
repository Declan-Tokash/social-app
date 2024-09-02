package handler
import (
 "github.com/Declan-Tokash/social-api/database"
 "github.com/Declan-Tokash/social-api/model"
 "github.com/Declan-Tokash/social-api/aws"
 "github.com/Declan-Tokash/social-api/utils"
 "github.com/gofiber/fiber/v2"
 "github.com/google/uuid"
 "github.com/golang-jwt/jwt/v5"
 "time"
 "log"
)

//Create a user
func CreateUser(c *fiber.Ctx) error {
	db := database.DB.Db
	user := new(model.User)
   // Store the body in the user and return error if encountered
	err := c.BodyParser(user)
	if err != nil {
	 return c.Status(500).JSON(fiber.Map{"status": "error", "message":  "Something's wrong with your input", "data": err})
	}
    err = db.Create(&user).Error
	if err != nil {
	 return c.Status(500).JSON(fiber.Map{"status": "error", "message":  "Could not create user", "data": err})
	}
	location := new(model.Location)
	location.UserID = user.ID.String()
	location.Latitude = -1
	location.Longitude = -1
	err = db.Create(&location).Error
	if err != nil {
	 return c.Status(500).JSON(fiber.Map{"status": "error", "message":  "Could not create user", "data": err})
	}
   // Return the created user
	return c.Status(201).JSON(fiber.Map{"status": "success", "message":  "User has created", "data": user})
}

func Login(c *fiber.Ctx) error {
	db := database.DB.Db

	userLogin := new(model.UserLogin)
	if err := c.BodyParser(&userLogin); err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user := new(model.User)

	// Check if the user exists with the given identity (username or email)
    if err := db.Where("username = ?", userLogin.Username).First(&user).Error; err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid credentials"})
    }

    // Check if the password is correct
    if userLogin.Password != user.Password { // Implement password hash checking
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid credentials"})
    }

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["identity"] = user.ID
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Success login", "data": t})
}

// Get All Users from db
func GetAllUsers(c *fiber.Ctx) error {
	db := database.DB.Db
	var users []model.User
   // find all users in the database
	db.Find(&users)
   // If no user found, return an error
	if len(users) == 0 {
	 return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Users not found", "data": nil})
	}
   // return users
	return c.Status(200).JSON(fiber.Map{"status": "sucess", "message": "Users Found", "data": users})
}

// GetSingleUser from db
func GetSingleUser(c *fiber.Ctx) error {
	db := database.DB.Db
   // get id params
	id := c.Params("id")
   var user model.User
   // find single user in the database by id
	db.Find(&user, "id = ?", id)
   if user.ID == uuid.Nil {
	 return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}
   return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User Found", "data": user})
}

   // update a user in db
func UpdateUser(c *fiber.Ctx) error {
	type updateUser struct {
	 Username string `json:"username"`
	}
   db := database.DB.Db
   var user model.User
   // get id params
	id := c.Params("id")
   // find single user in the database by id
	db.Find(&user, "id = ?", id)
   if user.ID == uuid.Nil {
	 return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
	}
   var updateUserData updateUser
	err := c.BodyParser(&updateUserData)
	if err != nil {
	 return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input", "data": err})
	}
   user.Username = updateUserData.Username
   // Save the Changes
	db.Save(&user)
   // Return the updated user
	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "users Found", "data": user})
}

   // delete user in db by ID
func DeleteUserByID(c *fiber.Ctx) error {
	db := database.DB.Db
	var user model.User
   // get id params
	id := c.Params("id")
   // find single user in the database by id
	db.Find(&user, "id = ?", id)
   if user.ID == uuid.Nil {
	 return c.Status(404).JSON(fiber.Map{"status": "error", "message": "User not found", "data": nil})
   }
   err := db.Delete(&user, "id = ?", id).Error
   if err != nil {
	 return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Failed to delete user", "data": nil})
	}
   return c.Status(200).JSON(fiber.Map{"status": "success", "message": "User deleted"})
}

func UpdateUserLocation(c *fiber.Ctx) error {
	db := database.DB.Db

    userID := c.Locals("userID").(string)
    location := new(model.Location)
    
    if err := c.BodyParser(location); err != nil {
        return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Invalid location data", "data": err.Error()})
    }

    // Update the location based on user ID
    err := db.Model(&model.Location{}).Where("user_id = ?", userID).Updates(location).Error
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not update location", "data": err.Error()})
    }

    return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Location updated", "data": location})
}

// GetUserLocation from db
func GetUserLocation(c *fiber.Ctx) error {
	db := database.DB.Db
   // get id params
	id := c.Params("id")
    var location model.Location
   // find single user in the database by id
   db.Find(&location, "user_id = ?", id)
   if location.UserID == "" {
	 return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Location not found", "data": nil})
	}
   return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Location Found", "data": location})
}

func CreatePost(c *fiber.Ctx) error {
	//Load file
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Failed to upload image", "data": err.Error()})
	}
	// Upload the file
	result, err := aws.UploadFile(file)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "error uploading file", "data": err.Error()})
	}

	db := database.DB.Db
	userID := c.Locals("userID").(string)
	latitude, longitude := utils.GetUserLocation(userID)

	if latitude == -1 && longitude == -1 {
		log.Println("Error: Invalid location data")
		return c.Status(500).SendString("Error getting user location")
	}

	post := new(model.Post)
	location := model.Location{
		UserID: userID,
		Latitude:  latitude,
		Longitude: longitude,
	}

	post.Title = "Post"
	post.Image = result.Location
	post.Location = location

	err = db.Create(&post).Error
	if err != nil {
	 return c.Status(500).JSON(fiber.Map{"status": "error", "message":  "Could not create user", "data": err})
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "message":  "Image has been posted", "file": result})
}

//GetUserPosts
func GetUserPosts(c *fiber.Ctx) error {
	db := database.DB.Db
    // get id params
	id := c.Params("id")
    var posts []model.Post
    // find posts in the database by user_id
    db.Find(&posts, "user_id = ?", id)
    if len(posts) == 0 {
	  return c.Status(404).JSON(fiber.Map{"status": "error", "message": "Posts not found", "data": nil})
	 }
    return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Posts Found", "data": posts})
}

//GetPostsInArea
func GetPostsInArea(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)
	latitude, longitude := utils.GetUserLocation(userID)

	if latitude == -1 && longitude == -1 {
		log.Println("Error: Invalid location data")
		return c.Status(500).SendString("Error getting user location")
	}

	db := database.DB.Db
    radiusKm := 10 * 1.60934

	query := `
        SELECT * FROM posts
        WHERE (
            6371 * acos(
                cos(radians(?)) * cos(radians(latitude)) * cos(radians(longitude) - radians(?)) +
                sin(radians(?)) * sin(radians(latitude))
            )
        ) < ?
    `
    
    var posts []model.Post
    err := db.Raw(query, latitude, longitude, latitude, radiusKm).Scan(&posts).Error
    if err != nil {
        return err
    }

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "Posts Found", "data": posts})

}