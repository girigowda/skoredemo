package controllers

import (
	"net/http"
	"sk-integrated-services/models"
	"strconv"
	"time"

	"sk-integrated-services/pkg/helpers"
	"sk-integrated-services/pkg/helpers/dbhelpers"
	"sk-integrated-services/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (base *Controller) Save_pin(ctx *gin.Context) {
	var requestBody helpers.SavepinRequestBody

	//check if body is present in payload
	if err := ctx.BindJSON(&requestBody); err != nil {
		logger.Errorf("error: %v", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid BODY",
		})
		return
	}

	//check if body has phone and password field
	if requestBody.Phone == "" && requestBody.Password == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Phone Number or password is not provided",
		})
		return
	}
	p, err := strconv.ParseUint(requestBody.Phone, 0, 64)
	println("p", p)
	if err != nil {
		println(err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Phone number should have only numbers",
		})
		return
	}

	p1, err := strconv.ParseUint(requestBody.Password, 0, 64)
	println("p", p1)
	if err != nil {
		println(err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Password number should have only numbers",
		})
		return
	}

	if len(requestBody.Phone) != 11 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Phone number length should be less than 11",
		})
		return
	}

	//check the length of phone and password field
	if len(requestBody.Phone) <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Phone Number or password is empty",
		})
		return
	}
	if len(requestBody.Password) > 6 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "password is empty or more than 6 character",
		})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), 8)
	if err != nil {
		logger.Errorf("error: %v", err.Error())
		// ctx.String(http.StatusBadRequest, "Invalid BODY")
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Error in password creation",
		})
		return
	}

	//get user details from user table
	var user_results models.SL_USER
	user_query_result := dbhelpers.FilterQuery(base.DB, "SL_USER", map[string]interface{}{"phone_number": requestBody.Phone}, &user_results)

	println(user_query_result.RowsAffected)
	if user_query_result.RowsAffected > 0 && user_results.Otp_verified_status {
		dbhelpers.UpdateQuery(base.DB, "SL_USER", map[string]interface{}{"phone_number": requestBody.Phone}, map[string]interface{}{"pin": hashedPassword, "pin_changed_on": time.Now(), "pin_updated_at": time.Now(), "pin_activated": true, "updated_at": time.Now()})

		// populate data in SL_PIN_Attempts table
		pin_attempts_result := dbhelpers.InsertQuery(base.DB, "SL_PIN_Attempts", map[string]interface{}{"user_id": user_results.User_id, "status": "Active", "pin_generated": time.Now()})
		println(pin_attempts_result.RowsAffected)

		// populate data in SL_login_Attempts table
		login_attempts_result := dbhelpers.InsertQuery(base.DB, "SL_Login_Attempts", map[string]interface{}{"user_id": user_results.User_id, "status": "Success", "device_id": user_results.Active_device_id, "login_timestamp": time.Now()})
		println(login_attempts_result.RowsAffected)

		jwt := helpers.GenerateJWT(requestBody.Phone)

		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "Login Successfullly",
			"jwt":     jwt,
		})
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "OTP not verified",
		})
	}

}
