package secret

import (
	"context"
	"fmt"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"gopkg.in/yaml.v3"
)

var projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")

func SetProjectID(projID string) {
	projectID = projID
}

func GetProjectID() string {
	return projectID
}

type Secret struct {
	data []byte
}

func (sec *Secret) Bytes() []byte {
	return sec.data
}

func (sec *Secret) YamlDecode(out interface{}) error {
	return yaml.Unmarshal(sec.data, out)
}

func GetSecret(name string, version string) (*Secret, error) {
	var hasil Secret
	var err error

	// Create the client.
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	// GCP project in which to store secrets in Secret Manager.

	secretName := fmt.Sprintf("projects/%s/secrets/%s/versions/%s", projectID, name, version)
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: secretName,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return &hasil, err
	}

	hasil.data = result.Payload.Data
	return &hasil, err
}
