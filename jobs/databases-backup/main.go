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
	// Defining variable or secret readings.
	VAR_ORG_ID = "SCW_DEFAULT_ORGANIZATION_ID"
	VAR_AK     = "SCW_ACCESS_KEY"
	VAR_SK     = "SCW_SECRET_KEY"
	VAR_REGION = "REGION"
	VAR_RDB_ID = "INSTANCE_ID"

	// optional, never expires if not definied.
	VAR_EXPIRE_AT_DAYS = "EXPIRE_AT_DAYS"
)

func main() {
	fmt.Println("creating backup of managed database...")

	// Create a Scaleway client with credentials from environment variables.
	client, err := scw.NewClient(
		// Get your organization ID at https://console.scaleway.com/organization/settings
		scw.WithDefaultOrganizationID(os.Getenv(VAR_ORG_ID)),

		// Get your credentials at https://console.scaleway.com/iam/api-keys
		scw.WithAuth(os.Getenv(VAR_AK), os.Getenv(VAR_SK)),

		// Get more about our availability zones at https://www.scaleway.com/en/docs/console/my-account/reference-content/products-availability/
		scw.WithDefaultRegion(scw.RegionFrPar),
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
		Region:     scw.Region(os.Getenv(VAR_REGION)),
		InstanceID: os.Getenv(VAR_RDB_ID),
	})
	if err != nil {
		return fmt.Errorf("error while getting database instance %w", err)
	}

	databasesList, err := rdbAPI.ListDatabases(&rdb.ListDatabasesRequest{
		Region:     scw.Region(os.Getenv(VAR_REGION)),
		InstanceID: rdbInstance.ID,
	})
	if err != nil {
		return fmt.Errorf("error while listing databases %w", err)
	}

	expiresAt, err := getExpirationDate()
	if err != nil {
		return fmt.Errorf("error while getting expiration date %w", err)
	}

	tn := time.Now()
	backupName := fmt.Sprintf("backup_%s_%d%d%d", rdbInstance.Name, tn.Year(), tn.Month(), tn.Day())

	for _, database := range databasesList.Databases {

		backup, err := rdbAPI.CreateDatabaseBackup(&rdb.CreateDatabaseBackupRequest{
			Region:       scw.Region(os.Getenv(VAR_REGION)),
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
	expireDays := os.Getenv(VAR_EXPIRE_AT_DAYS)

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

func init() {
	if os.Getenv(VAR_ORG_ID) == "" {
		panic("missing " + VAR_ORG_ID)
	}

	if os.Getenv(VAR_AK) == "" {
		panic("missing " + VAR_AK)
	}

	if os.Getenv(VAR_SK) == "" {
		panic("missing " + VAR_SK)
	}

	if os.Getenv(VAR_RDB_ID) == "" {
		panic("missing " + VAR_RDB_ID)
	}

	if os.Getenv(VAR_REGION) == "" {
		panic("missing " + VAR_REGION)
	}
}
