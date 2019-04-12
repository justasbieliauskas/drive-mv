package fs

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

type credentials map[string]string

func credsFromEnv(env []string) credentials {
	creds := make(map[string]string)
	for _, str := range env {
		pair := strings.Split(str, "=")
		creds[pair[0]] = pair[1]
	}
	return creds
}

func (creds credentials) json() ([]byte, error) {
	type info struct {
		ClientID     string   `json:"client_id"`
		ProjectID    string   `json:"project_id"`
		AuthURI      string   `json:"auth_uri"`
		TokenURI     string   `json:"token_uri"`
		CertURL      string   `json:"auth_provider_x509_cert_url"`
		Secret       string   `json:"client_secret"`
		RedirectURIs []string `json:"redirect_uris"`
	}
	bytes, err := json.Marshal(map[string]info{
		"installed": info{
			ClientID:  creds["DRIVE_CLIENT_ID"],
			ProjectID: creds["DRIVE_PROJECT_ID"],
			AuthURI:   "https://accounts.google.com/o/oauth2/auth",
			TokenURI:  "https://oauth2.googleapis.com/token",
			CertURL:   "https://www.googleapis.com/oauth2/v1/certs",
			Secret:    creds["DRIVE_CLIENT_SECRET"],
			RedirectURIs: []string{
				"urn:ietf:wg:oauth:2.0:oob",
				"http://localhost",
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal credentials json: %v", err)
	}
	return bytes, nil
}

func (creds credentials) token() (*oauth2.Token, error) {
	expiry, err := time.Parse(
		time.RFC3339,
		creds["DRIVE_TOKEN_EXPIRY"],
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse expiry date: %v", err)
	}
	token := &oauth2.Token{
		AccessToken:  creds["DRIVE_ACCESS_TOKEN"],
		TokenType:    "Bearer",
		RefreshToken: creds["DRIVE_REFRESH_TOKEN"],
		Expiry:       expiry,
	}
	return token, nil
}
