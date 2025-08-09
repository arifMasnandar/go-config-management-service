package service

import (
	"context"

	"example.com/go-config-management-service/internal/core/domain"
	"example.com/go-config-management-service/internal/core/port"
)

type ConfigurationServicer interface {
	PutConfiguration(ctx context.Context, config *domain.Config) (*domain.Config, error)
	GetConfiguration(ctx context.Context, name string) (*domain.Config, error)
	ListConfigurations(ctx context.Context, skip, limit uint64) ([]*domain.Config, error)
	ListConfigurationVersions(ctx context.Context, name string, skip, limit uint64) ([]*domain.Config, error)
	GetConfigurationVersion(ctx context.Context, name string, version int) (*domain.Config, error)
	RollbackConfigurationVersion(ctx context.Context, name string, version int) (*domain.Config, error)
}

type configurationService struct {
	repo port.ConfigurationRepository
}

func NewConfigurationService(repo port.ConfigurationRepository) ConfigurationServicer {
	return &configurationService{
		repo,
	}
}

func (s *configurationService) PutConfiguration(ctx context.Context, config *domain.Config) (*domain.Config, error) {
	return s.repo.PutConfiguration(ctx, config)
}

func (s *configurationService) GetConfiguration(ctx context.Context, name string) (*domain.Config, error) {
	return s.repo.GetConfiguration(ctx, name)
}

func (s *configurationService) ListConfigurations(ctx context.Context, skip, limit uint64) ([]*domain.Config, error) {
	return s.repo.ListConfigurations(ctx, skip, limit)
}

func (s *configurationService) ListConfigurationVersions(ctx context.Context, name string, skip, limit uint64) ([]*domain.Config, error) {
	return s.repo.ListConfigurationVersions(ctx, name, skip, limit)
}
func (s *configurationService) GetConfigurationVersion(ctx context.Context, name string, version int) (*domain.Config, error) {
	return s.repo.GetConfigurationVersion(ctx, name, version)
}
func (s *configurationService) RollbackConfigurationVersion(ctx context.Context, name string, version int) (*domain.Config, error) {
	return s.repo.RollbackConfigurationVersion(ctx, name, version)
}
