package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/hack/logs"
	"gitlab.com/faemproject/backend/hack/models"
)

// LoginHandler godoc
func LoginHandler(c echo.Context) error {

	var loginData models.LoginRequest
	if err := c.Bind(&loginData); err != nil {

		logs.Eloger.WithFields(logrus.Fields{
			"event":  "user login [crm]",
			"reason": "Error binding data",
		}).Error(err)

		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code: 400,
			Msg:  "ошибка бинда",
		})
	}
	loginResp, err := models.AuthenticateUser(loginData)
	if err != nil {
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "user login [crm]",
			"reason": "Authentification error",
		}).Error(err)

		return c.JSON(http.StatusUnauthorized, ResponseStruct{
			Code: 400,
			Msg:  "ошибка аутентификации",
		})
	}

	logs.Eloger.WithFields(logrus.Fields{
		"event": "user login [crm]",
		"value": loginResp.UserUUID,
	}).Info("User logged in")

	return c.JSON(http.StatusOK, loginResp)
}
