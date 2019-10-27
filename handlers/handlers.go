package handlers

import (
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-pg/pg/urlvalues"
	"github.com/labstack/echo"
)

const (
	emptyUUIDMsg        = "empty uuid"
	bindingDataErrorMsg = "error binding data"
)

// UserIDFromJWT return user ID
func userIDFromJWT(c echo.Context) int {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64)
	return int(userID)
}
func userUUIDFromJWT(c echo.Context) string {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userUUID := claims["user_uuid"].(string)
	return userUUID
}
func getPager(c echo.Context) (urlvalues.Pager, error) {
	var (
		pager urlvalues.Pager
	)
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return pager, fmt.Errorf("Error parsing pager limit,%s", err)
	}

	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return pager, fmt.Errorf("Error parse orders page,%s", err)
	}
	pager.Offset = (page - 1) * limit
	pager.Limit = limit
	return pager, nil
}

func getTimeParams(c echo.Context) (time.Time, time.Time, error) {
	var (
		minData, maxData time.Time
	)
	mindateParam := c.QueryParam("mindate")
	if mindateParam != "" {
		mindateUnix, err := strconv.ParseInt(mindateParam, 10, 64)
		if err != nil {
			return minData, maxData, fmt.Errorf("error mindate parsing,%s", err)
		}
		minData = time.Unix(mindateUnix, 0)
	}
	maxData = time.Now()
	maxdateParam := c.QueryParam("maxdate")
	if maxdateParam != "" {
		maxdateUnix, err := strconv.ParseInt(maxdateParam, 10, 64)
		if err != nil {
			return minData, maxData, fmt.Errorf("error maxdate parsing,%s", err)
		}
		maxData = time.Unix(maxdateUnix, 0)
	}
	return minData, maxData, nil
}

type ResponseStruct struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}
