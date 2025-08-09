package http

import (
	"example.com/go-config-management-service/internal/core/domain"
	"example.com/go-config-management-service/internal/core/service"
	"github.com/gin-gonic/gin"
)

// ConfigurationHandler represents the HTTP handler for configuration-related requests
type ConfigurationHandler struct {
	svc service.ConfigurationServicer
}

// NewConfigurationHandler creates a new ConfigurationHandler instance
func NewConfigurationHandler(svc service.ConfigurationServicer) *ConfigurationHandler {
	return &ConfigurationHandler{
		svc,
	}
}

type putConfigurationRequestUri struct {
	Name string `uri:"name" binding:"required" example:"app_config"`
	//Value string `json:"value" binding:"required" example:"config_value"`
}

type putConfigurationRequestJson struct {
	//Name  string `uri:"name" binding:"required" example:"app_config"`
	Value string `json:"value" binding:"required" example:"config_value"`
}

func (ch *ConfigurationHandler) PutConfiguration(ctx *gin.Context) {
	var reqUri putConfigurationRequestUri

	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		validationError(ctx, err)
		return
	}

	var reqJson putConfigurationRequestJson

	if err := ctx.ShouldBindJSON(&reqJson); err != nil {
		validationError(ctx, err)
		return
	}

	config := &domain.Config{
		Name:  reqUri.Name,
		Value: reqJson.Value,
	}

	createdConfig, err := ch.svc.PutConfiguration(ctx, config)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newConfigResponse(createdConfig)

	handleSuccess(ctx, rsp)
}

type getConfigurationRequest struct {
	Name string `uri:"name" binding:"required" example:"app_config"`
}

func (ch *ConfigurationHandler) GetConfiguration(ctx *gin.Context) {
	var req getConfigurationRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	config, err := ch.svc.GetConfiguration(ctx, req.Name)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newConfigResponse(config)

	handleSuccess(ctx, rsp)
}

type listConfigurationsRequest struct {
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

func (ch *ConfigurationHandler) ListConfigurations(ctx *gin.Context) {
	var req listConfigurationsRequest
	var configsList []configurationResponse

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	configs, err := ch.svc.ListConfigurations(ctx, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, config := range configs {
		configsList = append(configsList, newConfigResponse(&config))
	}

	total := uint64(len(configsList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, configsList, "configs")

	handleSuccess(ctx, rsp)
}

type getConfigurationVersionRequest struct {
	Name    string `uri:"name" binding:"required" example:"app_config"`
	Version int    `uri:"version" binding:"required" example:"1"`
}

func (ch *ConfigurationHandler) GetConfigurationVersion(ctx *gin.Context) {
	var req getConfigurationVersionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}
	config, err := ch.svc.GetConfigurationVersion(ctx, req.Name, req.Version)
	if err != nil {
		handleError(ctx, err)
		return
	}
	rsp := newConfigResponse(config)
	handleSuccess(ctx, rsp)
}

type listConfigurationVersionsRequest struct {
	Name  string `uri:"name" binding:"required" example:"app_config"`
	Skip  uint64 `form:"skip" binding:"required,min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"required,min=5" example:"5"`
}

func (ch *ConfigurationHandler) ListConfigurationVersions(ctx *gin.Context) {
	var req listConfigurationVersionsRequest
	var configsList []configurationResponse

	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		validationError(ctx, err)
		return
	}

	configs, err := ch.svc.ListConfigurationVersions(ctx, req.Name, req.Skip, req.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	for _, config := range configs {
		configsList = append(configsList, newConfigResponse(&config))
	}

	total := uint64(len(configsList))
	meta := newMeta(total, req.Limit, req.Skip)
	rsp := toMap(meta, configsList, "configs")

	handleSuccess(ctx, rsp)
}

type rollbackConfigurationVersionRequest struct {
	Name    string `uri:"name" binding:"required" example:"app_config"`
	Version int    `uri:"version" binding:"required" example:"1"`
}

func (ch *ConfigurationHandler) RollbackConfigurationVersion(ctx *gin.Context) {
	var req rollbackConfigurationVersionRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		validationError(ctx, err)
		return
	}

	config, err := ch.svc.RollbackConfigurationVersion(ctx, req.Name, req.Version)
	if err != nil {
		handleError(ctx, err)
		return
	}

	rsp := newConfigResponse(config)
	handleSuccess(ctx, rsp)
}
