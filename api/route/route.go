package route

import (
	"github.com/gin-gonic/gin"
	"tripservice/api/handler"
)

func Routing(router *gin.Engine) {
	router.POST("/trips", handler.AddTrips)
	router.GET("/trips/:id/postcodes", handler.GetTripPostcodes)
	router.GET("/trips/:id/speeds", handler.GetTripSpeeds)
	router.GET("/vehicles/:id/trips", handler.GetVehicleTrips)
}
