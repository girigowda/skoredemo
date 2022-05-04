package helpers

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"sk-integrated-services/models"
	"sk-integrated-services/pkg/helpers/dbhelpers"
	"sk-integrated-services/pkg/logger"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9'}

func Generatenumber(length int) []byte {
	max := length
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	println("otp is", string(b))
	return b
}

func Coverttobase64(img string) string {
	imgFile, err := os.Open(img) // a QR code image

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer imgFile.Close()

	// create a new buffer base on file size
	fInfo, _ := imgFile.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	// read file content into buffer
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)

	imgBase64Str := base64.StdEncoding.EncodeToString(buf)
	return imgBase64Str
}

func DigitsCount(number int) int {
	var count int = 0
	for number != 0 {
		number /= 10
		count += 1
	}
	return count

}

func GetEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		println("Error loading .env file")
	}

	return os.Getenv(key)
}

func Deletefile(filepath string) string {
	err1 := os.Remove(filepath)
	if err1 != nil {
		fmt.Println(err1)
	}
	return "deleted"
}

func Createfile(filename string, file io.Reader) string {
	out, _ := os.Create("public/" + filename)
	defer out.Close()

	_, err := io.Copy(out, file)
	if err != nil {
		logger.Errorf("error: %v", err.Error())
	}
	defer out.Close()
	return "created"
}

func Generateid(conn *gorm.DB) string {
	var id_count models.SL_numbers
	k_id := dbhelpers.FilterQuery(conn, "SL_numbers", map[string]interface{}{"id": 1}, &id_count)
	println(k_id.RowsAffected)
	user_id := "SL_" + strconv.Itoa(time.Now().Day()) + strconv.Itoa(time.Now().Year()) + strconv.Itoa(id_count.Count)
	conn.Table("SL_numbers").Where("id = ?", 1).Updates(map[string]interface{}{"count": id_count.Count + 1})
	return user_id
}
