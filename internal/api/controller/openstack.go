package controller

import (
	"forklift/internal/domain"
	"forklift/internal/domain/entity"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetProjects godoc
// @Produce  json
// @Success 200 {string} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Tags Openstack
// @Router /api/openstack/projectList [get]
func GetProjects(e *echo.Echo, openstackRepository domain.OpenstackRepository) {
	e.GET("/api/openstack/projectList", func(c echo.Context) error {

		var (
			result []entity.Project
			err    error
		)

		if result, err = openstackRepository.GetProjects(); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, result)
	})
}

// GetFlavors godoc
// @Produce  json
// @Success 200 {string} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Tags Openstack
// @Router /api/openstack/flavorList [get]
func GetFlavors(e *echo.Echo, openstackRepository domain.OpenstackRepository) {
	e.GET("/api/openstack/flavorList", func(c echo.Context) error {

		var (
			result []entity.Flavor
			err    error
		)

		if result, err = openstackRepository.GetFlavors(); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, result)
	})
}

// GetNetworks godoc
// @Produce  json
// @Success 200 {string} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Tags Openstack
// @Param project query string true "project"
// @Router /api/openstack/networkList [get]
func GetNetworks(e *echo.Echo, openstackRepository domain.OpenstackRepository) {
	e.GET("/api/openstack/networkList", func(c echo.Context) error {

		var (
			result      []entity.Network
			err         error
			projectName string
		)

		projectName = c.QueryParam("project")

		if result, err = openstackRepository.GetNetworks(projectName); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, result)
	})
}

// GetKeys godoc
// @Produce  json
// @Success 200 {string} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Tags Openstack
// @Param project query string true "project"
// @Router /api/openstack/keyList [get]
func GetKeys(e *echo.Echo, openstackRepository domain.OpenstackRepository) {
	e.GET("/api/openstack/keyList", func(c echo.Context) error {

		var (
			result      []entity.Key
			err         error
			projectName string
		)

		projectName = c.QueryParam("project")

		if result, err = openstackRepository.GetKeys(projectName); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, result)
	})
}

// GetSecurityGroups godoc
// @Produce  json
// @Success 200 {string} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Tags Openstack
// @Param project query string true "project"
// @Router /api/openstack/securityGroupList [get]
func GetSecurityGroups(e *echo.Echo, openstackRepository domain.OpenstackRepository) {
	e.GET("/api/openstack/securityGroupList", func(c echo.Context) error {

		var (
			result      []entity.SecurityGroup
			err         error
			projectName string
		)

		projectName = c.QueryParam("project")

		if result, err = openstackRepository.GetSecurityGroups(projectName); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, result)
	})
}
