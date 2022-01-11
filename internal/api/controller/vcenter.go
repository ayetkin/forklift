package controller

import (
	"fmt"
	"forklift/internal/api/model"
	"forklift/internal/domain"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
)

// GetPoweredOffVms godoc
// @Accept  json
// @Produce  json
// @Param request body model.VmListRequest true "request"
// @Success 200 {object} object
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Tags vCenter
// @Router /api/vcenter/vmList [post]
func GetPoweredOffVms(e *echo.Echo, vcenterRepository domain.VcenterRepository) {
	e.POST("/api/vcenter/vmList", func(c echo.Context) error {

		ctx := context.Background()

		var (
			poweredOffVmList []string
			err              error
		)

		request := new(model.VmListRequest)

		if err = c.Bind(request); err != nil {
			log.Error(ctx, "request deserialize error.", err)
			return c.JSON(http.StatusBadRequest, fmt.Sprintf("request deserialize error: %s", err))
		}

		if poweredOffVmList, err = vcenterRepository.GetAllPoweredOffVms(request.Dc); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, poweredOffVmList)
	})
}

// GetDatacenters godoc
// @Summary Get all dc list
// @Description Get all dc list
// @Produce  json
// @Success 200 {object} object
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Tags vCenter
// @Router /api/vcenter/dcList [get]
func GetDatacenters(e *echo.Echo, vcenterRepository domain.VcenterRepository) {
	e.GET("/api/vcenter/dcList", func(c echo.Context) error {

		var (
			dcList []string
			err    error
		)

		if dcList, err = vcenterRepository.GetAllDcs(); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, dcList)
	})
}
