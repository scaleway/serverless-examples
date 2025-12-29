package main

import (
	"fmt"
	"os"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Environment variable constants
const (
	envOrgID          = "SCW_DEFAULT_ORGANIZATION_ID"
	envAccessKey      = "SCW_ACCESS_KEY"
	envSecretKey      = "SCW_SECRET_KEY"
	envProjectID      = "SCW_DEFAULT_PROJECT_ID"
	envZone           = "SCW_ZONE"
	envBucket         = "SCW_BUCKET_NAME"
	envBucketEndpoint = "SCW_BUCKET_ENDPOINT"
)

type Config struct {
	OrgID          string
	AccessKey      string
	SecretKey      string
	ProjectID      string
	Zone           scw.Zone
	BucketName     string
	BucketEndpoint string
}

func LoadConfig() (*Config, error) {
	// Mandatory variables
	vars := map[string]*string{
		envAccessKey:      new(string),
		envSecretKey:      new(string),
		envProjectID:      new(string),
		envZone:           new(string),
		envBucket:         new(string),
		envBucketEndpoint: new(string),
	}

	// Optional variables
	orgID := os.Getenv(envOrgID)

	for envKey, valPtr := range vars {
		val := os.Getenv(envKey)
		if val == "" {
			return nil, fmt.Errorf("missing environment variable %s", envKey)
		}
		*valPtr = val
	}

	return &Config{
		OrgID:          orgID,
		AccessKey:      *vars[envAccessKey],
		SecretKey:      *vars[envSecretKey],
		ProjectID:      *vars[envProjectID],
		Zone:           scw.Zone(*vars[envZone]),
		BucketName:     *vars[envBucket],
		BucketEndpoint: *vars[envBucketEndpoint],
	}, nil
}
