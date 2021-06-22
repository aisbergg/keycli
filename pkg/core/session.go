package core

import (
	"net/url"
	"time"

	"github.com/Nerzal/gocloak/v8"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"
)

// -----------------------------------------------------------------------------
//
// Interfaces
//
// -----------------------------------------------------------------------------

// Session represents static session information, that are used to get access to
// resources of Keycloak.
type Session struct {
	Name       string      `json:"name"`
	URL        string      `json:"url"`
	Realm      string      `json:"realm"`
	ClientID   string      `json:"client_id"`
	Token      gocloak.JWT `json:"token"`
	Created    jwt.Time    `json:"created_at"`
	SkipVerify bool        `json:"skip_verify"`
}

// SessionService manages the creation and destruction of Keycloak sessions.
type SessionService interface {
	// Load loads a session from a stored session.
	Load(name string) (*Session, error)
	// LoadRefresh loads and refreshes a session.
	LoadRefresh(name string, refreshBeforeExpiry bool) (*Session, error)
	// CreateWithUsernamePassword creates a new session by loging into a session
	// provider with username and password. It also writes the newly created
	// session to a session repository.
	CreateWithUsernamePassword(name, url, realm, clientID, user, password string, skipVerify bool) (*Session, error)
	// CreateWithClientSecret creates a new session by loging into a session
	// provider with a client secret. It also writes the newly created session
	// to a session repository.
	CreateWithClientSecret(name, url, realm, clientID, secret string, skipVerify bool) (*Session, error)
	// Refresh refreshes a session, if it can be refreshed. Returns true, if the
	// access token was refreshed.
	Refresh(session *Session, beforeExpiry bool) (bool, error)
	// End ends the session and removes it from a session repository.
	End(session *Session, force bool) error
}

// SessionRepository is used for loading and storing from and to a repository.
// The access to the repository is managed to be exclusive between multiple
// processes.
type SessionRepository interface {
	// Exists indicates that a session with a given name exists or not.
	Exists(name string) (bool, error)
	// Open opens the session repository with exclusive access.
	Open(name string) error
	// Close closes the session repository.
	Close() error
	// Read reads the content of the session repository.
	Read() (*Session, error)
	// Write writes a session to the repository.
	Write(session *Session) error
	// Remove removes a stored session repository. Has no effect, if the file
	// doesn't exist.
	Remove(name string) error
}

// SessionProvider provides the means to create, refresh and end a session.
type SessionProvider interface {
	// CreateWithUsernamePassword creates a new session by logging into the
	// service provider using a username and password.
	CreateWithUsernamePassword(name, url, realm, clientID, user, password string, skipVerify bool) (*Session, error)
	// CreateWithUsernamePassword creates a new session by logging into the
	// service provider using a client secret.
	CreateWithClientSecret(name, url, realm, clientID, secret string, skipVerify bool) (*Session, error)
	// Logout logs out of the session provider and thereby ending a session.
	End(session *Session) error
	// Refresh refreshes an existing session.
	Refresh(session *Session) (bool, error)
}

// -----------------------------------------------------------------------------
//
// Implementation
//
// -----------------------------------------------------------------------------

// IsExpired returns true, if the access token is expired, else false. This doesn't mean it cannot be refreshed using the refresh token.
func (s *Session) IsExpired(beforeExpiry bool) bool {
	now := jwt.Now()
	accessTokenExpired := s.Created.Add(time.Second * time.Duration(s.Token.ExpiresIn)).After(now.Time)
	accessTokenExpiresSoon := s.Created.Add(time.Second * time.Duration(s.Token.ExpiresIn-60)).After(now.Time)
	return accessTokenExpired || (beforeExpiry && accessTokenExpiresSoon)
}

// CanBeRefreshed returns true, if the access token can be refreshed using the refresh token, else false.
func (s *Session) CanBeRefreshed() bool {
	now := jwt.Now()
	refreshTokenExpired := s.Created.Add(time.Second * time.Duration(s.Token.RefreshExpiresIn)).After(now.Time)
	return !refreshTokenExpired
}

// IsValid returns true, if the session object is indeed valid
func (s *Session) IsValid() bool {
	_, err := url.ParseRequestURI(s.URL)
	if err != nil ||
		s.URL == "" ||
		s.Realm == "" ||
		s.ClientID == "" ||
		s.Created.Equal(time.Time{}) ||
		s.Created.After(jwt.Now().Time) ||
		s.Token.ExpiresIn <= 0 ||
		s.Token.RefreshExpiresIn <= 0 ||
		s.Token.AccessToken == "" ||
		s.Token.RefreshToken == "" ||
		s.Token.TokenType != "Bearer" {
		return false
	}
	return true
}

