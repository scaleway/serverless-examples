package main

import (
	"fmt"
	"log/slog"

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
}

// NewRegistryAPI creates a new RegistryAPI to manage the Scaleway Container Registry API.
// It initializes the RegistryAPI struct with the provided Scaleway SDK client and a project
// ID. If the projectID is empty, it will not be passed to the Scaleway SDK, allowing operations
// on all projects.
func NewRegistryAPI(client *scw.Client, projectID string) *RegistryAPI {
	return &RegistryAPI{
		regClient: registry.NewAPI(client),
		projectID: scw.StringPtr(projectID),
	}
}

func (r *RegistryAPI) DeleteEmptyNamespace(dryRun bool) error {
	namespaces, err := r.regClient.ListNamespaces(&registry.ListNamespacesRequest{ProjectID: r.projectID}, scw.WithAllPages())
	if err != nil {
		return fmt.Errorf("error listing registry namespaces %w", err)
	}
	slog.Info("DryRun ENABLED")

	for _, namespace := range namespaces.Namespaces {
		if namespace.Status == registry.NamespaceStatusReady && namespace.ImageCount == 0 {
			slog.Info("deleteing namespace", slog.String("name", namespace.Name), slog.String("id", namespace.ID))
			if !dryRun {
				_, err := r.regClient.DeleteNamespace(&registry.DeleteNamespaceRequest{
					NamespaceID: namespace.ID,
				})
				if err != nil {
					return fmt.Errorf("error deleting namesapce %w", err)
				}
			}
		}
	}

	return nil
}
