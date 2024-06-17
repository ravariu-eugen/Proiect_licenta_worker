package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// downloadFile downloads a file from a specified S3 bucket and saves it to a given destination folder.
//
// Parameters:
//   - bucketName: the name of the S3 bucket from which to download the file (string)
//   - filename: the name of the file to download (string)
//   - destinationFolder: the folder where the downloaded file will be saved (string)
//
// Returns:
//   - error: an error if the file download fails, otherwise nil (error)
func downloadFile(bucketName, filename, destinationFolder string) error {

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-west-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(awsCred.AccessKeyID, awsCred.SecretAccessKey, "")))

	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	result, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
	}, func(options *s3.Options) {
		options.UsePathStyle = true
	})
	if err != nil {
		return err
	}

	body, _ := io.ReadAll(result.Body)

	// Create the destination folder if it doesn't exist
	err = os.MkdirAll(destinationFolder, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(destinationFolder+"/"+filename, body, 0644)
	if err != nil {
		return err
	}

	return nil
}

func unarchiveFile(filepath, destinationFolder string) error {

	filename := filepath[strings.LastIndex(filepath, "/")+1:]
	fileType := filename[:strings.Index(filepath, ".")]

	err := os.MkdirAll(destinationFolder, os.ModePerm)
	if err != nil {
		return err
	}
	var cmd *exec.Cmd
	switch fileType {
	case "tar":
		cmd = exec.Command("tar", "-xf", filepath, "-C", destinationFolder)
	case "zip":
		cmd = exec.Command("unzip", filepath, "-d", destinationFolder)
	default:
		return fmt.Errorf("unsupported file type: %s", fileType)
	}

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
