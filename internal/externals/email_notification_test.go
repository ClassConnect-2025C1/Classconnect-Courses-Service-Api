package externals

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockDoer struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockDoer) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func captureLogs(f func()) string {
	var buf bytes.Buffer
	logger := log.Default()
	origOutput := logger.Writer()
	logger.SetOutput(&buf)
	defer logger.SetOutput(origOutput)

	f()
	return buf.String()
}

func TestSendNotificationEmail_Success(t *testing.T) {
	os.Setenv("URL_NOTIFICATION", "http://notification")
	os.Setenv("URL_USERS", "http://users")

	mock := &mockDoer{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.String(), "/users/profile/") {
				body := `{"email":"test@example.com"}`
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(body)),
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
			}, nil
		},
	}

	client := NewNotificationClient(mock)
	logs := captureLogs(func() {
		client.SendNotificationEmail("123", "Juan, curso de Go")
	})

	assert.Empty(t, logs, "No se esperan logs en caso éxito")
}

func TestSendNotificationEmail_GetUserFails(t *testing.T) {
	os.Setenv("URL_NOTIFICATION", "http://notification")
	os.Setenv("URL_USERS", "http://users")

	mock := &mockDoer{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
			}, nil
		},
	}

	client := NewNotificationClient(mock)
	logs := captureLogs(func() {
		client.SendNotificationEmail("123", "Juan, curso de Go")
	})

	assert.Contains(t, logs, "unexpected status code from user", "Debe loguear error al obtener email")
}

func TestSendNotificationEmail_PostFails(t *testing.T) {
	os.Setenv("URL_NOTIFICATION", "http://notification")
	os.Setenv("URL_USERS", "http://users")

	mock := &mockDoer{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.String(), "/users/profile/") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"email":"test@example.com"}`)),
				}, nil
			}
			return nil, io.EOF
		},
	}

	client := NewNotificationClient(mock)
	logs := captureLogs(func() {
		client.SendNotificationEmail("123", "Juan, curso de Go")
	})

	assert.Contains(t, logs, "failed to send notification request", "Debe loguear error al enviar POST")
}

func TestSendNotificationEmail_PostUnexpectedStatus(t *testing.T) {
	os.Setenv("URL_NOTIFICATION", "http://notification")
	os.Setenv("URL_USERS", "http://users")

	mock := &mockDoer{
		DoFunc: func(req *http.Request) (*http.Response, error) {
			if strings.Contains(req.URL.String(), "/users/profile/") {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(strings.NewReader(`{"email":"test@example.com"}`)),
				}, nil
			}
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("error")),
			}, nil
		},
	}

	client := NewNotificationClient(mock)
	logs := captureLogs(func() {
		client.SendNotificationEmail("123", "Juan, curso de Go")
	})

	assert.Contains(t, logs, "unexpected status code from notification", "Debe loguear status inesperado de POST")
}
