package service

import (
	"context"
	"testing"

	"example.com/go-config-management-service/internal/core/domain"
	"example.com/go-config-management-service/internal/core/port"
)

func TestPutConfigurationSuccess(t *testing.T) {
	mockRepo := port.NewMockConfigurationRepository(t)

	configurationService := NewConfigurationService(mockRepo)

	mockRepo.On("PutConfiguration", context.Background(), &domain.Config{Name: "test-config", Value: map[string]interface{}{"name": "John", "age": 25}, Version: 1}).Return(&domain.Config{Name: "test-config", Version: 1}, nil)
	//	mockRepo.On("PutConfiguration", context.Background(), &domain.Config{Name: "error-config", Value: map[string]interface{}{"name": "John", "age": 25}, Version: 1}).Return(nil, domain.ErrDataNotFound)

	config, err := configurationService.PutConfiguration(context.Background(), &domain.Config{Name: "test-config", Value: map[string]interface{}{"name": "John", "age": 25}, Version: 1})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if config.Name != "test-config" || config.Version != 1 {
		t.Fatalf("expected config with name 'test-config' and version 1, got %v", config)
	}

	/*
		t.Run("Success", func(t *testing.T) {
			config, err := configurationService.PutConfiguration(context.Background(), &domain.Config{Name: "test-config", Value: map[string]interface{}{"name": "John", "age": 25}, Version: 1})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if config.Name != "test-config" || config.Version != 1 {
				t.Fatalf("expected config with name 'test-config' and version 1, got %v", config)
			}
		})
	*/

	/*
		t.Run("Invalid Schema", func(t *testing.T) {
			config, err := configurationService.PutConfiguration(context.Background(), &domain.Config{Name: "test-config", Value: map[string]interface{}{"name": "John"}, Version: 1})
			if config != nil || err == nil {
				t.Fatalf("expected error, got config: %v, error: %v", config, err)
			}

			if err != domain.ErrInvalidSchema {
				t.Fatalf("expected error %v, got %v", domain.ErrDataNotFound, err)
			}
		})

		t.Run("Error", func(t *testing.T) {
			config, err := configurationService.PutConfiguration(context.Background(), &domain.Config{Name: "error-config", Value: map[string]interface{}{"name": "John"}, Version: 1})
			if config != nil || err == nil {
				t.Fatalf("expected error, got config: %v, error: %v", config, err)
			}

			if err != domain.ErrDataNotFound {
				t.Fatalf("expected error %v, got %v", domain.ErrDataNotFound, err)
			}
		})
	*/
}

func TestPutConfigurationInvalidSchema(t *testing.T) {
	mockRepo := port.NewMockConfigurationRepository(t)

	configurationService := NewConfigurationService(mockRepo)

	config, err := configurationService.PutConfiguration(context.Background(), &domain.Config{Name: "test-config", Value: map[string]interface{}{"name": "John"}, Version: 1})
	if config != nil || err == nil {
		t.Fatalf("expected error, got config: %v, error: %v", config, err)
	}

	if err != domain.ErrInvalidSchema {
		t.Fatalf("expected error %v, got %v", domain.ErrDataNotFound, err)
	}
}

func TestPutConfigurationError(t *testing.T) {
	mockRepo := port.NewMockConfigurationRepository(t)

	configurationService := NewConfigurationService(mockRepo)

	//mockRepo.On("PutConfiguration", context.Background(), &domain.Config{Name: "test-config", Value: map[string]interface{}{"name": "John", "age": 25}, Version: 1}).Return(&domain.Config{Name: "test-config", Version: 1}, nil)
	mockRepo.On("PutConfiguration", context.Background(), &domain.Config{Name: "error-config", Value: map[string]interface{}{"name": "John", "age": 25}, Version: 1}).Return(nil, domain.ErrDataNotFound)

	config, err := configurationService.PutConfiguration(context.Background(), &domain.Config{Name: "error-config", Value: map[string]interface{}{"name": "John", "age": 25}, Version: 1})
	if config != nil || err == nil {
		t.Fatalf("expected error, got config: %v, error: %v", config, err)
	}

	if err != domain.ErrDataNotFound {
		t.Fatalf("expected error %v, got %v", domain.ErrDataNotFound, err)
	}
}

