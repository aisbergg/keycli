package cli

import (
	"fmt"

	"github.com/aisbergg/keycli/pkg/core"
	"github.com/aisbergg/keycli/pkg/infrastructure/jsonfile"
	"github.com/aisbergg/keycli/pkg/infrastructure/keycloak"
)

// Login is the implementation of the login command.
func Login(name, url, realm, clientID, secretKey, user, password string, skipVerify bool) error {
	sessionRepository := jsonfile.NewJSONFileSessionRepository()
	sessionProvider := keycloak.NewKeycloakSessionProvider()
	sessionService := core.NewSessionService(sessionRepository, sessionProvider)

	var err error
	if secretKey != "" {
		_, err = sessionService.CreateWithClientSecret(name, url, realm, clientID, secretKey, skipVerify)
	} else {
		_, err = sessionService.CreateWithUsernamePassword(name, url, realm, clientID, user, password, skipVerify)
	}

	if err != nil {
		return err
	}
	fmt.Printf("Created session '%s'.\nYour session was stored unencrypted in %s\n"+
		"When you are done, you can end the session by using the 'logout' command.\n",
		name, jsonfile.PathFromName(name))

	return nil
}
