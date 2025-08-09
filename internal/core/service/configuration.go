package service

import (
	"context"
	"time"

	"example.com/go-config-management-service/internal/core/domain"
)

type ConfigurationServicer interface {
	PutConfiguration(ctx context.Context, config *domain.Config) (*domain.Config, error)
	GetConfiguration(ctx context.Context, name string) (*domain.Config, error)
	ListConfigurations(ctx context.Context, skip, limit uint64) ([]domain.Config, error)
	ListConfigurationVersions(ctx context.Context, name string, skip, limit uint64) ([]domain.Config, error)
	GetConfigurationVersion(ctx context.Context, name string, version int) (*domain.Config, error)
	RollbackConfigurationVersion(ctx context.Context, name string, version int) (*domain.Config, error)
}

type configurationService struct {
}

func NewConfigurationService() ConfigurationServicer {
	return &configurationService{}
}

var configVersions = map[string][]domain.Config{
	"config1": {
		{
			Name:      "config1",
			Value:     "value1",
			Version:   1,
			CreatedAt: time.Now(),
		},
	},
}

func (s *configurationService) PutConfiguration(ctx context.Context, config *domain.Config) (*domain.Config, error) {
	// Looking for a config whose Name value matches the parameter.
	versions, ok := configVersions[config.Name]

	if ok {
		// If found, update the config and return it.

		config.Version = versions[len(versions)-1].Version + 1 // Increment the version for the updated config
		config.CreatedAt = time.Now()                          // Set the creation timestamp

		configVersions[config.Name] = append(configVersions[config.Name], *config)
		// Update the versions for this config.
		return config, nil
	}

	// If not found, create a new config.

	config.Version = 1            // Set the version to 1 for a new config
	config.CreatedAt = time.Now() // Set the creation timestamp

	// Initialize versions for the new config if not already present.
	configVersions[config.Name] = []domain.Config{*config}

	return config, nil
}

func (s *configurationService) GetConfiguration(ctx context.Context, name string) (*domain.Config, error) {
	// Loop over the list of configs, looking for
	// a config whose Name value matches the parameter.
	versions, ok := configVersions[name]

	if ok {
		return &versions[len(versions)-1], nil
	}

	return nil, domain.ErrDataNotFound
}

func (s *configurationService) ListConfigurations(ctx context.Context, skip, limit uint64) ([]domain.Config, error) {
	configs := make([]domain.Config, len(configVersions))

	idx := 0
	for _, versions := range configVersions {
		configs[idx] = versions[len(versions)-1] // Get the latest version of each config
		idx++
	}
	return configs, nil
}

func (s *configurationService) ListConfigurationVersions(ctx context.Context, name string, skip, limit uint64) ([]domain.Config, error) {
	versions := configVersions[name]

	if versions != nil {
		return versions, nil

	}

	// If no versions found, return a 404 status.
	return nil, domain.ErrDataNotFound
}
func (s *configurationService) GetConfigurationVersion(ctx context.Context, name string, version int) (*domain.Config, error) {
	// Loop over the list of configs, looking for
	// a config whose Name value matches the parameter.
	versions, ok := configVersions[name]

	if ok {
		for _, v := range versions {
			if v.Version == version {
				return &v, nil // Return the specific version
			}
		}

		return nil, domain.ErrDataNotFound
	}

	return nil, domain.ErrDataNotFound
}
func (s *configurationService) RollbackConfigurationVersion(ctx context.Context, name string, version int) (*domain.Config, error) {
	// Looking for a config whose Name value matches the parameter.
	versions, ok := configVersions[name]

	if ok {
		for _, v := range versions {
			if v.Version == version {
				// Rollback logic: Set the latest version to the specified version

				var newConfigVersion domain.Config
				newConfigVersion.Name = v.Name
				newConfigVersion.Value = v.Value                                 // Use the value from the rolled back version
				newConfigVersion.RollbackedVersion = version                     // Set the version to the rolled back version
				newConfigVersion.Version = versions[len(versions)-1].Version + 1 // Increment the version for the new config

				newConfigVersion.CreatedAt = time.Now() // Set the creation timestamp
				configVersions[name] = append(configVersions[name], newConfigVersion)

				return &newConfigVersion, nil // Return the rolled back version
			}
		}
	}

	return nil, domain.ErrDataNotFound

}
