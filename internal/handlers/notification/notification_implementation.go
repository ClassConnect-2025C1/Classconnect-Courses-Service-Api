package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// NewNotificationClient creates a new notification client
func NewNotificationClient(client HttpDoer) *NotificationClient {
	notificationURL := os.Getenv("URL_NOTIFICATION")
	usersServiceURL := os.Getenv("URL_USERS")

	// if notificationURL == "" || usersServiceURL == "" {
	// 	log.Fatalf("Algunas variables de entorno no están configuradas correctamente.")
	// }

	if client == nil {
		client = http.DefaultClient
	}

	return &NotificationClient{
		Client:                 client,
		NotificationServiceURL: notificationURL,
		UsersServiceURL:        usersServiceURL,
	}
}

// SendNotificationEmail sends an enrollment notification email to a user
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

// getUserEmailFromService retrieves a user's email and name from the users service
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

func (sender *NotificationClient) SendNotification(userId, courseName, notification_type string) {
	userEmail, userName := sender.getUserEmailFromService(userId)
	if userEmail == "" {
		log.Printf("failed to get user email for userId: %s", userId)
		return
	}
	subject, TextTemplate, HtmlTemplate := "", "", ""
	switch notification_type {
	case "enrollment":
		subject = subjectTemplate
		TextTemplate = textTemplate
		HtmlTemplate = htmlTemplate
	case "feedback":
		subject = feedbackSubjectTemplate
		TextTemplate = textFeedbackTemplate
		HtmlTemplate = htmlFeedbackTemplate
	case "course_approve":
		subject = approvedSubjectTemplate
		TextTemplate = textAppovedTemplate
		HtmlTemplate = htmlAppovedTemplate
	case "new_assignment":
		subject = newAssigmentSubjectTemplate
		TextTemplate = textNewAssigmentTemplate
		HtmlTemplate = htmlNewAssigmentTemplate
	default:
		subject, TextTemplate, HtmlTemplate = "", "", ""
	}

	payload := NotificationPayload{
		ID:               userId,
		ReceiverEmail:    userEmail,
		NotificationType: notification_type,
		Subject:          subject,
		Text:             fmt.Sprintf(TextTemplate, userName, courseName),
		HTML:             fmt.Sprintf(HtmlTemplate, userName, courseName),
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("failed to marshal payload: %v", err)
		return
	}

	url := sender.NotificationServiceURL + "/notifications/send"
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

func (sender *NotificationClient) SendNotificationToAll(allUsers []map[string]any, courseName, notification_type string) {
	for _, m := range allUsers {
		if userID, ok := m["user_id"]; ok {
			sender.SendNotification(userID.(string), courseName, notification_type)
		}
	}
}
