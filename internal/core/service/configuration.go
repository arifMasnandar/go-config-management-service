package service

import (
	"context"

	"github.com/arifMasnandar/go-config-management-service/internal/core/domain"
	"github.com/arifMasnandar/go-config-management-service/internal/core/port"
	"github.com/kaptinlin/jsonschema"
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
	repo    port.ConfigurationRepository
	schemas map[string]*jsonschema.Schema
}

func NewConfigurationService(repo port.ConfigurationRepository) ConfigurationServicer {
	schemas := make(map[string]*jsonschema.Schema)

	// Compile schema
	compiler := jsonschema.NewCompiler()
	schema, err := compiler.Compile([]byte(`{
			    "type": "object",
			    "properties": {
			        "name": {"type": "string", "minLength": 1},
			        "age": {"type": "integer", "minimum": 0}
			    },
			    "required": ["name","age"]
			}`))
	if err != nil {
		panic(err) // Handle schema compilation error
	}

	schemas["person"] = schema

	return &configurationService{
		repo,
		schemas,
	}
}

func (s *configurationService) PutConfiguration(ctx context.Context, config *domain.Config) (*domain.Config, error) {

	schema, ok := s.schemas[config.Type]

	if !ok {
		return nil, domain.ErrInvalidSchema // Schema not found for the config type
	}

	result := schema.ValidateMap(config.Value)
	if !result.IsValid() {
		/*
			for field, err := range result.Errors {
				fmt.Printf("- %s: %s\n", field, err.Message)
			}
		*/
		return nil, domain.ErrInvalidSchema
	}

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
