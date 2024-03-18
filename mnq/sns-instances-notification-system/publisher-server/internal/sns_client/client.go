package sns_client

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"os"
)

var AwsSns *sns.SNS

func InitAwsSnsClient() {
	awsSession := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("fr-par"),
		Endpoint:    aws.String("http://sns.mnq.fr-par.scaleway.com"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), ""),
	}))

	AwsSns = sns.New(awsSession)
}