type sessionService struct {
	repository SessionRepository
	provider   SessionProvider
}

// NewSessionService initializes a `SessionService`.
func NewSessionService(repo SessionRepository, provider SessionProvider) SessionService {
	return &sessionService{repository: repo, provider: provider}
}

func (ss *sessionService) Load(name string) (*Session, error) {
	if exists, _ := ss.repository.Exists(name); !exists {
		return nil, errors.Errorf("session '%s': does not exist", name)
	}

	// lock session repository for exclusive access
	if err := ss.repository.Open(name); err != nil {
		return nil, err
	}
	defer ss.repository.Close()

	// read the session information from the repository
	session, err := ss.repository.Read()
	if err != nil {
		return nil, errors.Wrapf(err, "session '%s': failed to retrieve from repository", name)
	}

	// check if the loaded session is valid
	if !session.IsValid() {
		return nil, errors.Errorf("session '%s': invalid. Login again to create a new session", name)
	}

	return session, nil
}

func (ss *sessionService) LoadRefresh(name string, refreshBeforeExpiry bool) (*Session, error) {
	if exists, _ := ss.repository.Exists(name); !exists {
		return nil, errors.Errorf("session '%s': does not exist", name)
	}

	// lock session repository for exclusive access
	if err := ss.repository.Open(name); err != nil {
		return nil, err
	}
	defer ss.repository.Close()

	// read the session information from the repository
	session, err := ss.repository.Read()
	if err != nil {
		return nil, errors.Wrapf(err, "session '%s': failed to retrieve from repository", name)
	}

	// check if the loaded session is valid
	if !session.IsValid() {
		return nil, errors.Errorf("session '%s': invalid. Login again to create a new session", name)
	}

	// refresh the access token
	refreshed, err := ss.Refresh(session, refreshBeforeExpiry)
	if err != nil {
		return nil, err
	}

	// write any changes to the session file
	if refreshed {
		if err := ss.repository.Write(session); err != nil {
			return nil, errors.Wrapf(err, "session '%s': failed to write to repository", name)
		}
	}

	return session, nil
}

func (ss *sessionService) CreateWithUsernamePassword(name, url, realm, clientID, user, password string, skipVerify bool) (*Session, error) {
	createFunc := func() (*Session, error) {
		return ss.provider.CreateWithUsernamePassword(name, url, realm, clientID, user, password, skipVerify)
	}
	return ss.create(name, createFunc)
}

func (ss *sessionService) CreateWithClientSecret(name, url, realm, clientID, secret string, skipVerify bool) (*Session, error) {
	createFunc := func() (*Session, error) {
		return ss.provider.CreateWithClientSecret(name, url, realm, clientID, secret, skipVerify)
	}
	return ss.create(name, createFunc)
}

func (ss *sessionService) create(name string, createFunc func() (*Session, error)) (*Session, error) {
	// lock session repository for exclusive access
	if err := ss.repository.Open(name); err != nil {
		return nil, errors.Wrapf(err, "session '%s': failed to open repository", name)
	}
	defer ss.repository.Close()

	// create a new session by requesting access and refresh tokens
	session, err := createFunc()
	if err != nil {
		return nil, errors.Wrapf(err, "session '%s': failed to login", name)
	}

	// write session information into the repository
	if err := ss.repository.Write(session); err != nil {
		return nil, errors.Wrapf(err, "session '%s': failed to write session information", name)
	}

	return session, nil
}

func (ss *sessionService) Refresh(session *Session, beforeExpiry bool) (bool, error) {
	if !session.IsExpired(beforeExpiry) {
		return false, nil
	}
	if !session.CanBeRefreshed() {
		return false, errors.Errorf("session '%s': is expired. Login again to create a new session", session.Name)
	}

	refreshed, err := ss.provider.Refresh(session)
	if err != nil {
		return false, errors.Wrapf(err, "session '%s': failed to refresh", session.Name)
	}

	return refreshed, nil
}

func (ss *sessionService) End(session *Session, force bool) error {
	if !session.IsValid() || !session.CanBeRefreshed() {
		return ss.repository.Remove(session.Name)
	}

	err := ss.provider.End(session)
	if err != nil && !force {
		return errors.Errorf("session '%s': failed to logout: %v", session.Name, err)
	}

	return ss.repository.Remove(session.Name)
}
