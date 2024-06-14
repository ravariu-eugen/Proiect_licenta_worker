package main

import (
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	awsAccessKeyID     = ""
	awsSecretAccessKey = ""
)

func setCredentials(awsAccessKey, awsSecretKey string) {

	awsAccessKeyID = awsAccessKey
	awsSecretAccessKey = awsSecretKey

}

func setupData(inputFile string, sharedFile string) error {

	// load data from s3 bucket
	bucket := "rav123bucket"

	// set session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(
			awsAccessKeyID,
			awsSecretAccessKey,
			""),
	})
	// download
	downloader := s3manager.NewDownloader(sess)
	file, err := os.Create(item)
	numBytes, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})

	// download data

	// untar the data

	cmd := exec.Command("tar", "-xvf", inputFile, "-C", InputFolder)
	err := cmd.Run()
	if err != nil {
		return err
	}

	cmd = exec.Command("tar", "-xvf", sharedFile, "-C", SharedFolder)
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
