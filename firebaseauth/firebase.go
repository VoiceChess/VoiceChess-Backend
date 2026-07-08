package firebaseauth

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context) (*auth.Client, error) {
	if serviceAccountJSON := os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON"); serviceAccountJSON != "" {
		app, err := firebase.NewApp(ctx, nil, option.WithCredentialsJSON([]byte(serviceAccountJSON)))
		if err != nil {
			return nil, fmt.Errorf("firebase.NewApp WithCredentialsJSON: %w", err)
		}
		return app.Auth(ctx)
	}

	credentialsFile := os.Getenv("FIREBASE_CREDENTIALS_FILE")
	if credentialsFile == "" {
		credentialsFile = "firebase-service-account.json"
	}

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("firebase.NewApp WithCredentialsFile: %w", err)
	}

	return app.Auth(ctx)
}
