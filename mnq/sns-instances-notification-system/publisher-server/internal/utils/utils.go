package utils

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
	"os"
	"publisher-server/internal/sns_client"
)

func PublishMessage(message string) error {
	_, err := sns_client.AwsSns.Publish(&sns.PublishInput{
		Message:  aws.String(message),
		TopicArn: aws.String(os.Getenv("TOPIC_ARN")),
	})
	if err != nil {
		return fmt.Errorf("publish message to topic failed: %w", err)
	}

	log.Println("Message published successfully")
	return nil
}
