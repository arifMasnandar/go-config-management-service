package memory

import (
	"context"
	"time"

	"github.com/arifMasnandar/go-config-management-service/internal/core/domain"
)

type ConfigurationRepository struct {
	configurations map[string][]*domain.Config
}

func NewConfigurationRepository() *ConfigurationRepository {
	return &ConfigurationRepository{
		configurations: make(map[string][]*domain.Config),
	}
}
func (r *ConfigurationRepository) PutConfiguration(ctx context.Context, config *domain.Config) (*domain.Config, error) {
	// Looking for a config whose Name value matches the parameter.
	versions, ok := r.configurations[config.Name]

	if ok {
		// If found, update the config and return it.
		config.Version = versions[len(versions)-1].Version + 1 // Increment the version for the updated config
		config.CreatedAt = time.Now()                          // Set the creation timestamp

		r.configurations[config.Name] = append(r.configurations[config.Name], config)
		return config, nil
	}

	// If not found, create a new config.
	config.Version = 1            // Set the version to 1 for a new config
	config.CreatedAt = time.Now() // Set the creation timestamp

	r.configurations[config.Name] = []*domain.Config{config}

	return config, nil
}

func (r *ConfigurationRepository) GetConfiguration(ctx context.Context, name string) (*domain.Config, error) {
	// Looking for a config whose Name value matches the parameter.
	versions, ok := r.configurations[name]

	if ok && len(versions) > 0 {
		return versions[len(versions)-1], nil // Return the latest version of the config
	}

	return nil, domain.ErrDataNotFound
}

func (r *ConfigurationRepository) ListConfigurations(ctx context.Context, skip, limit uint64) ([]*domain.Config, error) {
	var configs []*domain.Config

	for _, versions := range r.configurations {
		configs = append(configs, versions[len(versions)-1]) // Get the latest version of each config
	}

	if skip >= uint64(len(configs)) {
		return nil, nil // No configs to return
	}

	end := skip + limit
	if end > uint64(len(configs)) {
		end = uint64(len(configs))
	}

	return configs[skip:end], nil
}
func (r *ConfigurationRepository) ListConfigurationVersions(ctx context.Context, name string, skip, limit uint64) ([]*domain.Config, error) {
	versions, ok := r.configurations[name]

	if !ok {
		return nil, domain.ErrDataNotFound
	}

	if skip >= uint64(len(versions)) {
		return nil, nil // No versions to return
	}

	end := skip + limit
	if end > uint64(len(versions)) {
		end = uint64(len(versions))
	}

	return versions[skip:end], nil
}
func (r *ConfigurationRepository) GetConfigurationVersion(ctx context.Context, name string, version int) (*domain.Config, error) {
	// Looking for a config whose Name value matches the parameter.
	versions, ok := r.configurations[name]

	if ok {
		for _, v := range versions {
			if v.Version == version {
				return v, nil // Return the specific version
			}
		}
	}

	return nil, domain.ErrDataNotFound
}
func (r *ConfigurationRepository) RollbackConfigurationVersion(ctx context.Context, name string, version int) (*domain.Config, error) {
	// Looking for a config whose Name value matches the parameter.
	versions, ok := r.configurations[name]

	if !ok || len(versions) == 0 {
		return nil, domain.ErrDataNotFound
	}

	for _, v := range versions {
		if v.Version == version {
			newConfigVersion := *v                                           // Create a copy of the found version
			newConfigVersion.RollbackedVersion = version                     // Set the version to the rolled back version
			newConfigVersion.Version = versions[len(versions)-1].Version + 1 // Increment the version for the new config

			newConfigVersion.CreatedAt = time.Now() // Set the creation timestamp
			r.configurations[name] = append(r.configurations[name], &newConfigVersion)

			return &newConfigVersion, nil // Return the rolled back version
		}
	}

	return nil, domain.ErrDataNotFound
}
