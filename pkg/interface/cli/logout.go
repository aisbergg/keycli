package cli

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/aisbergg/keycli/pkg/core"
	"github.com/aisbergg/keycli/pkg/infrastructure/jsonfile"
	"github.com/aisbergg/keycli/pkg/infrastructure/keycloak"
)

// Logout is the implementation of the logout command.
func Logout(name string, force bool) error {
	sessionRepository := jsonfile.NewJSONFileSessionRepository()
	sessionProvider := keycloak.NewKeycloakSessionProvider()
	sessionService := core.NewSessionService(sessionRepository, sessionProvider)
	session, err := sessionService.Load(name)
	if err != nil {
		return errors.Wrap(err, "Failed to load session")
	}
	if err := sessionService.End(session, force); err != nil {
		return errors.Wrap(err, "Failed to end session")
	}
	fmt.Printf("Ended '%s' session and removed login credentials\n", name)

	return nil
}
