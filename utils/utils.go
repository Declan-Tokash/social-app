package utils

import (
	"github.com/Declan-Tokash/social-api/database"
	"github.com/Declan-Tokash/social-api/model"
	"math"
)



// Haversine formula to calculate the great-circle distance between two points
func haversine(lon1, lat1, lon2, lat2 float64) float64 {
    const R = 6371 // Earth radius in kilometers

    dLat := (lat2 - lat1) * math.Pi / 180
    dLon := (lon2 - lon1) * math.Pi / 180
    lat1 = lat1 * math.Pi / 180
    lat2 = lat2 * math.Pi / 180

    a := math.Sin(dLat/2)*math.Sin(dLat/2) +
        math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    
    return R * c // Distance in kilometers
}

func GetUserLocation(userID string) (float64, float64) {
	db := database.DB.Db
	location := new(model.Location)

	// Use First if you expect only one result
	result := db.Where("user_id = ?", userID).First(location)
	if result.Error != nil {
		// Return default values if there's an error
		return -1, -1
	}

	// Return latitude and longitude values as strings
	return location.Latitude, location.Longitude
}