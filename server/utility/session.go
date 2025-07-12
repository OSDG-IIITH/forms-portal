package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

const SessionCookieName = "session"

var secret = []byte(Config.SessionSecret)

type Session struct {
	ID      string    `json:"id"`
	Expires time.Time `json:"expires"`
}

const DefaultSessionTtl = 7 * 24 * time.Hour

func CreateSession(id string, ttl time.Duration) Session {
	return Session{
		ID:      id,
		Expires: time.Now().Add(ttl),
	}
}

func (session *Session) String() string {
	raw, err := json.Marshal(session)
	if err != nil {
		log.Error(
			"failed to marshal session for cookie",
			"error", err.Error(),
		)
		os.Exit(1)
	}

	payload := base64.URLEncoding.EncodeToString(raw)
	hash := getHash([]byte(raw))

	return payload + "." + hash
}

func (session *Session) Cookie() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = SessionCookieName
	cookie.Expires = session.Expires
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Secure = Config.Domain[4] == 's' // only https
	cookie.Value = session.String()
	return cookie
}

func ValidateSession(raw string) (*Session, error) {
	splits := strings.Split(raw, ".")
	if len(splits) != 2 {
		return nil, errors.New("Invalid session.")
	}

	payload, err := base64.URLEncoding.DecodeString(splits[0])
	if err != nil {
		return nil, err
	}

	hash, err := base64.URLEncoding.DecodeString(splits[1])
	if err != nil {
		return nil, err
	}

	if !hmac.Equal(getRawHash([]byte(payload)), hash) {
		return nil, errors.New("Invalid hash.")
	}

	var session Session
	err = json.Unmarshal([]byte(payload), &session)
	if err != nil {
		return nil, err
	}

	if session.Expires.Before(time.Now()) {
		return nil, errors.New("Session has expired.")
	}
	return &session, nil
}

func getHash(payload []byte) string {
	return base64.URLEncoding.EncodeToString(getRawHash(payload))
}

func getRawHash(payload []byte) []byte {
	mac := hmac.New(sha256.New, []byte(Config.SessionSecret))
	return mac.Sum(payload)
}

func DeletionCookie() *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = SessionCookieName
	cookie.Value = ""
	cookie.HttpOnly = true
	cookie.MaxAge = -1
	cookie.HttpOnly = true
	cookie.Path = "/"
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Secure = Config.Domain[4] == 's'
	return cookie
}
