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
	Name string `uri:"name" binding:"required" example:"person_config"`
}

type putConfigurationRequestJson struct {
	Type  string                 `json:"type" binding:"required" example:"person"`
	Value map[string]interface{} `json:"value" swaggertype:"object,string" binding:"required"`
}

// PutConfiguration godoc
//
//	@Summary		Create a new configuration or replace an existing one
//	@Description	Create a new configuration with the specified name and value, or replace an existing
//	@Tags			Configurations
//	@Accept			json
//	@Produce		json
//	@Param			name					path		string						true	"Configuration name"	example:"person_config"
//	@Param			createCategoryRequest	body		putConfigurationRequestJson	true	"Create or Replace Configuration request"
//	@Success		200						{object}	configurationResponse		"Configuration created"
//	@Failure		400						{object}	errorResponse				"Validation error"
//	@Failure		401						{object}	errorResponse				"Unauthorized error"
//	@Failure		403						{object}	errorResponse				"Forbidden error"
//	@Failure		404						{object}	errorResponse				"Data not found error"
//	@Failure		409						{object}	errorResponse				"Data conflict error"
//	@Failure		500						{object}	errorResponse				"Internal server error"
//	@Router			/cms/configs/{name} [put]
//	@Security		BearerAuth
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
		Type:  reqJson.Type,
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

// GetConfiguration godoc
//
//	@Summary		Retrieve the latest version of a configuration
//	@Description	Retrieve the latest version of a configuration by its name
//	@Tags			Configurations
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string					true	"Configuration name"	example:"person_config"
//	@Success		200		{object}	configurationResponse	"Configuration found"
//	@Failure		400		{object}	errorResponse			"Validation error"
//	@Failure		401		{object}	errorResponse			"Unauthorized error"
//	@Failure		403		{object}	errorResponse			"Forbidden error"
//	@Failure		404		{object}	errorResponse			"Data not found error"
//	@Failure		409		{object}	errorResponse			"Data conflict error"
//	@Failure		500		{object}	errorResponse			"Internal server error"
//	@Router			/cms/configs/{name} [get]
//	@Security		BearerAuth
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
	Skip  uint64 `form:"skip" binding:"min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"min=1,max=100" example:"5"`
}

// ListConfigurations godoc
//
//	@Summary		Retrieve configuration list
//	@Description	Retrieve a list configurations with pagination support.
//	@Tags			Configurations
//	@Accept			json
//	@Produce		json
//	@Param			skip	query		int						false	"Starting offset"	example:"0"
//	@Param			limit	query		int						true	"Page size"			example:"5"
//	@Success		200		{object}	configurationResponse	"Configuration found"
//	@Failure		400		{object}	errorResponse			"Validation error"
//	@Failure		401		{object}	errorResponse			"Unauthorized error"
//	@Failure		403		{object}	errorResponse			"Forbidden error"
//	@Failure		404		{object}	errorResponse			"Data not found error"
//	@Failure		409		{object}	errorResponse			"Data conflict error"
//	@Failure		500		{object}	errorResponse			"Internal server error"
//	@Router			/cms/configs [get]
//	@Security		BearerAuth
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
		configsList = append(configsList, newConfigResponse(config))
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

// GetConfigurationVersion godoc
//
//	@Summary		Retrieve a particular version of a configuration
//	@Description	Retrieve a particular version of a configuration by its name and version number
//	@Tags			Configurations
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string					true	"Configuration name"	example:"person_config"
//	@Param			version	path		int						true	"Version Number"	example:"1"
//	@Success		200		{object}	configurationResponse	"Configuration found"
//	@Failure		400		{object}	errorResponse			"Validation error"
//	@Failure		401		{object}	errorResponse			"Unauthorized error"
//	@Failure		403		{object}	errorResponse			"Forbidden error"
//	@Failure		404		{object}	errorResponse			"Data not found error"
//	@Failure		409		{object}	errorResponse			"Data conflict error"
//	@Failure		500		{object}	errorResponse			"Internal server error"
//	@Router			/cms/configs/{name}/versions/{version} [get]
//	@Security		BearerAuth
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

type listConfigurationVersionsRequestUri struct {
	Name string `uri:"name" binding:"required" example:"app_config"`
}

type listConfigurationVersionsRequestForm struct {
	Skip  uint64 `form:"skip" binding:"min=0" example:"0"`
	Limit uint64 `form:"limit" binding:"min=5,max=100" example:"5"`
}

// ListConfigurationVersions godoc
//
//	@Summary		Retrieve a historical version list of a configuration
//	@Description	Retrieve a historical version list of a configuration by its name
//	@Tags			Configurations
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string					true	"Configuration name"	example:"person_config"
//	@Param			skip	query		int						false	"Starting offset"		example:"0"
//	@Param			limit	query		int						true	"Page size"				example:"5"
//
//	@Success		200		{object}	configurationResponse	"Configuration found"
//	@Failure		400		{object}	errorResponse			"Validation error"
//	@Failure		401		{object}	errorResponse			"Unauthorized error"
//	@Failure		403		{object}	errorResponse			"Forbidden error"
//	@Failure		404		{object}	errorResponse			"Data not found error"
//	@Failure		409		{object}	errorResponse			"Data conflict error"
//	@Failure		500		{object}	errorResponse			"Internal server error"
//	@Router			/cms/configs/{name}/versions [get]
//	@Security		BearerAuth
func (ch *ConfigurationHandler) ListConfigurationVersions(ctx *gin.Context) {
	var reqUri listConfigurationVersionsRequestUri
	if err := ctx.ShouldBindUri(&reqUri); err != nil {
		validationError(ctx, err)
		return
	}

	var reqForm listConfigurationVersionsRequestForm
	if err := ctx.ShouldBindQuery(&reqForm); err != nil {
		validationError(ctx, err)
		return
	}

	configs, err := ch.svc.ListConfigurationVersions(ctx, reqUri.Name, reqForm.Skip, reqForm.Limit)
	if err != nil {
		handleError(ctx, err)
		return
	}

	var configsList []configurationResponse
	for _, config := range configs {
		configsList = append(configsList, newConfigResponse(config))
	}

	total := uint64(len(configsList))
	meta := newMeta(total, reqForm.Limit, reqForm.Skip)
	rsp := toMap(meta, configsList, "configs")

	handleSuccess(ctx, rsp)
}

type rollbackConfigurationVersionRequest struct {
	Name    string `uri:"name" binding:"required" example:"app_config"`
	Version int    `uri:"version" binding:"required" example:"1"`
}

// RollbackConfigurationVersion godoc
//
//	@Summary		Rollback a configuration to a previous version
//	@Description	Rollback a configuration to a previous version
//	@Tags			Configurations
//	@Accept			json
//	@Produce		json
//	@Param			name	path		string					true	"Configuration name"	example:"person_config"
//	@Param			version	path		int						true	"Version Number"	example:"1"
//	@Success		200		{object}	configurationResponse	"Configuration rolled back"
//	@Failure		400		{object}	errorResponse			"Validation error"
//	@Failure		401		{object}	errorResponse			"Unauthorized error"
//	@Failure		403		{object}	errorResponse			"Forbidden error"
//	@Failure		404		{object}	errorResponse			"Data not found error"
//	@Failure		409		{object}	errorResponse			"Data conflict error"
//	@Failure		500		{object}	errorResponse			"Internal server error"
//	@Router			/cms/configs/{name}/versions/{version}/rollback [post]
//	@Security		BearerAuth
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
