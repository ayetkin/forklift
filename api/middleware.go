package api

import (
	"bytes"
	"encoding/json"
	"forklift/internal/api/model"
	"forklift/pkg/config"
	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
)

func Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		var configuration config.Configuration
		err := viper.Unmarshal(&configuration)

		if strings.Contains(c.Request().RequestURI, "healthz") || strings.Contains(c.Request().RequestURI, "swagger") {
			return next(c)
		}

		auth := c.Request().Header.Get("authorization")

		var jsonStr = []byte(`{"token":"` + auth + `"}`)

		req, err := http.NewRequest("POST", configuration.Auth.ValidationUrl, bytes.NewBuffer(jsonStr))
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Error(err)
		}
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		authRes := model.AuthResponse{}
		_ = json.Unmarshal(body, &authRes)

		if !authRes.Success {
			return c.JSON(http.StatusUnauthorized, "Unauthorized")
		}
		out, err := json.Marshal(authRes.AuthUser)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		user := string(out)

		c.Request().Header.Set("user", user)
		return next(c)
	}
}
