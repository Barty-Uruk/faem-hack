package db

import (
	"fmt"

	"gitlab.com/faemproject/backend/hack/logs"
	"gitlab.com/faemproject/backend/hack/models"

	pg "github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/sirupsen/logrus"
)

type dbLogger struct{}

var (
	modelsList []interface{}
)

func init() {
	modelsList = append(modelsList, (*models.Projects)(nil))
	modelsList = append(modelsList, (*models.Group)(nil))
	modelsList = append(modelsList, (*models.UsersCRM)(nil))
	modelsList = append(modelsList, (*models.ProjectStatuses)(nil))
	// modelsList = append(modelsList, (*model.RegionCRM)(nil))
	// modelsList = append(modelsList, (*model.DistrictCRM)(nil))
	// modelsList = append(modelsList, (*model.ServiceCRM)(nil))
	// modelsList = append(modelsList, (*model.DriverCRM)(nil))
	// modelsList = append(modelsList, (*model.ClientCRM)(nil))
	// modelsList = append(modelsList, (*model.FeatureCRM)(nil))
	// modelsList = append(modelsList, (*orders.OrderCRM)(nil))
	// modelsList = append(modelsList, (*model.DriverStatesCRM)(nil))

	// modelsList = append(modelsList, (*Drivers)(nil))
}

func (d dbLogger) BeforeQuery(q *pg.QueryEvent) {}

func (d dbLogger) AfterQuery(q *pg.QueryEvent) {

	logs.Eloger.Trace(q.FormattedQuery())
}

// Connect return DB connection
func Connect() (*pg.DB, error) {

	var conn *pg.DB

	addr := fmt.Sprintf("%s:%s", "78.110.156.74", "6001")

	conn = pg.Connect(&pg.Options{
		Addr:     addr,
		User:     "barman",
		Password: "ba4man80",
		Database: "hack",
	})
	var n int

	conn.AddQueryHook(dbLogger{})
	_, err := conn.QueryOne(pg.Scan(&n), "SELECT 1")
	if err != nil {
		return conn, err
	}

	if err := createSchema(conn); err != nil {
		return conn, fmt.Errorf("Error creating DB schemas. %v", err)
	}
	return conn, nil
}

// CloseDbConnection closing connection for defer in main
func CloseDbConnection(db *pg.DB) {
	db.Close()
}

func createSchema(db *pg.DB) error {
	logrus.Info("Creatind tables if not exist...")
	for _, m := range modelsList {
		err := db.CreateTable(m, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
