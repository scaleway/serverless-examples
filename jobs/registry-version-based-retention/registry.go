package main

import (
	"fmt"
	"log/slog"
	"strings"

	registry "github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// RegistryAPI represents a Scaleway Container Registry accessor and extends
// capabilities to clean images. It allows you to manage container images and tags
// across one or more projects, and provides options to delete image tags safely or
// with caution.
type RegistryAPI struct {
	// regClient Scaleway Container Registry accessor.
	regClient *registry.API

	// projectID specifies a project to be scoped for operations. If this field
	// is nil, operations will be performed on all available projects.
	projectID *string

	// disableProtection if set to true, it will allow deletion of images that
	// might be in use by Serverless Jobs, Functions, or Containers. This should
	// be used with caution.
	disableProtection bool
}

// NewRegistryAPI creates a new RegistryAPI to manage the Scaleway Container Registry API.
// It initializes the RegistryAPI struct with the provided Scaleway SDK client and a project
// ID. If the projectID is empty, it will not be passed to the Scaleway SDK, allowing operations
// on all projects.
func NewRegistryAPI(client *scw.Client, projectID string) *RegistryAPI {
	return &RegistryAPI{
		regClient:         registry.NewAPI(client),
		projectID:         scw.StringPtr(projectID),
		disableProtection: false,
	}
}

// GetTagsAfterNVersions returns a list of image tags that should be deleted, based on the number of
// versions to keep. This function lists all container images and their associated tags, and determines
// which tags are beyond the specified count of versions to retain.
//
// The numberVersionsToKeep parameter specifies how many versions of each image should be preserved.
// If this value is less than or equal to 1, an error is returned to prevent accidental deletion of
// all image tags.
//
// The function returns a slice of pointers to registry.Tag structures representing the tags to be
// deleted, or an error if any issues arise during the process.
func (r *RegistryAPI) GetTagsAfterNVersions(numberVersionsToKeep int) ([]*registry.Tag, error) {
	images, err := r.regClient.ListImages(&registry.ListImagesRequest{ProjectID: r.projectID}, scw.WithAllPages())
	if err != nil {
		return nil, fmt.Errorf("error listing container images: %w", err)
	}

	if numberVersionsToKeep <= 1 {
		return nil, fmt.Errorf("number of versions to keep <= 1 is dangerous")
	}

	tagsToDelete := make([]*registry.Tag, 0)

	for _, image := range images.Images {
		tags, err := r.regClient.ListTags(&registry.ListTagsRequest{
			ImageID: image.ID,
			OrderBy: registry.ListTagsRequestOrderByCreatedAtDesc,
		}, scw.WithAllPages())
		if err != nil {
			return nil, fmt.Errorf("error listing tags for image %s: %w", image.Name, err)
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

// DeleteTags deletes the specified image tags. If the dryRun parameter is set to true,
// the function will log the tags that would be deleted without actually performing the
// deletion. If dryRun is false, the function will proceed to delete the tags.
//
// The tagsToDelete parameter is a slice of pointers to registry.Tag structures representing
// the tags to be deleted.
//
// The function logs informational messages about the operations being performed and returns
// an error if any issues arise during the process.
func (r *RegistryAPI) DeleteTags(tagsToDelete []*registry.Tag, dryRun bool) error {
	if dryRun {
		slog.Info("Dry run mode ENABLED")

		for _, tag := range tagsToDelete {
			slog.Info("dry-run: deleting tag:", slog.String("tag name", tag.Name), slog.String("tagID", tag.ID))
		}
	} else {
		slog.Warn("Dry run DISABLED")

		for _, tag := range tagsToDelete {
			if strings.EqualFold(tag.Name, "latest") {
				slog.Info("skipping deletion of latest tag", slog.String("tag name", tag.Name))

				continue
			}

			_, err := r.regClient.DeleteTag(&registry.DeleteTagRequest{TagID: tag.ID})
			if err != nil {
				return fmt.Errorf("error deleting registry tag %s (id %s): %w", tag.Name, tag.ID, err)
			}
		}
	}

	return nil
}
