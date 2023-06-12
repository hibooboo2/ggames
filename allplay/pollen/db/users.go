package db

import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
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

func (s *UserSession) Cookie() *http.Cookie {
	return &http.Cookie{
		Name:    constants.SessionCookieName,
		Value:   s.ID,
		Expires: s.Expiry,
		Path:    "/",
	}
}

func init() {
	RegisterUser("jj@jhrb.us", "jj", "jj")
	RegisterUser("ff@jhrb.us", "ff", "ff")
}

func (s *UserSession) MustReauthenticate() bool {
	if time.Now().After(s.Expiry) {
		return true
	}
	_, idUUID := uuid.FromString(s.Username)
	if idUUID == nil {
		s.Expiry = time.Now().Add(constants.SessionTimeout)
	}
	return false
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

func Login(username string, password string, timeout time.Duration, newSession bool) (*UserSession, error) {
	if !CheckPassword(username, password) {
		return nil, errors.New("invalid username or password")
	}

	if !newSession {
		for _, s := range sessions {
			if s.Username == username {
				if s.MustReauthenticate() {
					delete(sessions, s.ID)
					return nil, errors.New("must reauthenticate")
				}
				s.Expiry = time.Now().Add(timeout)
				return s, nil
			}
		}
		return nil, errors.New("session not found")
	}

	sessionID := uuid.Must(uuid.NewV4()).String()
	sessions[sessionID] = &UserSession{
		ID:       sessionID,
		Username: username,
		Expiry:   time.Now().Add(timeout),
	}

	return sessions[sessionID], nil
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

func tmpLogin(r *http.Request) (*UserSession, bool) {
	tempID := r.FormValue("anonymous")
	if tempID != "" && IsTempID(tempID) {
		InvalidateTempID(tempID)
		log.Println("No active session and tempID found, creating temporary session")
		tmpID := uuid.Must(uuid.NewV4())
		err := RegisterUser(fmt.Sprintf("%s@jhrb.us", tmpID.String()), tmpID.String(), tmpID.String())
		if err != nil {
			log.Println("FAiled to create temporary session: ", err)
			return nil, false
		}
		us, err := Login(tmpID.String(), tmpID.String(), time.Hour*24*7, true)
		if err != nil {
			log.Println("Failed to login temporary session: ", err)
			return nil, false
		}
		return us, true
	}
	return nil, false
}

func loginBasicAuth(r *http.Request, newSession bool) (*UserSession, bool) {
	username, password, ok := r.BasicAuth()
	log.Println("Basicauth: ", username)
	if !ok {
		log.Printf("No basic auth found")
		return nil, false
	}

	log.Println("Attempting to login with basic auth: ", username)
	log.Println("password: ", password)
	us, err := Login(username, password, constants.SessionTimeout, newSession)
	if err != nil {
		log.Printf("Login failed: %q %v", username, err)
		return nil, false
	}
	return us, true
}

func hasCookieSession(r *http.Request) (*UserSession, bool) {
	c, err := r.Cookie(constants.SessionCookieName)
	if err != nil {
		log.Printf("No cookie found: %v", err)
		return nil, false
	}

	sessionID := c.Value

	us, ok := sessions[sessionID]
	if !ok {
		log.Println("No session found")
		return nil, false
	}
	if time.Now().After(us.Expiry) {
		log.Println("Session expired")
		delete(sessions, sessionID)
		return nil, false
	}

	_, idUUID := uuid.FromString(us.Username)
	if idUUID == nil {
		us.Expiry = time.Now().Add(constants.SessionTimeout)
	}

	log.Println("Is logged in")
	return us, true
}

func IsLoggedIn(r *http.Request, newSession bool) (string, *http.Cookie, bool) {
	us, ok := hasCookieSession(r)
	if !ok {
		//basic auth
		us, ok = loginBasicAuth(r, newSession)
		if !ok {
			//tmp auth
			us, ok = tmpLogin(r)
		}
	}

	if !ok {
		//no auth
		return "", nil, false
	}

	return us.Username, us.Cookie(), true
}
