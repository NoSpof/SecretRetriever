package retriever

import (
	"context"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"google.golang.org/api/option"
)

func CreateSecretClient() *secretmanager.Client {

	client, err := secretmanager.NewClient(context.Background(), option.WithCredentialsFile("./account.json"))
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	return client
}
