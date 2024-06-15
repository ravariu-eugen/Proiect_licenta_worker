package main

var (
	awsAccessKeyID     = ""
	awsSecretAccessKey = ""
)

func setCredentials(awsAccessKey, awsSecretKey string) {

	awsAccessKeyID = awsAccessKey
	awsSecretAccessKey = awsSecretKey

}
