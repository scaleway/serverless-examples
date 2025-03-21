package main

import (
	"log/slog"
	"os"
	"strconv"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Constants for environment variable names used to configure the application
const (
	envOrgID     = "SCW_DEFAULT_ORGANIZATION_ID" // Scaleway organization ID
	envAccessKey = "SCW_ACCESS_KEY"              // Scaleway API access key
	envSecretKey = "SCW_SECRET_KEY"              // Scaleway API secret key
	envProjectID = "SCW_PROJECT_ID"              // Scaleway project ID

	envNTagsToKeep = "SCW_NUMBER_VERSIONS_TO_KEEP" // Number of container registry tags to keep
	// If set to "true", older tags will be deleted.
	// Otherwise, only a dry run will be performed
	envNoDryRun = "SCW_NO_DRY_RUN"
)

// Check for mandatory variables before starting to work.
func init() {
	// Slice of environmental variables that must be set for the application to run
	mandatoryVariables := [...]string{envOrgID, envAccessKey, envSecretKey, envProjectID, envNTagsToKeep}

	// Iterate through the slice and check if any variables are not set
	for idx := range mandatoryVariables {
		if os.Getenv(mandatoryVariables[idx]) == "" {
			panic("missing environment variable " + mandatoryVariables[idx])
		}
	}
}

func main() {
	slog.Info("cleaning container registry tags...")

	// Create a Scaleway client with credentials provided via environment variables.
	// The client is used to interact with the Scaleway API
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

	// Create a new instance of RegistryAPI, passing the Scaleway client and the project ID
	// RegistryAPI is a custom interface for interacting with the Scaleway container registry
	regAPI := NewRegistryAPI(client, os.Getenv(scw.ScwDefaultProjectIDEnv))

	// Parse the number of tags to keep from the environment variable, which should be a string
	numberTagsToKeep, err := strconv.Atoi(os.Getenv(envNTagsToKeep))
	if err != nil {
		panic(err)
	}

	// Get the tags to delete by specifying the number of tags to keep
	// The function returns both the tags to be deleted and an error, if any
	tagsToDelete, err := regAPI.GetTagsAfterNVersions(numberTagsToKeep)
	if err != nil {
		panic(err)
	}

	// Determine whether to perform a dry run or delete the tags
	// Default behavior is to perform a dry run (no deletion)
	dryRun := true
	noDryRunEnv := os.Getenv(envNoDryRun)

	// If the SCW_NO_DRY_RUN environment variable is set to "true",
	// the tags will be deleted; otherwise, only a dry run will be performed
	if strings.EqualFold(noDryRunEnv, "true") {
		dryRun = false
	}

	// Delete the tags or perform a dry run, depending on the dryRun flag
	if err := regAPI.DeleteTags(tagsToDelete, dryRun); err != nil {
		panic(err)
	}
}
