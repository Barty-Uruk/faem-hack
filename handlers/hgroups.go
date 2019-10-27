package handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/hack/logs"
	"gitlab.com/faemproject/backend/hack/models"
)

// GetUser godoc
func GetGroupsList(c echo.Context) error {
	uuid := c.Param("uuid")
	if uuid == "" {
		msg := emptyUUIDMsg
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "getting group list",
			"reason": msg,
		}).Error()

		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code: 400,
			Msg:  emptyUUIDMsg,
		})
	}
	groups := models.GetAllGroups()

	return c.JSON(http.StatusOK, groups)
}

// NewProjectStatus godoc
func AddUsersToGroup(c echo.Context) error {
	var data struct {
		UsersID   []int `json:"users_id"`
		ProjectID int   `json:"project_id"`
	}

	err := c.Bind(&data)
	if err != nil {
		msg := bindingDataErrorMsg
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "AddUsersToGroup",
			"reason": msg,
		}).Error(err)

		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code: 400,
			Msg:  "ошибка бинда",
		})
	}

	group := models.AddUserToGroup(data.UsersID, data.ProjectID)
	if err != nil {
		msg := "add error"
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "AddUsersToGroup",
			"reason": msg,
		}).Error(err)
		return c.JSON(http.StatusNotFound, ResponseStruct{
			Code: 400,
			Msg:  msg,
		})
	}
	return c.JSON(http.StatusOK, group)
}
