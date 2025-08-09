package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// config represents data about a record config.
type config struct {
	Name              string    `json:"name"`
	Value             string    `json:"value"`
	Version           int       `json:"version"`
	RollbackedVersion int       `json:"rollbacked_version,omitempty"` // Optional field for copied version
	CreatedAt         time.Time `json:"created_at,omitempty"`         // Optional field for creation timestamp
}

var configVersions = make(map[string][]config) // Map to hold configs by name with their versions

// getConfigs responds with the list of all configs as JSON.
func getConfigs(c *gin.Context) {

	configs := make([]config, len(configVersions))

	idx := 0
	for _, versions := range configVersions {
		configs[idx] = versions[len(versions)-1] // Get the latest version of each config
		idx++
	}
	c.IndentedJSON(http.StatusOK, configs)
}

// putConfigByName inserts or raplaces a config from JSON received in the request body.
func putConfigByName(c *gin.Context) {
	Name := c.Param("name")

	var newConfig config

	// Call BindJSON to bind the received JSON to
	// newConfig.
	if err := c.BindJSON(&newConfig); err != nil {
		return
	}

	newConfig.Name = Name // Ensure the Name field is set from the URL parameter

	// Looking for a config whose Name value matches the parameter.
	versions, ok := configVersions[Name]

	if ok {
		// If found, update the config and return it.

		newConfig.Version = versions[len(versions)-1].Version + 1 // Increment the version for the updated config
		newConfig.CreatedAt = time.Now()                          // Set the creation timestamp

		c.IndentedJSON(http.StatusOK, newConfig)

		configVersions[Name] = append(configVersions[Name], newConfig)
		// Update the versions for this config.
		return
	}

	// If not found, create a new config.

	newConfig.Version = 1            // Set the version to 1 for a new config
	newConfig.CreatedAt = time.Now() // Set the creation timestamp

	// Initialize versions for the new config if not already present.
	configVersions[Name] = []config{newConfig}

	// Add the new config to the slice.
	//configs = append(configs, newConfig)
	c.IndentedJSON(http.StatusOK, newConfig)
}

// getConfigByName locates the config whose Name value matches the Name
// parameter sent by the client, then returns that config as a response.
func getConfigByName(c *gin.Context) {
	Name := c.Param("name")

	// Loop over the list of configs, looking for
	// a config whose Name value matches the parameter.
	versions, ok := configVersions[Name]

	if ok {
		c.IndentedJSON(http.StatusOK, versions[len(versions)-1]) // Return the latest version
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "config not found"})
}

func getConfigVersionsByName(c *gin.Context) {
	Name := c.Param("name")

	versions := configVersions[Name]

	if versions != nil {
		c.IndentedJSON(http.StatusOK, versions)
		return
	}

	// If no versions found, return a 404 status.
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "config not found"})
}

// getConfigVersionsByNameAndVersion locates the config whose Name and version values matches the Name and
// the Version parameters sent by the client, then returns that config as a response.
func getConfigVersionsByNameAndVersion(c *gin.Context) {
	Name := c.Param("name")
	Version := c.Param("version")

	// Convert Version to an integer
	versionInt, err := strconv.Atoi(Version)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid version"})
		return
	}

	// Loop over the list of configs, looking for
	// a config whose Name value matches the parameter.
	versions, ok := configVersions[Name]

	if ok {
		for _, v := range versions {
			if v.Version == versionInt {
				c.IndentedJSON(http.StatusOK, v) // Return the specific version
				return
			}
		}
		// If no specific version found, return the latest version.
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "version not found"})

		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "config not found"})
}

func rollbackConfigVersion(c *gin.Context) {
	Name := c.Param("name")
	Version := c.Param("version")

	// Convert Version to an integer
	versionInt, err := strconv.Atoi(Version)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid version"})
		return
	}

	// Looking for a config whose Name value matches the parameter.
	versions, ok := configVersions[Name]

	if ok {
		for _, v := range versions {
			if v.Version == versionInt {
				// Rollback logic: Set the latest version to the specified version

				var newConfigVersion config
				newConfigVersion.Name = Name
				newConfigVersion.Value = v.Value                                 // Use the value from the rolled back version
				newConfigVersion.RollbackedVersion = versionInt                  // Set the version to the rolled back version
				newConfigVersion.Version = versions[len(versions)-1].Version + 1 // Increment the version for the new config

				configVersions[Name] = append(configVersions[Name], newConfigVersion)
				c.IndentedJSON(http.StatusOK, newConfigVersion) // Return the rolled back version
				return
			}
		}
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "version not found"})
}

func main() {
	router := gin.Default()
	router.GET("/configs", getConfigs)
	router.GET("/configs/:name", getConfigByName)
	router.PUT("/configs/:name", putConfigByName)
	router.GET("/configs/:name/versions", getConfigVersionsByName)
	router.GET("/configs/:name/versions/:version", getConfigVersionsByNameAndVersion)
	router.POST("configs/:name/versions/:version/rollback", rollbackConfigVersion)

	router.Run("localhost:8080")
}
