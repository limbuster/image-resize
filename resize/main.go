package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// GetOriginalS3Key will return s3key for the original prefix
func GetOriginalS3Key(s3key string) string {
	var re = regexp.MustCompile(`^upload\/(.*\.jpg)($)`)
	return re.ReplaceAllString(s3key, `original/$1`)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, s3Event events.S3Event) {
	for _, record := range s3Event.Records {
		s3 := record.S3
		fmt.Printf("[%s - %s] Bucket = %s, Key = %s, Region =%s \n", record.EventSource, record.EventTime, s3.Bucket.Name, s3.Object.Key, record.AWSRegion)
		s3Helper := S3Helper{region: record.AWSRegion, bucket: s3.Bucket.Name}
		s3key := s3.Object.Key
		s3bytes := s3Helper.GetS3Object(s3key)
		hash := md5.Sum(s3bytes)
		md5checksum := hex.EncodeToString(hash[:])
		fmt.Printf("MD checksum: %s \n", md5checksum)
		s3Helper.CopyS3Object(s3key, "original/"+md5checksum+".jpg")
		thumbnail := ResizeImage(s3bytes, 500, 500)
		s3Helper.PutS3Object(thumbnail, "thumbnail/"+md5checksum+".jpg")
		s3Helper.DeleteS3Object(s3key)
	}
	return
}

func main() {
	lambda.Start(Handler)
}
