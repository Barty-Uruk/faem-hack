package logs

import (
	"context"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func LoggerForContext(ctx context.Context) logrus.FieldLogger {
	if ctx == nil {
		return Eloger
	}
	rid, ok := ctx.Value(echo.HeaderXRequestID).(string)
	if !ok || rid == "" {
		return Eloger
	}
	return Eloger.WithField("request_id", rid)
}
