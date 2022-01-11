package controller

import (
	"encoding/json"
	"forklift/internal/api/model"
	"forklift/internal/domain"
	"forklift/internal/domain/entity"
	"github.com/labstack/echo/v4"
	"net/http"
)

// AddUser godoc
// @Accept  json
// @Produce  json
// @Success 200 {object} object
// @Failure 400 {string} string
// @Failure 401 {string} string
// @Failure 500 {string} string
// @Tags User
// @Router /api/user [post]
func AddUser(e *echo.Echo, userRepository domain.UserRepository) {
	e.POST("/api/user", func(c echo.Context) error {
		auth := c.QueryParams().Get("user")
		response := model.AuthResponse{}
		err := json.Unmarshal([]byte(auth), &response)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		user := entity.NewUser(response.AuthUser.FullName, response.AuthUser.UserName, response.AuthUser.UserEmail, response.AuthUser.UserId)

		if err = userRepository.Upsert(user); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}

		return c.JSON(http.StatusOK, "OK")
	})
}
