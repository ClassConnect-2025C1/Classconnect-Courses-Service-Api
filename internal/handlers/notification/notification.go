package notification

import (
	"net/http"
)

var (
	subjectTemplate = "ClassConnect - Inscripción exitosa"
	textTemplate    = `Felicitaciones %s!.
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

	approvedSubjectTemplate = "ClassConnect - Curso aprobado"
	textAppovedTemplate     = `Felicitaciones %s!.
Aprobaste el curso %s.`

	htmlAppovedTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Inscripción exitosa</title>
</head>
<body>
  <p>Felicitaciones %s!<br>
  Aprobaste el curso %s.</p>
</body>
</html>`

	feedbackSubjectTemplate = "ClassConnect - Feedback disponible"
	textFeedbackTemplate    = `Importante %s!.
Tenes feedback del curso %s.`

	htmlFeedbackTemplate = `<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Inscripción exitosa</title>
</head>
<body>
  <p>Importante %s!<br>
  Tenes feedback del curso %s.</p>
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

type NotificationPayload struct {
	ID               string `json:"id"`
	ReceiverEmail    string `json:"receiver_email"`
	NotificationType string `json:"notification_type"`
	Subject          string `json:"subject"`
	Text             string `json:"text"`
	HTML             string `json:"html"`
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
