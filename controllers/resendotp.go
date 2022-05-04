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
)

func (base *Controller) Resend_otp(ctx *gin.Context) {
	var requestBody helpers.OTPResendBody
	//check if body is present in payload
	if err := ctx.BindJSON(&requestBody); err != nil {
		logger.Errorf("error: %v", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid BODY",
		})
		return
	}

	//check if body has phone
	if requestBody.Phone == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Phone Number is not provided",
		})
		return
	}

	//check the length of phone
	if len(requestBody.Phone) == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Phone Number is empty",
		})
		return
	}

	// GENERATE OTP
	b := helpers.Generatenumber(6)
	value := string(b)
	number, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		logger.Errorf("error: %v", err.Error())
	}
	finalIntNum := int(number)

	//get user details from user table
	var user_results models.SL_USER
	user_query_result := dbhelpers.FilterQuery(base.DB, "SL_USER", map[string]interface{}{"phone_number": requestBody.Phone}, &user_results)
	println(user_query_result.RowsAffected)

	//get data from otp_attempts for particular user
	var otp_results []models.SL_OTP_Attempts
	otp_query_results := dbhelpers.FilterQuery(base.DB, "SL_OTP_Attempts", map[string]interface{}{"phone_number": requestBody.Phone}, &otp_results)
	println(otp_query_results.RowsAffected)

	//check condition for number of attempts in a day
	var total_count = 0
	for _, product := range otp_results {
		if product.Otp_generated.Day() == time.Now().Day() {
			total_count += 1
		}
		dbhelpers.UpdateQuery(base.DB, "SL_OTP_Attempts", map[string]interface{}{"id": product.Id}, map[string]interface{}{"status": "invalid", "active_otp": finalIntNum})
	}

	//if already 4 attempts throught error
	otp_count := helpers.GetEnvVariable("OTP_ATTEMPTS")
	parsed_otp_count, err := strconv.ParseInt(otp_count, 0, 32)
	if err != nil {
		logger.Errorf("error: %v", err.Error())
	}
	if total_count >= int(parsed_otp_count) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "4 attempt completed",
		})
		return
	}

	t := time.Now()
	otp_time := helpers.GetEnvVariable("OTP_TIME")
	parsed_otp_time, err := strconv.Atoi(otp_time)
	if err != nil {
		logger.Errorf("error: %v", err.Error())
	}
	newT := t.Add(time.Second * time.Duration(parsed_otp_time))

	//populate data in SL_OTP_Attempts table
	otp_result := dbhelpers.InsertQuery(base.DB, "SL_OTP_Attempts", map[string]interface{}{"user_id": user_results.User_id, "otp": finalIntNum, "otp_generated": time.Now(), "active_otp": finalIntNum, "status": "Valid", "phone_number": requestBody.Phone, "otp_expiry": newT})
	println(otp_result.RowsAffected)

	ctx.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "OTP sent",
		"otp":     finalIntNum,
	})
}
