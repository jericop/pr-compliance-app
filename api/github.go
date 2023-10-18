package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

func (server *Server) validatWebhookRequest(r *http.Request) (interface{}, error) {
	// Validate payload from request using webhook secret
	payload, err := github.ValidatePayload(r, []byte(server.githubWebhookSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to validate payload: %v", err)
	}

	// Parse event from payload
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse webhook into event: %v", err)
	}

	return event, nil
}

func (server *Server) newSignedHttpClient(ctx context.Context) (*http.Client, error) {
	now := time.Now()

	// Token expires after 10 minutes
	claims := jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
		Issuer:    server.githubAppId,
	}

	// Sign token with jwt
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	jwtToken, err := token.SignedString(server.githubPrivateKey)
	if err != nil {
		log.Fatalf("Failed to generate JWT token: %v", err)
		return &http.Client{}, fmt.Errorf("failed to generate JWT token: %v", err)
	}

	oauthTransport := oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: jwtToken,
		}),
	)

	return oauthTransport, nil
}

// Create and return new GitHub app client
func (server *Server) newGithubAppClient(ctx context.Context) (*github.Client, error) {
	oauthTransport, err := server.newSignedHttpClient(ctx)
	if err != nil {
		return &github.Client{}, err
	}

	return github.NewClient(oauthTransport), nil
}

// Create and return new GitHub installation client
func newGithubAppInstallationClient(ctx context.Context, appClient *github.Client, installationID int64) (*github.Client, error) {
	installationToken, _, err := appClient.Apps.CreateInstallationToken(ctx, installationID, nil)
	if err != nil {
		return &github.Client{}, fmt.Errorf("Failed to create installation token: %v", err)
	}

	installationTransport := oauth2.NewClient(
		ctx,
		oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: installationToken.GetToken(),
		}),
	)

	return github.NewClient(installationTransport), nil
}

type githubFactoryInterface interface {
	NewAppClient(context.Context) (*github.Client, error)
	NewInstallationClient(context.Context, int64) (*github.Client, error)
	ValidatWebhookRequest(*http.Request) (interface{}, error)
}

type githubFactory struct {
	server *Server
}

func NewGithubFactory(server *Server) *githubFactory {
	return &githubFactory{server: server}
}

func (f *githubFactory) NewAppClient(ctx context.Context) (*github.Client, error) {
	appClient, err := f.server.newGithubAppClient(ctx)
	if err != nil {
		return &github.Client{}, err
	}
	return appClient, nil
}

func (f *githubFactory) NewInstallationClient(ctx context.Context, installationID int64) (*github.Client, error) {
	appClient, err := f.server.newGithubAppClient(ctx)
	if err != nil {
		return &github.Client{}, err
	}
	appInstallationClient, err := newGithubAppInstallationClient(ctx, appClient, installationID)
	if err != nil {
		return &github.Client{}, err
	}
	return appInstallationClient, nil
}

func (f *githubFactory) ValidatWebhookRequest(r *http.Request) (interface{}, error) {
	// Validate payload from request using webhook secret
	payload, err := github.ValidatePayload(r, []byte(f.server.githubWebhookSecret))
	if err != nil {
		return nil, fmt.Errorf("failed to validate payload: %v", err)
	}

	// Parse event from payload
	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		return nil, fmt.Errorf("failed to parse webhook into event: %v", err)
	}

	return event, nil
}
