package awsshared

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type AWSSession struct {
	REGION  string
	SESSION *session.Session
}

func GetSession() AWSSession {
	AWS_REGION := os.Getenv("AWS_REGION")

	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.EnvProvider{},
		})
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWS_REGION),
		Credentials: creds,
	})
	if err != nil {
		log.Fatal(err)
	}

	return AWSSession{AWS_REGION, s}
}
