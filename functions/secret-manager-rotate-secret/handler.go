package secretmanagerrotatesecret

import (
	"encoding/json"
	"fmt"
	"net/http"
	"secret-manager-rotate-secret/random"

	rdb "github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type rotateRequest struct {
	SecretID      string `json:"secret_id"`
	RDBInstanceID string `json:"rdb_instance_id"`
}

type databaseCredentials struct {
	Engine   string `json:"engine"`
	Username string `json:"username"`
	Password string `json:"password"`
	Hostname string `json:"host"`
	DBName   string `json:"dbname"`
	Port     string `json:"port"`
}

func Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req rotateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	client, err := scw.NewClient(scw.WithEnv())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rdbApi := rdb.NewAPI(client)
	secretApi := secret.NewAPI(client)

	// access current secret version to get revision and payload
	currentVersion, err := secretApi.AccessSecretVersion(&secret.AccessSecretVersionRequest{
		SecretID: req.SecretID,
		Revision: "latest_enabled",
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// generate new password
	newPassword, err := random.CreateString(random.StringParams{
		Length:     16,
		Upper:      true,
		MinUpper:   1,
		Lower:      true,
		MinLower:   1,
		Numeric:    true,
		MinNumeric: 1,
		Special:    true,
		MinSpecial: 1,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// deserialize secret payload to access values
	var payload databaseCredentials
	err = json.Unmarshal(currentVersion.Data, &payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update RDB user with new password
	_, err = rdbApi.UpdateUser(&rdb.UpdateUserRequest{
		InstanceID: req.RDBInstanceID,
		Password:   scw.StringPtr(string(newPassword)),
		Name:       payload.Username,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	payload.Password = string(newPassword)

	newData, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create new version of the secret with same content and password updated
	_, err = secretApi.CreateSecretVersion(&secret.CreateSecretVersionRequest{
		SecretID:        req.SecretID,
		Data:            newData,
		DisablePrevious: scw.BoolPtr(true),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "database credentials updated")
}
