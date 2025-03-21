package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	envOrgID     = "SCW_DEFAULT_ORGANIZATION_ID" // Scaleway organization ID
	envAccessKey = "SCW_ACCESS_KEY"              // Scaleway API access key
	envSecretKey = "SCW_SECRET_KEY"              // Scaleway API secret key
	envProjectID = "SCW_PROJECT_ID"              // Scaleway project ID

	envRegion               = "SCW_REGION"
	envDatabaseID           = "SCW_RDB_ID"
	envBackupExpirationDays = "SCW_EXPIRATION_DAYS"
)

// Check for mandatory variables before starting to work.
func init() {
	// Slice of environmental variables that must be set for the application to run
	mandatoryVariables := [...]string{envOrgID, envAccessKey, envSecretKey, envProjectID, envRegion}

	// Iterate through the slice and check if any variables are not set
	for idx := range mandatoryVariables {
		if os.Getenv(mandatoryVariables[idx]) == "" {
			panic("missing environment variable " + mandatoryVariables[idx])
		}
	}
}

func main() {
	fmt.Println("creating backup of managed database...")

	// Create a Scaleway client with credentials provided via environment variables.
	// The client is used to interact with the Scaleway API
	client, err := scw.NewClient(
		// Get your organization ID at https://console.scaleway.com/organization/settings
		scw.WithDefaultOrganizationID(os.Getenv(envOrgID)),

		// Get your credentials at https://console.scaleway.com/iam/api-keys
		scw.WithAuth(os.Getenv(envAccessKey), os.Getenv(envSecretKey)),

		// Get more about our availability
		// zones at https://www.scaleway.com/en/docs/console/my-account/reference-content/products-availability/
		scw.WithDefaultRegion(scw.Region(os.Getenv(envRegion))),
	)
	if err != nil {
		panic(err)
	}

	rdbAPI := rdb.NewAPI(client)

	if err := createRdbSnapshot(rdbAPI); err != nil {
		panic(err)
	}
}

func createRdbSnapshot(rdbAPI *rdb.API) error {
	rdbInstance, err := rdbAPI.GetInstance(&rdb.GetInstanceRequest{
		Region:     scw.Region(scw.Region(os.Getenv(envRegion))),
		InstanceID: os.Getenv(envDatabaseID),
	})
	if err != nil {
		return fmt.Errorf("error while getting database instance %w", err)
	}

	databasesList, err := rdbAPI.ListDatabases(&rdb.ListDatabasesRequest{
		Region:     scw.Region(os.Getenv(envRegion)),
		InstanceID: rdbInstance.ID,
	})
	if err != nil {
		return fmt.Errorf("error while listing databases %w", err)
	}

	expiresAt, err := getExpirationDate()
	if err != nil {
		return fmt.Errorf("error while getting expiration date %w", err)
	}

	now := time.Now()

	for _, database := range databasesList.Databases {
		backupName := fmt.Sprintf("backup-%s-%s-%s",
			database.Name,
			now,
			os.Getenv(envRegion))

		backup, err := rdbAPI.CreateDatabaseBackup(&rdb.CreateDatabaseBackupRequest{
			Region:       scw.Region(os.Getenv(envRegion)),
			InstanceID:   rdbInstance.ID,
			Name:         backupName,
			DatabaseName: database.Name,
			ExpiresAt:    expiresAt,
		})
		if err != nil {
			return fmt.Errorf("error while creating database backup request %w", err)
		}

		fmt.Println("Created backup ", backup.Name)
	}

	return nil
}

func getExpirationDate() (*time.Time, error) {
	var expiresAt *time.Time
	expireDays := os.Getenv(envBackupExpirationDays)

	if expireDays != "" {
		expireDaysInt, err := strconv.Atoi(expireDays)
		if err != nil {
			return nil, fmt.Errorf("error while getting %w", err)
		}

		if expireDaysInt > 0 {
			expiration := time.Now().AddDate(0, 0, expireDaysInt)
			expiresAt = &expiration
		}
	}

	return expiresAt, nil
}
