package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"github.com/nospof/secretretriever/kate"
	"github.com/nospof/secretretriever/retriever"
	"github.com/nospof/secretretriever/tools"
	"google.golang.org/api/iterator"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	// Init env var
	var projectID = ""
	var namespace = "default"
	var app_name = ""
	var env = "stg"
	// Init Kubernetes client
	kate := kate.GetKubeConfig()
	// Init Secret Client
	secretClient := retriever.CreateSecretClient()
	// set var from env var
	if len(os.Getenv("PROJECT_ID")) != 0 {
		projectID = os.Getenv("PROJECT_ID")
	}
	if len(os.Getenv("SECRET_NAMESPACE")) != 0 {
		namespace = os.Getenv("SECRET_NAMESPACE")
	}
	if len(os.Getenv("APP_NAME")) != 0 {
		app_name = os.Getenv("APP_NAME")
	}
	if len(os.Getenv("ENV")) != 0 {
		env = os.Getenv("ENV")
	}
	// Add ObjectMeta for kubernetes
	var objectMeta metav1.ObjectMeta
	objectMeta.Name = app_name + "-secrets-" + env
	objectMeta.Namespace = namespace
	data := make(map[string][]byte)
	var filter = "labels.app=" + app_name
	fmt.Println(filter)
	// Build the request.
	req := &secretmanagerpb.ListSecretsRequest{
		Parent: "projects/" + projectID,
		//	Filter: filter,
	}
	// Call the API.
	it := secretClient.ListSecrets(context.Background(), req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Panic(err)
		}
		log.Println("Get secret from : " + resp.Name)
		// Get version from secret
		reqVersions := &secretmanagerpb.ListSecretVersionsRequest{
			Parent: resp.Name,
		}
		it := secretClient.ListSecretVersions(context.Background(), reqVersions)
		for {
			resp, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Println(err)
			}
			if resp.State.String() == "ENABLED" {
				log.Println("The " + resp.Name + " is : " + resp.State.String())
				reqVerion := &secretmanagerpb.AccessSecretVersionRequest{
					Name: resp.Name,
				}
				result, err := secretClient.AccessSecretVersion(context.Background(), reqVerion)
				tools.CheckIfError(err)
				// Check if checksum is correct
				retriever.CheckChecksum(result.Payload)
				secretData := strings.TrimSuffix(string(result.Payload.Data), "\r\n")
				tools.CheckTheEndline(secretData)
				data[strings.Split(resp.Name, "/")[3]] = []byte(secretData)
			}

		}

	}
	// Construct the secret
	payloadSecret := &v1.Secret{
		Data:       data,
		ObjectMeta: objectMeta,
	}
	checkifexist, err := kate.CoreV1().Secrets(namespace).Get(context.Background(), app_name+"-secrets-"+env, metav1.GetOptions{})
	if len(checkifexist.Name) != 0 {
		err := kate.CoreV1().Secrets(namespace).Delete(context.Background(), app_name+"-secrets-"+env, metav1.DeleteOptions{})
		tools.CheckIfError(err)

	}
	secret, err := kate.CoreV1().Secrets(namespace).Create(context.TODO(), payloadSecret, metav1.CreateOptions{})
	tools.CheckIfError(err)
	log.Println(secret.ObjectMeta.Name + " Has been created ")
	defer secretClient.Close()
}
