package routers

import (
	"sk-integrated-services/controllers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routings(route *gin.Engine, db *gorm.DB) {
	ctrl := controllers.Controller{DB: db}
	v1 := route.Group("/v1")
	v1.GET("/test", ctrl.Createquery)
	v1.POST("/register", ctrl.Register)
	v1.POST("/login", ctrl.Login)
	v1.POST("/otp_validate", ctrl.OTP_Validate)
	v1.POST("/otp", ctrl.Resend_otp)
	v1.POST("/pin", ctrl.Save_pin)
	v1.POST("/face_match", ctrl.Validate_face)
	v1.POST("/ocr_card_details", ctrl.Validate_card)
	v1.GET("/ping", ctrl.Ping)

}
