// Package keycloak provides an implementation of the application repository
// interfaces.
package keycloak

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/Nerzal/gocloak/v8"
	"github.com/aisbergg/keycli/pkg/core"
)

const timeout = 15 * time.Second

type client struct {
	gocloakClient *gocloak.GoCloak
	session       *core.Session
}

func NewClient(session core.Session) *client {
	gocloakClient := createGoclaokClient(session.URL, session.SkipVerify)
	return &client{gocloakClient: gocloakClient, session: &session}
}

func createGoclaokClient(url string, skipVerify bool) *gocloak.GoCloak {
	gocloakClient := gocloak.NewClient(url)
	if skipVerify {
		restyClient := gocloakClient.RestyClient()
		restyClient.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	}
	return &gocloakClient
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
