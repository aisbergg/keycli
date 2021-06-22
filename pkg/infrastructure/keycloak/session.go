package keycloak

import (
	"github.com/Nerzal/gocloak/v8"
	"github.com/aisbergg/keycli/pkg/core"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"
)

// keyclaokSessionProvider implements `core.SessionProvider`
type keyclaokSessionProvider struct{}

// NewKeycloakSessionProvider initializes a new `keyclaokSessionProvider`.
func NewKeycloakSessionProvider() core.SessionProvider {
	return &keyclaokSessionProvider{}
}

func (sp *keyclaokSessionProvider) CreateWithUsernamePassword(name, url, realm, clientID, user, password string, skipVerify bool) (*core.Session, error) {
	topt := gocloak.TokenOptions{
		ClientID:  gocloak.StringP(clientID),
		GrantType: gocloak.StringP("password"),
		Username:  &user,
		Password:  &password,
	}
	return sp.create(name, url, realm, skipVerify, topt)
}

func (sp *keyclaokSessionProvider) CreateWithClientSecret(name, url, realm, clientID, secret string, skipVerify bool) (*core.Session, error) {
	topt := gocloak.TokenOptions{
		ClientID:     &clientID,
		ClientSecret: &secret,
		GrantType:    gocloak.StringP("client_credentials"),
	}
	return sp.create(name, url, realm, skipVerify, topt)
}

// create sends an auth request to Keycloak and returns a new session.
func (sp *keyclaokSessionProvider) create(name, url, realm string, skipVerify bool, tokenOptions gocloak.TokenOptions) (*core.Session, error) {
	ctx, cancel := createContext()
	defer cancel()
	gocloakClient := createGoclaokClient(url, skipVerify)

	token, err := (*gocloakClient).GetToken(ctx, realm, tokenOptions)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get token")
	}

	session := core.Session{
		ClientID:   *tokenOptions.ClientID,
		Token:      *token,
		Created:    *jwt.Now(),
		Name:       name,
		URL:        url,
		Realm:      realm,
		SkipVerify: skipVerify,
	}
	return &session, nil
}

// End ends the session.
func (sp *keyclaokSessionProvider) End(session *core.Session) error {
	ctx, cancel := createContext()
	defer cancel()
	gocloakClient := createGoclaokClient(session.URL, session.SkipVerify)

	err := (*gocloakClient).Logout(ctx, session.ClientID, "", session.Realm, session.Token.RefreshToken)
	if err != nil {
		return errors.Wrap(err, "failed to logout")
	}

	return nil
}

func (sp *keyclaokSessionProvider) Refresh(session *core.Session) (bool, error) {
	ctx, cancel := createContext()
	defer cancel()
	gocloakClient := createGoclaokClient(session.URL, session.SkipVerify)

	token, err := (*gocloakClient).GetToken(ctx, session.Realm, gocloak.TokenOptions{
		ClientID:     &session.ClientID,
		GrantType:    gocloak.StringP("refresh_token"),
		RefreshToken: &session.Token.RefreshToken,
	})
	if err != nil {
		return false, errors.Wrap(err, "failed to refresh token")
	}

	// update token information
	now := jwt.Now()
	session.Created = *now
	session.Token = *token

	return true, nil
}
