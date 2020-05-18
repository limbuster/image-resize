package main

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3Helper will be used for storing s3 info
type S3Helper struct {
	region string // aws region
	bucket string // s3 bucket
}

// DeleteS3Object deletes S3 object depending on the s3Helper struct
func (s3Helper S3Helper) DeleteS3Object(key string) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(s3Helper.region)})
	svc := s3.New(sess)
	_, err = svc.DeleteObject(&s3.DeleteObjectInput{Bucket: aws.String(s3Helper.bucket), Key: aws.String(key)})
	if err != nil {
		log.Printf("Unable to delete object %q from bucket %q, %v\n", key, s3Helper.bucket, err)
	}
	log.Printf("Successfully deleted %q from %q\n", key, s3Helper.bucket)
}

// CopyS3Object copies S3 object depending on the s3Helper struct and targetBucket
func (s3Helper S3Helper) CopyS3Object(sourceKey string, targetKey string) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(s3Helper.region)})
	svc := s3.New(sess)
	_, err = svc.CopyObject(&s3.CopyObjectInput{Bucket: aws.String(s3Helper.bucket), CopySource: aws.String(s3Helper.bucket + "/" + sourceKey), Key: aws.String(targetKey)})
	if err != nil {
		log.Printf("Unable to copy object %q from %q in %q bucket, %v\n", targetKey, targetKey, s3Helper.bucket, err)
	}
	log.Printf("Successfully copied %q to %q\n", sourceKey, targetKey)
}

// GetS3Object gets s3 object
func (s3Helper S3Helper) GetS3Object(key string) []byte {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(s3Helper.region)})
	svc := s3.New(sess)
	s3obj, err := svc.GetObject(&s3.GetObjectInput{Bucket: aws.String(s3Helper.bucket), Key: aws.String(key)})
	body, err := ioutil.ReadAll(s3obj.Body)
	if err != nil {
		log.Printf("Unable to get object %q from %q bucket, %v\n", key, s3Helper.bucket, err)
	}
	log.Printf("Successfully got %q\n", key)
	return body
}

// PutS3Object gets s3 object
func (s3Helper S3Helper) PutS3Object(imageBytes []byte, key string) {
	sess, err := session.NewSession(&aws.Config{Region: aws.String(s3Helper.region)})
	svc := s3.New(sess)
	body := bytes.NewReader(imageBytes)
	_, err = svc.PutObject(&s3.PutObjectInput{Bucket: aws.String(s3Helper.bucket), Key: aws.String(key), ContentType: aws.String("image/jpeg"), Body: body})

	if err != nil {
		log.Printf("Unable to put object %q to %q bucket, %v\n", key, s3Helper.bucket, err)
	}
	log.Printf("Successfully put %q\n", key)
}
