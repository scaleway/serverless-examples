package gomail

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	tem "github.com/scaleway/scaleway-sdk-go/api/tem/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Body to send to the function (via curl for example or Messaging & Queuing).

//	{
//		"to": "mail@mail.mail",
//		"subject": "from console test",
//		"message": "very very long msg"
//	}

// Region used to call the API.
const region = scw.RegionFrPar

// Data holds the body of the HTTP call. Fields in Data must be completed
// to send an email
type Data struct {
	Subject string `json:"subject" 	validate:"required"`
	Message string `json:"message" 	validate:"required"`
	To      string `json:"to" 		validate:"required,email"`
}

// Handler is the entrypoint of the function.
func Handler(respWriter http.ResponseWriter, req *http.Request) {
	// Only allow POST verbs
	if req.Method != http.MethodPost {
		respWriter.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	var body Data

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		respWriter.WriteHeader(http.StatusBadRequest)

		return
	}

	// Validate the email data, return a 400 error if not valid

	if err := validator.New().Struct(body); err != nil {
		respWriter.WriteHeader(http.StatusBadRequest)
		_, _ = respWriter.Write([]byte(err.Error()))

		return
	}

	if err := sendMail(body.Subject, body.Message, body.To, "CHANGE_ME", false); err != nil {
		panic(err)
	}
}

// sendMail sends a mail to "mailTo" using "from" email.
// If checkMailStatus is at true, the function will take more time to run some calls
// to the API in order to get mail status over time.
func sendMail(subject, content, mailTo, from string, checkMailStatus bool) error {
	// Create a Scaleway client
	client, err := scw.NewClient(
		// Get your organization ID at https://console.scaleway.com/organization/settings
		scw.WithDefaultOrganizationID(os.Getenv("SCW_DEFAULT_ORGANIZATION_ID")),
		// Get your credentials at https://console.scaleway.com/iam/api-keys
		scw.WithAuth(os.Getenv("SCW_ACCESS_KEY"), os.Getenv("SCW_SECRET_KEY")),
		// Get more about our availability zones at
		// https://www.scaleway.com/en/docs/console/my-account/reference-content/products-availability/
		scw.WithDefaultRegion(region),
	)
	if err != nil {
		return fmt.Errorf("error creating scaleway client with sdk %w", err)
	}

	// Create SDK object to manipulate transactional email.
	temAPI := tem.NewAPI(client)

	// Create email is used to send the email to the destination.
	mailResp, err := temAPI.CreateEmail(&tem.CreateEmailRequest{
		Subject:   subject,
		Text:      content,
		ProjectID: os.Getenv("SCW_DEFAULT_ORGANIZATION_ID"),
		From:      &tem.CreateEmailRequestAddress{Email: from},
		To:        []*tem.CreateEmailRequestAddress{{Email: mailTo}},
		Region:    region,
	})
	if err != nil {
		return fmt.Errorf("error trying to create and send mail %w", err)
	}

	if checkMailStatus {
		// Now we get the email ID in order to get it's status
		emailID := mailResp.Emails[0].ID

		mail, err := temAPI.GetEmail(&tem.GetEmailRequest{EmailID: emailID})
		if err != nil {
			return fmt.Errorf("error trying to get mail on first time %w", err)
		}

		fmt.Println("mail status :", mail.Status.String())

		// This sleep is only for educational purposes, it should be removed for industrial applications.
		time.Sleep(5 * time.Second)

		mail, err = temAPI.GetEmail(&tem.GetEmailRequest{EmailID: emailID})
		if err != nil {
			return fmt.Errorf("error trying to get mail on second time %w", err)
		}

		fmt.Println("mail status after 5 seconds :", mail.Status.String())
		// expected output:
		// mail status : new
		// mail status after 5 seconds : sent
	}

	return nil
}
