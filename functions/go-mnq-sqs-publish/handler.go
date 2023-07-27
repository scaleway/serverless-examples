package myfunc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-playground/validator/v10"
)

var (
	sqsClient *sqs.SQS
	queueURL  *string
)

func init() {
	accessKey := os.Getenv("SQS_ACCESS_KEY")
	secretKey := os.Getenv("SQS_SECRET_KEY")

	sqsSession := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("fr-par"),
		Endpoint:    aws.String("http://sqs-sns.mnq.fr-par.scw.cloud"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	}))

	sqsClient = sqs.New(sqsSession)
	queueUrl, err := getOrCreateQueue("chat-message")
	if err != nil {
		panic(err)
	}
	queueURL = queueUrl
}

type Body struct {
	Username string `json:"username" validate:"required"`
	Message  string `json:"message" validate:"required"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var body Body
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validator.New().Struct(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = sqsClient.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageAttributes: map[string]*sqs.MessageAttributeValue{
			"Username": {
				DataType:    aws.String("String"),
				StringValue: aws.String(body.Username),
			},
		},
		MessageBody: aws.String(body.Message),
		QueueUrl:    queueURL,
	})
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getOrCreateQueue(name string) (*string, error) {
	result, err := sqsClient.GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(name),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == sqs.ErrCodeQueueDoesNotExist {
			result, err := sqsClient.CreateQueue(&sqs.CreateQueueInput{
				QueueName: aws.String(name),
			})
			if err != nil {
				return nil, err
			}

			return result.QueueUrl, nil
		}

		return nil, err
	}

	return result.QueueUrl, nil
}
