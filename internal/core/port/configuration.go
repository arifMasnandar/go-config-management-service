package port

import (
	"context"

	"github.com/arifMasnandar/go-config-management-service/internal/core/domain"
)

type ConfigurationRepository interface {
	PutConfiguration(ctx context.Context, config *domain.Config) (*domain.Config, error)
	GetConfiguration(ctx context.Context, name string) (*domain.Config, error)
	ListConfigurations(ctx context.Context, skip, limit uint64) ([]*domain.Config, error)
	ListConfigurationVersions(ctx context.Context, name string, skip, limit uint64) ([]*domain.Config, error)
	GetConfigurationVersion(ctx context.Context, name string, version int) (*domain.Config, error)
	RollbackConfigurationVersion(ctx context.Context, name string, version int) (*domain.Config, error)
}
