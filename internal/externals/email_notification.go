package externals

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	textTemplate = `Felicitaciones %s!.
Tu inscripción al curso %s fue exitosa.`

	htmlTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Inscripción exitosa</title>
</head>
<body>
  <p>Felicitaciones %s!<br>
  Tu inscripción al curso %s fue exitosa.</p>
</body>
</html>`
)

type CurseEnrollNotification struct {
	ReceiverEmail string `json:"receiver_email"`
	Subject       string `json:"subject"`
	Text          string `json:"text"`
	HTML          string `json:"html"`
}

type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type NotificationClient struct {
	Client                 HttpDoer
	NotificationServiceURL string
	UsersServiceURL        string
}

func NewNotificationClient(client HttpDoer) *NotificationClient {
	notificationURL := os.Getenv("URL_NOTIFICATION")
	usersServiceURL := os.Getenv("URL_USERS")
	if notificationURL == "" || usersServiceURL == "" {
		log.Fatalf("Algunas variables de entorno no están configuradas correctamente.")
	}

	if client == nil {
		client = http.DefaultClient
	}

	return &NotificationClient{
		Client:                 client,
		NotificationServiceURL: notificationURL,
		UsersServiceURL:        usersServiceURL,
	}
}

func (sender *NotificationClient) SendNotificationEmail(userId, courseName string) {
	userEmail, userName := sender.getUserEmailFromService(userId)
	if userEmail == "" {
		log.Printf("failed to get user email for userId: %s", userId)
		return
	}

	payload := CurseEnrollNotification{
		ReceiverEmail: userEmail,
		Subject:       "ClassConnect - Inscripción exitosa",
		Text:          fmt.Sprintf(textTemplate, userName, courseName),
		HTML:          fmt.Sprintf(htmlTemplate, userName, courseName),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal payload: %v", err)
		return
	}

	url := sender.NotificationServiceURL + "/notifications/email"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := sender.Client.Do(req)
	if err != nil {
		log.Printf("failed to send notification request: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("unexpected status code from notification: %d", resp.StatusCode)
	}
}

func (sender *NotificationClient) getUserEmailFromService(userId string) (string, string) {
	url := sender.UsersServiceURL + "/users/profile/" + userId

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("failed to create user request: %v", err)
		return "", ""
	}

	resp, err := sender.Client.Do(req)
	if err != nil {
		log.Printf("failed to send user request: %v", err)
		return "", ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("unexpected status code from user: %d", resp.StatusCode)
		return "", ""
	}

	var response struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Printf("failed to decode response from user: %v", err)
		return "", ""
	}

	return response.Email, response.Name
}
