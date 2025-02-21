package main

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	envOrgID     = "SCW_DEFAULT_ORGANIZATION_ID"
	envAccessKey = "SCW_ACCESS_KEY"
	envSecretKey = "SCW_SECRET_KEY"
	envProjectID = "SCW_PROJECT_ID"

	envNTagsToKeep = "SCW_NUMBER_VERSIONS_TO_KEEP"

	// If set to true, older tags will be deleted.
	envNoDryRun = "SCW_NO_DRY_RUN"
)

// Check for mandatory variables before starting to work.
func init() {
	mandatoryVariables := [...]string{envOrgID, envAccessKey, envSecretKey, envProjectID, envNTagsToKeep}

	for idx := range mandatoryVariables {
		if os.Getenv(mandatoryVariables[idx]) == "" {
			panic("missing environment variable " + mandatoryVariables[idx])
		}
	}
}

func main() {
	slog.Info("cleaning container registry tags...")

	// Create a Scaleway client with credentials from environment variables.
	client, err := scw.NewClient(
		// Get your organization ID at https://console.scaleway.com/organization/settings
		scw.WithDefaultOrganizationID(os.Getenv(envOrgID)),

		// Get your credentials at https://console.scaleway.com/iam/api-keys
		scw.WithAuth(os.Getenv(envAccessKey), os.Getenv(envSecretKey)),

		// Get more about our availability
		// zones at https://www.scaleway.com/en/docs/console/my-account/reference-content/products-availability/
		scw.WithDefaultRegion(scw.RegionFrPar),
	)
	if err != nil {
		panic(err)
	}

	regAPI := NewRegistryAPI(client, os.Getenv(scw.ScwDefaultProjectIDEnv))

	numberTagsToKeep, err := strconv.Atoi(os.Getenv(envNTagsToKeep))
	if err != nil {
		panic(err)
	}

	tagsToDelete, err := regAPI.GetTagsAfterNVersions(numberTagsToKeep)
	if err != nil {
		panic(err)
	}

	dryRun := true
	noDryRunEnv := os.Getenv(envNoDryRun)

	if strings.EqualFold(noDryRunEnv, "true") {
		dryRun = false
	}

	if err := regAPI.DeleteTags(tagsToDelete, dryRun); err != nil {
		panic(err)
	}
}
