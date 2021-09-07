package s3_storage

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/DuckLuckBreakout/currency-exchange-backend/pkg/configer"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type ApplicationType int

var (
	sess  *session.Session
	types = map[string]string{
		"image/png":  "png",
		"image/jpg":  "jpg",
		"image/jpeg": "jpeg",
	}
)

const (
	Avatar ApplicationType = iota
	Background
)

func InitNewConnection(s3Config *configer.S3Data) {
	var err error
	sess, err = session.NewSession(&aws.Config{
		Region:   aws.String(s3Config.Region),
		Endpoint: aws.String(s3Config.Endpoint),
		Credentials: credentials.NewStaticCredentials(
			s3Config.AccessKeyId,
			s3Config.SecretAccessKey,
			"",
		),
	})

	if err != nil {
		log.Fatal(err)
	}
}

func UploadMultipartFile(file *multipart.File, fileType string,
	appType ApplicationType, userId uint64, s3Config *configer.S3Data) (string, error) {
	// Select path to upload file
	var path string
	switch appType {
	case Avatar:
		path = "avatar/"
	case Background:
		path = "background/"
	default:
		return "", errors.New("incorrect app type")
	}

	// Select file type
	extension, ok := types[fileType]
	if !ok {
		return "", errors.New("incorrect file type")
	}

	fileName := md5.Sum([]byte(fmt.Sprintf("%d-%d-%d", userId, appType, time.Now().UnixNano())))
	fileKey := fmt.Sprintf("%s%x.%s", path, fileName, extension)
	_, err := s3manager.NewUploader(sess).Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3Config.Bucket),
		ACL:    aws.String(s3Config.Acl),
		Key:    aws.String(fileKey),
		Body:   *file,
	})
	if err != nil {
		return "", err
	}

	return fileKey, nil
}

func DeleteFileByKey(fileName string, s3Config *configer.S3Data) error {
	_, err := s3.New(sess).DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s3Config.Bucket),
		Key:    aws.String(fileName),
	})

	return err
}

func PathToFile(fileName string, appType ApplicationType, s3Config *configer.S3Data) string {
	if fileName == "" {
		switch appType {
		case Avatar:
			fileName = "avatar/default.jpeg"
		case Background:
			fileName = "background/default.jpg"
		}
	}
	return fmt.Sprintf(
		"https://%s.%s/%s",
		s3Config.Bucket,
		s3Config.Endpoint,
		fileName,
	)
}
