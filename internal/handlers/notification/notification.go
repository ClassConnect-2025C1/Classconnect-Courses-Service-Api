package notification

import (
	"net/http"
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

// CurseEnrollNotification represents the data structure for an enrollment notification
type CurseEnrollNotification struct {
	ReceiverEmail string `json:"receiver_email"`
	Subject       string `json:"subject"`
	Text          string `json:"text"`
	HTML          string `json:"html"`
}

// HttpDoer defines an interface for HTTP client capabilities
type HttpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// NotificationClient handles sending notifications to users
type NotificationClient struct {
	Client                 HttpDoer
	NotificationServiceURL string
	UsersServiceURL        string
}

// NotificationSender defines the interface for sending notifications
type NotificationSender interface {
	// SendNotificationEmail sends an enrollment notification email to a user
	SendNotificationEmail(userId, courseName string)
}
