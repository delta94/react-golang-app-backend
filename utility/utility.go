package utility

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"labix.org/v2/mgo/bson"
)

//ResponseStruct struct message with token
type ResponseStruct struct {
	Message		string
	Token		string
	Data		interface{}
}

//Response struct
type Response struct {
	Message		string
}

//Error struct
type Error struct {
	ErrorMsg	string
	Error		error
}

//ResponseJSON struct
type ResponseJSON struct {
	Message		string
	Data		interface{}
}

//ResponseWithJSON to the http body
func ResponseWithJSON(msg string, data interface{}) (resp []byte) {
	msgResp := ResponseJSON{msg, data}
	response, err := json.Marshal(msgResp)
	if err != nil {
		log.Print("Error in marshall", err)
		return
	}
	return response
}

//SucessResponse to http body
func SucessResponse(msg string) (resp []byte) {
	msgResp := Response{msg}
	response, err := json.Marshal(msgResp)
	if err != nil {
		log.Print("Error in marshall", err)
		return
	}
	return response
}

//ErrorResponse to http body
func ErrorResponse(msg string, err error) (resp []byte) {
	errorMsg := Error{msg, err}
	response, err := json.Marshal(errorMsg)
	if err != nil {
		log.Print("Error in marshall", err)
		return
	}
	return response
}

//ResponseWithToken response body with token
func ResponseWithToken(msg string, token string, user interface{}) (resp []byte) {
	msgToken := ResponseStruct{msg, token, user}
	response, err := json.Marshal(msgToken)
	if err != nil {
		log.Print("Error in marshall", err)
		return
	}
	return response
}

//VerifyToken verification of JWT token
func VerifyToken(r *http.Request) (status bool) {
	auth := strings.Replace(r.Header.Get("Authorization"), "Bearer ", "", 1)
	if auth == "undefined" {
		return false
	}
	godotenv.Load(".env")
	jwtKey := []byte(os.Getenv("JWT_KEY"))

	token, _ := jwt.Parse(auth, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("There was an error")
        }
        return jwtKey, nil
	})

	return token.Valid
}

//UploadImageToS3 - get session and file to upload to AWS S3
func UploadImageToS3(ses *session.Session, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	godotenv.Load(".env")
	bucket := os.Getenv("AWS_BUCKET")
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	fileName := "avatar/" + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)
	_, err := s3.New(ses).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(fileName),
		ACL: aws.String("public-read"),
		Body: bytes.NewReader(buffer),
		ContentLength: aws.Int64(int64(size)),
		ContentType: aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass: aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return fileName, err
}

//ConvertToTime get slice of uint8 to return time.Time
func ConvertToTime(timeInt []uint8) (time.Time, error ) {
	string := string(timeInt)
	build := string[0:10] + "T" + string[11:] + "Z"
	conv, err := time.Parse(time.RFC3339, build)
	if err != nil {
		return time.Time{}, err
	}
	return conv, nil
}