func TestGetConfiguration(t *testing.T) {
	mockRepo := port.NewMockConfigurationRepository(t)

	configurationService := NewConfigurationService(mockRepo)

	mockRepo.On("GetConfiguration", context.Background(), "test-config").Return(&domain.Config{Name: "test-config", Version: 1}, nil)
	mockRepo.On("GetConfiguration", context.Background(), "non-existent-config").Return(nil, domain.ErrDataNotFound)

	t.Run("Success", func(t *testing.T) {
		config, err := configurationService.GetConfiguration(context.Background(), "test-config")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if config.Name != "test-config" || config.Version != 1 {
			t.Fatalf("expected config with name 'test-config' and version 1, got %v", config)
		}
	})

	t.Run("Error", func(t *testing.T) {
		config, err := configurationService.GetConfiguration(context.Background(), "non-existent-config")
		if config != nil || err == nil {
			t.Fatalf("expected error, got config: %v, error: %v", config, err)
		}

		if err != domain.ErrDataNotFound {
			t.Fatalf("expected error %v, got %v", domain.ErrDataNotFound, err)
		}
	})
}

func TestListConfigurations(t *testing.T) {
	mockRepo := port.NewMockConfigurationRepository(t)

	configurationService := NewConfigurationService(mockRepo)

	mockRepo.On("ListConfigurations", context.Background(), uint64(0), uint64(10)).Return([]*domain.Config{
		{Name: "config1", Version: 1},
		{Name: "config2", Version: 2},
	}, nil)

	t.Run("Success", func(t *testing.T) {
		configs, err := configurationService.ListConfigurations(context.Background(), 0, 10)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(configs) != 2 || configs[0].Name != "config1" || configs[1].Name != "config2" {
			t.Fatalf("expected two configs with names 'config1' and 'config2', got %v", configs)
		}
	})
}

func TestListConfigurationVersions(t *testing.T) {
	mockRepo := port.NewMockConfigurationRepository(t)

	configurationService := NewConfigurationService(mockRepo)

	mockRepo.On("ListConfigurationVersions", context.Background(), "test-config", uint64(0), uint64(10)).Return([]*domain.Config{
		{Name: "test-config", Version: 1},
		{Name: "test-config", Version: 2},
	}, nil)

	t.Run("Success", func(t *testing.T) {
		configs, err := configurationService.ListConfigurationVersions(context.Background(), "test-config", 0, 10)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(configs) != 2 || configs[0].Version != 1 || configs[1].Version != 2 {
			t.Fatalf("expected two versions of 'test-config', got %v", configs)
		}
	})
}

func TestGetConfigurationVersion(t *testing.T) {
	mockRepo := port.NewMockConfigurationRepository(t)

	configurationService := NewConfigurationService(mockRepo)

	mockRepo.On("GetConfigurationVersion", context.Background(), "test-config", 1).Return(&domain.Config{Name: "test-config", Version: 1}, nil)
	mockRepo.On("GetConfigurationVersion", context.Background(), "test-config", 999).Return(nil, domain.ErrDataNotFound)

	t.Run("Success", func(t *testing.T) {
		config, err := configurationService.GetConfigurationVersion(context.Background(), "test-config", 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if config.Name != "test-config" || config.Version != 1 {
			t.Fatalf("expected config with name 'test-config' and version 1, got %v", config)
		}
	})

	t.Run("Error", func(t *testing.T) {
		config, err := configurationService.GetConfigurationVersion(context.Background(), "test-config", 999)
		if config != nil || err == nil {
			t.Fatalf("expected error, got config: %v, error: %v", config, err)
		}

		if err != domain.ErrDataNotFound {
			t.Fatalf("expected error %v, got %v", domain.ErrDataNotFound, err)
		}
	})
}

func TestRollbackConfigurationVersion(t *testing.T) {
	mockRepo := port.NewMockConfigurationRepository(t)

	configurationService := NewConfigurationService(mockRepo)

	mockRepo.On("RollbackConfigurationVersion", context.Background(), "test-config", 1).Return(&domain.Config{Name: "test-config", Version: 1}, nil)
	mockRepo.On("RollbackConfigurationVersion", context.Background(), "test-config", 999).Return(nil, domain.ErrDataNotFound)

	t.Run("Success", func(t *testing.T) {
		config, err := configurationService.RollbackConfigurationVersion(context.Background(), "test-config", 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if config.Name != "test-config" || config.Version != 1 {
			t.Fatalf("expected config with name 'test-config' and version 1, got %v", config)
		}
	})

	t.Run("Error", func(t *testing.T) {
		config, err := configurationService.RollbackConfigurationVersion(context.Background(), "test-config", 999)
		if config != nil || err == nil {
			t.Fatalf("expected error, got config: %v, error: %v", config, err)
		}

		if err != domain.ErrDataNotFound {
			t.Fatalf("expected error %v, got %v", domain.ErrDataNotFound, err)
		}
	})
}
