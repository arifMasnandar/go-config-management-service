package memory

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/arifMasnandar/go-config-management-service/internal/core/domain"
)

func TestEmptyConfigurations(t *testing.T) {
	repo := NewConfigurationRepository()

	if repo == nil || repo.configurations == nil || len(repo.configurations) != 0 {
		t.Error("Expected new repository to be initialized with empty configurations")
	}

	// List configurations
	configs, err := repo.ListConfigurations(context.Background(), 0, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(configs) != 0 {
		t.Errorf("Expected no configurations, got %d", len(configs))
	}

	config := &domain.Config{
		Name:  "test_config",
		Value: map[string]interface{}{"name": "John"},
	}

	// Get configuration
	got, err := repo.GetConfiguration(context.Background(), config.Name)

	if err == nil || got != nil {
		t.Errorf("Failed to get configuration: %v", err)
	}

	if err != domain.ErrDataNotFound {
		t.Errorf("Expected error %v, got %v", domain.ErrDataNotFound, err)
	}

	// Get historical versions
	versions, err := repo.ListConfigurationVersions(context.Background(), config.Name, 0, 10)
	if err == nil || versions != nil {
		t.Errorf("Failed to get configuration: %v", err)
	}

	if err != domain.ErrDataNotFound {
		t.Errorf("Expected error %v, got %v", domain.ErrDataNotFound, err)
	}

	// Get specific version
	_, err = repo.GetConfigurationVersion(context.Background(), config.Name, 1)
	if err == nil {
		t.Errorf("Expected error when getting non-existing version, got nil")
	}
	if err != domain.ErrDataNotFound {
		t.Errorf("Expected error %v, got %v", domain.ErrDataNotFound, err)
	}

	// Rollback configuration version
	_, err = repo.RollbackConfigurationVersion(context.Background(), config.Name, 1)
	if err == nil {
		t.Errorf("Expected error when rolling back non-existing version, got nil")
	}
	if err != domain.ErrDataNotFound {
		t.Errorf("Expected error %v, got %v", domain.ErrDataNotFound, err)
	}
}

func TestPutConfiguration(t *testing.T) {
	repo := NewConfigurationRepository()
	config := &domain.Config{
		Name:  "test_config",
		Value: map[string]interface{}{"name": "John"},
	}

	t1 := time.Now()

	// Put configuration
	createdConfig, err := repo.PutConfiguration(context.Background(), config)
	if err != nil {
		t.Fatalf("Failed to put configuration: %v", err)
	}
	t2 := time.Now()

	validateConfig(t, createdConfig, config.Name, config.Value, 1, t1, time.Now())

	if len(repo.configurations) != 1 {
		t.Errorf("Expected configurations to have 1 entry, got %d", len(repo.configurations))
	}

	if repo.configurations[config.Name] == nil || len(repo.configurations[config.Name]) != 1 {
		t.Errorf("Expected configurations for %s to have 1 entry, got %d", config.Name, len(repo.configurations[config.Name]))
	}

	// List configurations
	configs, err := repo.ListConfigurations(context.Background(), 0, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(configs) != 1 {
		t.Errorf("Expected 1 configurations, got %d", len(configs))
	}

	configs, err = repo.ListConfigurations(context.Background(), 1, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(configs) != 0 {
		t.Errorf("Expected 0 configurations, got %d", len(configs))
	}

	// Get configuration
	got, err := repo.GetConfiguration(context.Background(), config.Name)

	if err != nil {
		t.Errorf("Failed to get configuration: %v", err)
	}

	validateConfig(t, got, config.Name, config.Value, 1, t1, t2)

	// Get historical versions
	versions, err := repo.ListConfigurationVersions(context.Background(), config.Name, 0, 10)
	if err != nil {
		t.Errorf("Failed to get configuration: %v", err)
	}

	if len(versions) != 1 {
		t.Errorf("Expected 1 version, got %d", len(versions))
	}

	validateConfig(t, versions[0], config.Name, config.Value, 1, t1, t2)

	versions, err = repo.ListConfigurationVersions(context.Background(), config.Name, 1, 10)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(versions) != 0 {
		t.Errorf("Expected 0 configurations, got %d", len(configs))
	}

	// Get specific version
	ver, err := repo.GetConfigurationVersion(context.Background(), config.Name, 1)
	if err != nil {
		t.Errorf("Failed to get configuration version: %v", err)
	}

	validateConfig(t, ver, config.Name, config.Value, 1, t1, t2)

	_, err = repo.GetConfigurationVersion(context.Background(), config.Name, 2)
	if err == nil {
		t.Errorf("Expected error when getting non-existing version, got nil")
	}
	if err != domain.ErrDataNotFound {
		t.Errorf("Expected error %v, got %v", domain.ErrDataNotFound, err)
	}

	// Rollback configuration version
	t1 = time.Now()
	rolledBackConfig, err := repo.RollbackConfigurationVersion(context.Background(), config.Name, 1)
	if err != nil {
		t.Errorf("Failed to rollback configuration version: %v", err)
	}
	t2 = time.Now()

	validateConfig(t, rolledBackConfig, config.Name, config.Value, 2, t1, t2)

	if len(repo.configurations) != 1 {
		t.Errorf("Expected configurations to have 1 entry, got %d", len(repo.configurations))
	}
	if repo.configurations[config.Name] == nil || len(repo.configurations[config.Name]) != 2 {
		t.Errorf("Expected configurations for %s to have 2 entries, got %d", config.Name, len(repo.configurations[config.Name]))
	}

	ver, err = repo.RollbackConfigurationVersion(context.Background(), config.Name, 3)
	if ver != nil || err == nil {
		t.Errorf("Expected error when rolling back non-existing version, got nil")
	}
	if err != domain.ErrDataNotFound {
		t.Errorf("Expected error %v, got %v", domain.ErrDataNotFound, err)
	}

}

func TestReplaceConfiguration(t *testing.T) {
	repo := NewConfigurationRepository()
	config := &domain.Config{
		Name:  "test_config",
		Value: map[string]interface{}{"name": "John"},
	}

	// Put initial configuration
	_, err := repo.PutConfiguration(context.Background(), config)
	if err != nil {
		t.Fatalf("Failed to put initial configuration: %v", err)
	}

	t1 := time.Now()
	// Update configuration
	config.Value = map[string]interface{}{"name": "John II"}

	updatedConfig, err := repo.PutConfiguration(context.Background(), config)
	if err != nil {
		t.Fatalf("Failed to update configuration: %v", err)
	}

	t2 := time.Now()
	validateConfig(t, updatedConfig, config.Name, config.Value, 2, t1, t2)

	if len(repo.configurations) != 1 {
		t.Errorf("Expected configurations to have 1 entry, got %d", len(repo.configurations))
	}

	if repo.configurations[config.Name] == nil || len(repo.configurations[config.Name]) != 2 {
		t.Errorf("Expected configurations for %s to have 2 entries, got %d", config.Name, len(repo.configurations[config.Name]))
	}
}

func validateConfig(t *testing.T, config *domain.Config, expectedName string, expectedValue map[string]interface{}, expectedVersion int, expectedCreatedAtAfter time.Time, expectedCreatedAtBefore time.Time) {
	if config.Name != expectedName {
		t.Errorf("Expected Name %s, got %s", expectedName, config.Name)
	}
	if reflect.DeepEqual(config.Value, expectedValue) == false {
		t.Errorf("Expected Value %s, got %s", expectedValue, config.Value)
	}
	if config.Version != expectedVersion {
		t.Errorf("Expected Version %d, got %d", expectedVersion, config.Version)
	}
	if config.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set, but it is zero")
	}
	if config.CreatedAt.Before(expectedCreatedAtAfter) || config.CreatedAt.After(expectedCreatedAtBefore) {
		t.Errorf("Expected CreatedAt to be within the range of %v and %v, got %v", expectedCreatedAtAfter, expectedCreatedAtBefore, config.CreatedAt)
	}
}
