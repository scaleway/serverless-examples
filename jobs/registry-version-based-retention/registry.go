package main

import (
	"fmt"
	"log/slog"
	"strings"

	registry "github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// RegistryAPI represents a Scaleway Container Registry accessor and extends
// capabilities to clean images.
type RegistryAPI struct {
	// regClient Scaleway Container Registry accessor
	regClient *registry.API

	// projectID if null, all projects will be checked.
	projectID *string

	// disableProtection if set to true (DANGEROUS) it will delete images potentially used
	// in Serverless Jobs, Functions and Containers.
	disableProtection bool
}

// NewRegistryAPI used to create a new RegistryAPI to manage the Scaleway Container Registry API.
func NewRegistryAPI(client *scw.Client, projectID string) *RegistryAPI {
	// if projectID is empty, no project ID will be passed to Scaleway SDK to use default settings. Generally is to apply
	// on all projects.
	var ptrProjectID *string
	if projectID != "" {
		ptrProjectID = &projectID
	}

	return &RegistryAPI{
		regClient:         registry.NewAPI(client),
		projectID:         ptrProjectID,
		disableProtection: false,
	}
}

func (r *RegistryAPI) GetTagsAfterNVersions(numberVersionsToKeep int) ([]*registry.Tag, error) {
	images, err := r.regClient.ListImages(&registry.ListImagesRequest{ProjectID: r.projectID}, scw.WithAllPages())
	if err != nil {
		return nil, fmt.Errorf("error listing container images %w", err)
	}

	if numberVersionsToKeep <= 1 {
		return nil, fmt.Errorf("number of versions to keep <= 1 is dangereous")
	}

	tagsToDelete := make([]*registry.Tag, 0)

	for _, image := range images.Images {
		// Unfortunately a request to list tags has to be done for each image.
		tags, err := r.regClient.ListTags(&registry.ListTagsRequest{
			ImageID: image.ID,
			OrderBy: registry.ListTagsRequestOrderByCreatedAtDesc,
		}, scw.WithAllPages())
		if err != nil {
			return nil, fmt.Errorf("error listing tags %w", err)
		}

		if len(tags.Tags) <= numberVersionsToKeep {
			// not enough versions to delete, skipping
			continue
		}

		slog.Info("appending tags for image: " + image.Name)

		tagsToDelete = append(tagsToDelete, tags.Tags[numberVersionsToKeep:]...)
	}

	return tagsToDelete, nil
}

func (r *RegistryAPI) DeleteTags(tagsToDelete []*registry.Tag, dryRun bool) error {
	if dryRun {
		slog.Info("Dry run mode ENABLED")

		for k := range tagsToDelete {
			slog.Info("dry-run: deleting tag: " + tagsToDelete[k].Name + " id: " + tagsToDelete[k].ID)
		}
	} else {
		slog.Warn("Dry run DISABLED")

		for k := range tagsToDelete {
			// dont delete latest:
			if !strings.EqualFold(tagsToDelete[k].Name, "latest") {
				_, err := r.regClient.DeleteTag(&registry.DeleteTagRequest{TagID: tagsToDelete[k].ID})
				if err != nil {
					return fmt.Errorf("error deleting registry tag %w", err)
				}
			}
		}
	}

	return nil
}
