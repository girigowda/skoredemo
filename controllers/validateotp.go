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

func (base *Controller) OTP_Validate(ctx *gin.Context) {
	var requestBody helpers.OTPRequestBody

	//check if body is present in payload
	if err := ctx.BindJSON(&requestBody); err != nil {
		logger.Errorf("error: %v", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid BODY",
		})
		return
	}

	//check if body has otp field
	if requestBody.Otp == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "otp is not provided",
		})
		return
	}

	//check otp length
	if len(requestBody.Otp) <= 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "otp is empty",
		})
		return
	}

	//convert otp value from string to int
	value := requestBody.Otp
	number, err := strconv.ParseUint(value, 10, 32)
	if err != nil {
		logger.Errorf("error: %v", err.Error())
	}
	finalIntNum := int(number)
	var key = helpers.DigitsCount(finalIntNum)
	if key == 6 {

		//get data from SL_OTP_Attempts for above otp
		var otp_results models.SL_OTP_Attempts
		otp_query_results := dbhelpers.FilterQuery(base.DB, "SL_OTP_Attempts", map[string]interface{}{"otp": requestBody.Otp, "status": "Valid"}, &otp_results)
		println(otp_query_results.RowsAffected)

		//RowsAffected is greater than 0 then otp already exist or expired
		if otp_query_results.RowsAffected > 0 {

			//check expiry time and current time difference
			t1 := time.Now()
			t2 := otp_results.Otp_expiry
			diff := t2.Sub(t1)

			//if time diff >0 than not expired
			if diff >= 0 {

				//check 2 otps
				if otp_results.Otp == finalIntNum {

					//update SL_USER and SL_OTP_Attempts tables
					dbhelpers.UpdateQuery(base.DB, "SL_USER", map[string]interface{}{"phone_number": otp_results.Phone_number}, map[string]interface{}{"otp_verified_status": true})
					dbhelpers.UpdateQuery(base.DB, "SL_OTP_Attempts", map[string]interface{}{"id": otp_results.Id}, map[string]interface{}{"otp_entered_time": time.Now(), "status": "invalid"})
					ctx.JSON(http.StatusOK, gin.H{
						"code":    http.StatusOK,
						"message": "OTP validated",
					})
				} else {
					ctx.JSON(http.StatusOK, gin.H{
						"code":    http.StatusBadRequest,
						"message": "Invalid otp",
					})
				}
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"code":    http.StatusBadRequest,
					"message": "OTP time is expried",
				})
			}
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusBadRequest,
				"message": "Invalid otp",
			})
		}

	} else {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "OTP length is greater or lesser",
		})
		return
	}
}
