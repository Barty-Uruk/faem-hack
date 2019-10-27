package handlers

import (
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"gitlab.com/faemproject/backend/hack/logs"
	"gitlab.com/faemproject/backend/hack/models"
)

// CreateUser godoc
func CreateUser(c echo.Context) error {

	var usient models.UsersCRM
	err := c.Bind(&usient)
	if err != nil {
		msg := bindingDataErrorMsg
		logs.Eloger.WithFields(logrus.Fields{
			"event": "creating user",

			"reason": msg,
		}).Error(err)

		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code: 400,
			Msg:  "ошибка",
		})
	}

	nUser, err := usient.Create()
	if err != nil {
		msg := "creating error"
		logs.Eloger.WithFields(logrus.Fields{
			"event": "creating user",

			"reason": msg,
		}).Error(err)
		return c.JSON(http.StatusNotFound, ResponseStruct{
			Code: 400,
			Msg:  "ошибка",
		})
	}

	logs.Eloger.WithFields(logrus.Fields{
		"event":       "creating user",
		"newUserUUID": nUser.UUID,
	}).Info("new user created")

	return c.JSON(http.StatusOK, nUser)
}

// GetUser godoc
func GetUsers(c echo.Context) error {

	users, err := models.UsersList()
	if err != nil {
		msg := "getting error"
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "getting users",
			"reason": msg,
		}).Error(err)

		return c.JSON(http.StatusNotFound, ResponseStruct{
			Code: 400,
			Msg:  "ошибка",
		})
	}
	return c.JSON(http.StatusOK, users)
}

// GetUser godoc
func GetUser(c echo.Context) error {
	uuid := c.Param("uuid")
	if uuid == "" {
		msg := emptyUUIDMsg
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "getting user",
			"reason": msg,
		}).Error()

		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code: 400,
			Msg:  "ошибка бинда",
		})
	}
	var user models.UsersCRM
	err := models.GetByUUID(uuid, &user)
	if err != nil {
		msg := "getting error"
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "getting user",
			"reason": msg,
		}).Error(err)

		return c.JSON(http.StatusNotFound, ResponseStruct{
			Code: 400,
			Msg:  "ошибка бинда",
		})
	}
	return c.JSON(http.StatusOK, user)
}

// GetStatusesForTeacher godoc
func GetStatusesForTeacher(c echo.Context) error {
	proj := models.GetStatusesForTeacher()
	return c.JSON(http.StatusOK, proj)
}

// GetStatusesByProject godoc
func GetStatusesByProject(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if id == 0 {
		msg := emptyUUIDMsg
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "getting  statuses project",
			"reason": msg,
		}).Error()

		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code: 400,
			Msg:  msg,
		})
	}
	proj := models.GetStatusesByProjectID(id)
	if err != nil {
		msg := "getting error"
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "getting statuses project",
			"reason": msg,
		}).Error(err)

		return c.JSON(http.StatusNotFound, ResponseStruct{
			Code: 400,
			Msg:  msg,
		})
	}
	return c.JSON(http.StatusOK, proj)
}

// GetUser godoc
func GetProject(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if id == 0 {
		msg := emptyUUIDMsg
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "getting project",
			"reason": msg,
		}).Error()

		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code: 400,
			Msg:  msg,
		})
	}
	proj := models.GetProjectByUUID(id)
	if err != nil {
		msg := "getting error"
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "getting project",
			"reason": msg,
		}).Error(err)

		return c.JSON(http.StatusNotFound, ResponseStruct{
			Code: 400,
			Msg:  msg,
		})
	}
	group := models.GetGroup(id)

	return c.JSON(http.StatusOK, struct {
		Group   models.Group    `json:"group"`
		Project models.Projects `json:"project"`
	}{group, proj})
}
func userIsTeacher(c echo.Context) bool {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userrole := claims["role"].(string)
	return userrole == "teacher"
}

// NewProjectStatus godoc
func NewProjectStatus(c echo.Context) error {
	var status models.ProjectStatuses

	err := c.Bind(&status)
	if err != nil {
		msg := bindingDataErrorMsg
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "NewProjectStatust",
			"reason": msg,
		}).Error(err)

		return c.JSON(http.StatusBadRequest, ResponseStruct{
			Code: 400,
			Msg:  "ошибка бинда",
		})
	}
	if !userIsTeacher(c) {
		status.UserID = userIDFromJWT(c)
	}
	err = status.Save()
	if err != nil {
		msg := "Save error"
		logs.Eloger.WithFields(logrus.Fields{
			"event":  "NewProjectStatus",
			"reason": msg,
		}).Error(err)
		return c.JSON(http.StatusNotFound, ResponseStruct{
			Code: 400,
			Msg:  msg,
		})
	}
	return c.JSON(http.StatusOK, ResponseStruct{
		Code: http.StatusOK,
		Msg:  "OK",
	})
}
