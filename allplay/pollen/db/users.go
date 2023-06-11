package db

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
	"github.com/hibooboo2/ggames/allplay/pollen/constants"
)

var (
	users    = map[string][32]byte{}
	sessions = map[string]*UserSession{}
	tempIDs  = map[string]struct{}{}
)

type UserSession struct {
	ID       string
	Username string
	Expiry   time.Time
}

func init() {
	RegisterUser("root@jhrb.us", "root", "root")
}

func (s *UserSession) Cookie() *http.Cookie {
	return &http.Cookie{
		Name:    constants.SessionCookieName,
		Value:   s.ID,
		Expires: s.Expiry,
	}
}

func RegisterUser(email, username, password string) error {
	if _, ok := users[username]; ok {
		return errors.New("user already exists")
	}

	_, isTemp := tempIDs[username]
	if isTemp {
		delete(tempIDs, username)
	}

	users[username] = sha256.Sum256([]byte(password))
	return nil
}

func CheckPassword(username string, password string) bool {
	passHash, ok := users[username]
	if !ok {
		return false
	}

	passwordHash := sha256.Sum256([]byte(password))

	return subtle.ConstantTimeCompare(passwordHash[:], passHash[:]) == 1
}

func Login(username string, password string, timeout time.Duration) (*http.Cookie, error) {
	if !CheckPassword(username, password) {
		return nil, errors.New("invalid username or password")
	}

	sessionID := uuid.Must(uuid.NewV4()).String()
	expires := time.Now().Add(timeout)
	sessions[sessionID] = &UserSession{
		ID:       sessionID,
		Username: username,
		Expiry:   expires,
	}

	return sessions[sessionID].Cookie(), nil
}

func GetTempID() string {
	tempID := uuid.Must(uuid.NewV4()).String()
	tempIDs[tempID] = struct{}{}
	return tempID
}

func IsTempID(tempID string) bool {
	_, ok := tempIDs[tempID]
	return ok
}

func InvalidateTempID(tempID string) {
	delete(tempIDs, tempID)
}

func IsLoggedIn(r *http.Request) (string, *http.Cookie, bool) {
	c, err := r.Cookie(constants.SessionCookieName)
	if err != nil {
		log.Printf("No cookie found: %v", err)
		return "", nil, false
	}

	sessionID := c.Value

	us, ok := sessions[sessionID]
	if !ok {
		log.Println("No session found")
		return "", nil, false
	}
	if time.Now().After(us.Expiry) {
		log.Println("Session expired")
		delete(sessions, sessionID)
		return "", nil, false
	}

	_, idUUID := uuid.FromString(us.Username)
	if idUUID == nil {
		us.Expiry = time.Now().Add(constants.SessionTimeout)
	}

	log.Println("Is logged in")
	return us.Username, us.Cookie(), true
}
