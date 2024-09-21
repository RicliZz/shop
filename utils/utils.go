package utils

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
	"net/http"
)

var Validator = validator.New()

func ParseJSON(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is empty")
	}

	return json.NewDecoder(r.Body).Decode(v)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func ErrorJSON(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func GetTokenFromRequest(r *http.Request) (string, error) {
	token := r.Header.Get("Authorization")
	if len(token) == 0 {
		return "", fmt.Errorf("token is empty")
	}
	return token, nil
}

func GenerateUUID() (uuid.UUID, error) {
	return uuid.New(), nil
}

func EmailSend(email string, token uuid.UUID) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "ricliz7@yandex.ru")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Confirm Email")
	confLink := fmt.Sprintf("http://localhost:8080/api/v1/confirm?token=%s", token)
	htmlBody := fmt.Sprintf(`<a href="%s">Confirm your account</a>`, confLink)
	m.SetBody("text/html", htmlBody)
	d := gomail.NewDialer("smtp.yandex.ru", 587, "ricliz7@yandex.ru", "qvytuvqibjjhbhhg")
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
