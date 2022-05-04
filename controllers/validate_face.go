package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sk-integrated-services/models"
	"sk-integrated-services/pkg/helpers"
	"sk-integrated-services/pkg/helpers/dbhelpers"
	"sk-integrated-services/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (base *Controller) Validate_face(ctx *gin.Context) {

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

	//get multiple files from request
	form, er := ctx.MultipartForm()
	if form == nil {
		logger.Errorf("error: %v", er.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "No payload in body",
		})
		return
	}
	//store file details in an array
	files := form.File["file"]
	s := []multipart.File{}
	filenames := []string{}
	for _, file := range files {
		s1, _ := file.Open()
		s = append(s, s1)
		filenames = append(filenames, file.Filename)
	}

	if len(filenames) != 2 {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "2 images are required",
		})
		return
	}

	//create locally for 1st file
	helpers.Createfile(filenames[0], s[0])
	helpers.Createfile(filenames[1], s[1])

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("firstImage", filenames[0])
	if err != nil {
		logger.Errorf("error: %v", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Somthing failed while creating form field",
		})
		return
	}

	f, err := os.Open("public/" + filenames[0])

	if err != nil {
		logger.Errorf("error: %v", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Failed to open File",
		})
		return
	}
	_, err = io.Copy(fw, f)

	if err != nil {
		logger.Errorf("error: %v", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "File Copy Failed",
		})
		return
	}

	selfie_fw, err := writer.CreateFormFile("secondImage", filenames[1])
	if err != nil {
		logger.Errorf("error: %v", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Somthing failed while creating form field",
		})
		return
	}

	selfie_f, err := os.Open("public/" + filenames[1])

	if err != nil {
		logger.Errorf("error: %v", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "Failed to open File",
		})
		return
	}

	defer selfie_f.Close()
	defer f.Close()

	_, err = io.Copy(selfie_fw, selfie_f)

	if err != nil {
		logger.Errorf("error: %v", err.Error())
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusBadRequest,
			"message": "File Copy Failed",
		})
		return
	}
	writer.Close()

	firstImage := helpers.Coverttobase64("public/" + filenames[0])
	secondImage := helpers.Coverttobase64("public/" + filenames[1])

	url := helpers.GetEnvVariable("VIDA_FACE_API_URL")

	postBody, _ := json.Marshal(map[string]string{
		"firstImage":  firstImage,
		"secondImage": secondImage,
	})
	responseBody := bytes.NewBuffer(postBody)

	req, err := http.NewRequest("POST", url, responseBody)
	vida_token := helpers.Generatevidatoken()
	vida_token = "Bearer " + vida_token
	if err != nil {
		logger.Errorf("error: %v", err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", vida_token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("error: %v", err.Error())
	}
	// defer resp.Body.Close()

	resp_body, _ := ioutil.ReadAll(resp.Body)

	var data map[string]interface{}
	json.Unmarshal([]byte(resp_body), &data)

	helpers.Deletefile("public/" + filenames[0])
	helpers.Deletefile("public/" + filenames[1])

	if data["data"] == nil {
		var errors helpers.Error_response
		json.Unmarshal([]byte(string(resp_body)), &errors)
		ctx.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": errors.Errors.Detail,
		})
		return
	} else {
		var success helpers.Success_response
		json.Unmarshal([]byte(string(resp_body)), &success)
		score := success.Data.Score

		conn := bytes.NewReader([]byte(score))
		buffer := make([]byte, 256)

		i, _ := conn.Read(buffer)

		result := string(buffer[:i])
		floatResult, _ := strconv.ParseFloat(result, 64)
		var results models.SL_USER
		user_response := dbhelpers.FilterQuery(base.DB, "SL_USER", map[string]interface{}{"phone_number": token_validation_repsonse.Phone}, &results)
		println(user_response.RowsAffected)

		VIDA_FACE_API_SCORE := helpers.GetEnvVariable("VIDA_FACE_API_SCORE")
		parsed_otp_count, err := strconv.ParseFloat(VIDA_FACE_API_SCORE, 64)
		if err != nil {
			logger.Errorf("error: %v", err.Error())
		}
		if floatResult > float64(parsed_otp_count) {
			face_result := dbhelpers.InsertQuery(base.DB, "SL_KYC_facematch", map[string]interface{}{"user_id": results.User_id, "status": "Matched", "match": true, "score": floatResult})
			println(face_result.RowsAffected)
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusOK,
				"message": "Face matched successfully",
			})
			return
		} else {
			face_result := dbhelpers.InsertQuery(base.DB, "SL_KYC_facematch", map[string]interface{}{"user_id": results.User_id, "status": "Failure", "match": false, "score": floatResult})
			println(face_result.RowsAffected)
			ctx.JSON(http.StatusOK, gin.H{
				"code":    http.StatusBadRequest,
				"message": success.Data.Message,
			})
			return
		}
	}
}
