package logs

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	// Eloger logrus instance
	Eloger *logrus.Logger
)

// Логика ЛОГОВ

// 	"event": "set driver state" //  Само событие
// 	"driverUUID": de00a4e1-7c63-47cd-b426-23302f37cb86
// 	"offerUUID": 0738ac67-c18f-4040-8801-d5de68a5d1e3
//  "orderUUID": e5f16e62-fbdd-4444-b2d1-8f3de3ef55a5
// 	"reason": "driver request" // Причина срабатываения
// 	"value": newState // Значение

// ConnectElastic - передаем адрес и название индекса
// "crm_logs"
// func ConnectElastic(url, eindex string) (*logrus.Logger, error) {
// 	httpURL := fmt.Sprintf("http://%s", url)
// 	client, err := elastic.NewClient(elastic.SetURL(httpURL), elastic.SetSniff(false))
// 	if err != nil {
// 		formError := fmt.Sprintf("Connection error. %s", err)
// 		return Eloger, fmt.Errorf(formError)
// 	}
// 	hook, err := elogrus.NewElasticHook(client, url, logrus.DebugLevel, eindex)
// 	if err != nil {
// 		formError := fmt.Sprintf("Error logrus hooking Elasticsearch. %s", err)
// 		return Eloger, fmt.Errorf(formError)
// 	}
// 	Eloger.Hooks.Add(hook)

// 	return Eloger, nil
// }
func InitElogger(eindex string) *logrus.Logger {
	return Eloger
}

func InitEloggerWithLevel(level string) (*logrus.Logger, error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse a log level")
	}
	Eloger.SetLevel(lvl)
	return Eloger, nil
}

func init() {
	Eloger = logrus.New()
	// default logging level
	Eloger.SetLevel(logrus.TraceLevel)
}

// OutputUserEvent godoc godoc
func OutputUserEvent(event, msg string, userID int) {
	Eloger.WithFields(logrus.Fields{
		"event":  event,
		"userID": userID,
	}).Info(msg)
}

// OutputEvent godoc godoc
func OutputEvent(event, msg string, userID int) {
	fmt.Printf("\n\n EVENT must be cutted - %s, Message - %s\n\n\n", event, msg)
	// Eloger.WithFields(logrus.Fields{
	// 	"event":  event,
	// 	"userID": userID,
	// }).Info(msg)
}

// OutputDebugEvent godoc
func OutputDebugEvent(event, msg string, userID int) {
	Eloger.WithFields(logrus.Fields{
		"event":  event,
		"userID": userID,
	}).Debug(msg)
}

// OutputDriverEvent godoc
func OutputDriverEvent(event, msg string, driver string) {
	// Eloger.WithFields(logrus.Fields{
	// 	"event":    event,
	// 	"driverID": driver,
	// }).Info(msg)
}

// OutputClientEvent godoc
func OutputClientEvent(event, msg string, client string) {
	Eloger.WithFields(logrus.Fields{
		"event":    event,
		"clientID": client,
	}).Info(msg)
}

// OutputOrderEvent godoc
func OutputOrderEvent(event, msg string, driver string) {
	Eloger.WithFields(logrus.Fields{
		"event":   event,
		"orderID": driver,
	}).Info(msg)
}

// OutputInfo godoc
func OutputInfo(msg string) {
	Eloger.Info(msg)
}

// OutputError godoc
func OutputError(msg string) {
	// Eloger.Error(msg)
}
