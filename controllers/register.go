package controllers

import (
	"net/http"
	"sk-integrated-services/models"
	"sk-integrated-services/pkg/helpers"
	"sk-integrated-services/pkg/helpers/dbhelpers"
	"sk-integrated-services/pkg/logger"
	"time"

	"strconv"

	"github.com/gin-gonic/gin"
)

func (base *Controller) Register(ctx *gin.Context) {
	var requestBody helpers.RegisterRequestBody
	//check if body is present in payload
	if err := ctx.BindJSON(&requestBody); err != nil {
		logger.Errorf("error:", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid Body",
		})
		return
	}

	//check the length of phone field
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

	if len(requestBody.Phone) != 11 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Phone number length should be 11",
		})
		return
	}

	//check the length of phone field
	if !requestBody.Privacy {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Privacy is not accepted",
		})
		return
	}

	//check whether phone number is already present
	var results models.SL_USER
	k := dbhelpers.FilterQuery(base.DB, "SL_USER", map[string]interface{}{"phone_number": requestBody.Phone}, &results)
	println(k.RowsAffected)
	if k.RowsAffected > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Phone Number already exist",
		})
		return
	}

	// GENERATE OTP
	otp_count := helpers.GetEnvVariable("OTP_COUNT")
	parsed_otp_count, err := strconv.ParseInt(otp_count, 0, 32)
	println(err)
	b := helpers.Generatenumber(int(parsed_otp_count))
	value := string(b)
	number, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		println(err)
	}
	finalIntNum := int(number)

	// Generate user_id
	user_id := helpers.Generateid(base.DB)

	//populate data in otp table
	t := time.Now()
	otp_time := helpers.GetEnvVariable("OTP_TIME")
	parsed_otp_time, err := strconv.Atoi(otp_time)
	println(err)
	newT := t.Add(time.Second * time.Duration(parsed_otp_time))
	otp_result := dbhelpers.InsertQuery(base.DB, "SL_OTP_Attempts", map[string]interface{}{"user_id": user_id, "otp": finalIntNum, "otp_generated": time.Now(), "active_otp": finalIntNum, "status": "Valid", "phone_number": requestBody.Phone, "otp_expiry": newT})
	println(otp_result.RowsAffected)

	// save data in user_device
	user_device := &models.SL_User_Device_Details{User_id: user_id, Country: requestBody.Country, Ip_address: requestBody.Ip_address, Location: requestBody.Location, Udid: requestBody.Udid}
	user_device_result := base.DB.Create(&user_device)
	println(user_device_result.RowsAffected)

	// save data in user table
	user_result := dbhelpers.InsertQuery(base.DB, "SL_USER", map[string]interface{}{"phone_number": requestBody.Phone, "created_at": time.Now(), "updated_at": time.Now(), "user_id": user_id, "active_device_id": user_device.Id, "privacy_policy": requestBody.Privacy})
	println(user_result.RowsAffected)

	if user_result.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Error in user creation",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":         http.StatusOK,
		"message":      "Registered Successfully",
		"phone_number": requestBody.Phone,
		"otp":          finalIntNum,
	})
}
