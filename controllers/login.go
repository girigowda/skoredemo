package controllers

import (
	"net/http"
	"sk-integrated-services/models"
	"sk-integrated-services/pkg/helpers"
	"sk-integrated-services/pkg/helpers/dbhelpers"
	"sk-integrated-services/pkg/logger"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (base *Controller) Login(ctx *gin.Context) {
	var requestBody helpers.LoginRequestBody

	//check if body is present in payload
	if err := ctx.BindJSON(&requestBody); err != nil {
		logger.Errorf("error: %v", err.Error())
		// ctx.String(http.StatusBadRequest, "Invalid BODY")
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid BODY",
		})
		return
	}

	//check the length of phone and password field
	if len(requestBody.Phone) <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Phone Number is empty",
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

	if len(requestBody.Password) > 6 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "password is empty or more than 6 character",
		})
		return
	}

	if len(requestBody.Password) <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "password is empty",
		})
		return
	}

	var results models.SL_USER
	//get data based on phone number
	k := dbhelpers.FilterQuery(base.DB, "SL_USER", map[string]interface{}{"phone_number": requestBody.Phone}, &results)
	println(k.RowsAffected)
	println(results.User_id)

	if k.RowsAffected != 0 {
		if err := bcrypt.CompareHashAndPassword([]byte(results.Pin), []byte(requestBody.Password)); err != nil {

			// save data in user_device
			user_device := &models.SL_User_Device_Details{User_id: results.User_id, Country: requestBody.Country, Ip_address: requestBody.Ip_address, Location: requestBody.Location, Udid: requestBody.Udid}
			user_device_result := base.DB.Create(&user_device)
			println(user_device_result.RowsAffected)

			// populate data in SL_login_Attempts table
			login_attempts_result := dbhelpers.InsertQuery(base.DB, "SL_Login_Attempts", map[string]interface{}{"user_id": results.User_id, "status": "Success", "device_id": user_device.Id})
			println(login_attempts_result.RowsAffected)

			// If the two passwords don't match, return a 401 status
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Password is incorrect",
			})
			return
		} else {
			jwt := helpers.GenerateJWT(requestBody.Phone)
			// save data in user_device
			user_device := &models.SL_User_Device_Details{User_id: results.User_id, Country: requestBody.Country, Ip_address: requestBody.Ip_address, Location: requestBody.Location, Udid: requestBody.Udid}
			user_device_result := base.DB.Create(&user_device)
			println(user_device_result.RowsAffected)

			// populate data in SL_login_Attempts table
			login_attempts_result := dbhelpers.InsertQuery(base.DB, "SL_Login_Attempts", map[string]interface{}{"user_id": results.User_id, "status": "Success", "device_id": user_device.Id, "login_timestamp": time.Now()})
			println(login_attempts_result.RowsAffected)

			//update active device id
			dbhelpers.UpdateQuery(base.DB, "SL_USER", map[string]interface{}{"phone_number": requestBody.Phone}, map[string]interface{}{"active_device_id": user_device.Id, "updated_at": time.Now()})
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": "Login Successfullly",
				"jwt":     jwt,
			})
		}
	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "User doesnot exist",
		})
		return
	}

}
