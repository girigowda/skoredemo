package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sk-integrated-services/models"
	"sk-integrated-services/pkg/helpers"
	"sk-integrated-services/pkg/helpers/dbhelpers"
	"sk-integrated-services/pkg/logger"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func (base *Controller) Validate_card(ctx *gin.Context) {
	if ctx.Request.Header.Get("Authorization") == "" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "No token",
		})
		return
	}
	token_status := helpers.ValidateJWT(ctx.Request.Header.Get("Authorization"))
	var token_validation_repsonse helpers.TOKEN_RESP
	json.Unmarshal([]byte(token_status), &token_validation_repsonse)

	if token_validation_repsonse.Status != "success" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Invalid Token",
		})
		return
	}

	file, header, err := ctx.Request.FormFile("file")
	if file == nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "No file",
		})
		return
	}
	if err != nil {
		logger.Errorf("error: %v", err.Error())
		return
	}

	//extract filename and create a file
	filename := header.Filename
	helpers.Createfile(filename, file)

	url := helpers.GetEnvVariable("AKSATA_OCR_URL")
	fmt.Println("URL:>", url)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	//open that file and write file to body
	fw, err := writer.CreateFormFile("ocrImage", filename)
	if err != nil {
		logger.Errorf("error: %v", err.Error())
	}

	f, err := os.Open("public/" + filename)

	if err != nil {
		logger.Errorf("error: %v", err.Error())
	}
	_, err = io.Copy(fw, f)

	if err != nil {
		logger.Errorf("error: %v", err.Error())
	}
	writer.Close()
	req, err := http.NewRequest("POST", url, bytes.NewReader(body.Bytes()))

	if err != nil {
		logger.Errorf("error: %v", err.Error())
	}

	//call aksata api
	aksata_key := helpers.GetEnvVariable("AKSATA_KEY")

	req.Header.Set("X-AKSATA-KEY", aksata_key)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		logger.Errorf("error: %v", err.Error())

	}

	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal([]byte(respBody), &data)
	code := data["code"]

	//remove image from public
	helpers.Deletefile("public/" + filename)

	//check error or success message
	if code == "SUCCESS" {
		var success helpers.Card_success_response
		json.Unmarshal([]byte(string(respBody)), &success)
		var results models.SL_USER
		user_response := dbhelpers.FilterQuery(base.DB, "SL_USER", map[string]interface{}{"phone_number": token_validation_repsonse.Phone}, &results)
		println(user_response.RowsAffected)

		split := strings.Split(success.Data.BirthPlaceBirthday, " ")
		birth_place := strings.Replace(split[0], ",", "", -1)
		birth_date := split[1]

		kyc_result := dbhelpers.InsertQuery(base.DB, "SL_KYC_Details", map[string]interface{}{"gender": success.Data.Gender, "full_name": success.Data.Name, "blood_group": success.Data.BloodType, "relegion": success.Data.Religion, "province": success.Data.Province, "city": success.Data.City, "district": success.Data.District, "village": success.Data.Village, "rt_rw": success.Data.Rtrw, "occupation": success.Data.Occupation, "expiry_date": success.Data.ExpiryDate, "marital_status": success.Data.MaritalStatus, "nationality": success.Data.Nationality, "nik_number": success.Data.IdNumber, "address": success.Data.Address, "created_at": time.Now(), "updated_at": time.Now(), "user_id": results.User_id, "card_upload_staus": true, "place_of_birth": birth_place, "dob": birth_date})
		println(kyc_result.RowsAffected)

		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "KYC details captured",
		})
		return
	} else if code == "ERROR" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Error in parameter",
		})
		return
	} else if code == "OCR_NO_RESULT" {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "OCR check failed, unable to find any ktp field in the uploaded picture",
		})
		return
	}
}
