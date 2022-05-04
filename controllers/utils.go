package controllers

import (
	"net/http"
	"sk-integrated-services/models"
	"sk-integrated-services/pkg/helpers/dbhelpers"

	"github.com/gin-gonic/gin"
)

func (base *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "ping")
}

func (base *Controller) Createquery(ctx *gin.Context) {
	example := "Tables are created"
	User := new(models.SL_USER)
	SL_OTP_Attempts := new(models.SL_OTP_Attempts)
	SL_Login_Attempts := new(models.SL_Login_Attempts)
	SL_PIN_Attempts := new(models.SL_PIN_Attempts)
	SL_User_Device_Details := new(models.SL_User_Device_Details)
	SL_KYC_Details := new(models.SL_KYC_Details)
	SL_KYC_facematch := new(models.SL_KYC_facematch)
	SL_numbers := new(models.SL_numbers)

	base.DB.Migrator().CreateTable(&User)

	base.DB.Migrator().CreateTable(&SL_OTP_Attempts)

	base.DB.Migrator().CreateTable(&SL_Login_Attempts)

	base.DB.Migrator().CreateTable(&SL_PIN_Attempts)

	base.DB.Migrator().CreateTable(&SL_User_Device_Details)

	base.DB.Migrator().CreateTable(&SL_KYC_Details)

	base.DB.Migrator().CreateTable(&SL_KYC_facematch)

	base.DB.Migrator().CreateTable(&SL_numbers)

	k := dbhelpers.InsertQuery(base.DB, "SL_numbers", map[string]interface{}{"count": 1, "id": 1})
	println(k.RowsAffected)

	ctx.JSON(http.StatusOK, &example)
}
