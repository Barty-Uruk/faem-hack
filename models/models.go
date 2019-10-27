package models

import (
	"fmt"
	"time"

	"github.com/go-pg/pg"
)

var (
	db *pg.DB
)

const (
	NeedStatus     = "need_check"
	RejectedStatus = "rejected"
	ApprovedStatus = "approved"
)

type ormStruct struct {
	Deleted   bool      `json:"-" sql:"default:false"`
	ID        int       `json:"id" sql:",pk"`
	CreatedAt time.Time `sql:"default:now()" json:"-" description:"Дата создания"`
	// UpdatedAt time.Time `json:"updated_at" `
	UpdatedAt time.Time `json:"-"`
}

// ConnectDB initialize connection to package var
func ConnectDB(conn *pg.DB) {
	db = conn
}

// CheckExistsUUID return error if UUID in mdl model not found
func CheckExistsUUID(mdl interface{}, uuid string) error {
	exs, _ := db.Model(mdl).Column("uuid").Where("uuid = ? AND deleted is not true", uuid).Exists()

	if !exs {
		return fmt.Errorf("Запись UUID=%v не найдена или удалена", uuid)
	}
	return nil
}

// CheckExistsUUIDWithError godoc
// TODO: перевести все с CheckExistsUUID на CheckExistsUUIDWithError
func CheckExistsUUIDWithError(mdl interface{}, uuid string) (bool, error) {
	exs, err := db.Model(mdl).Column("uuid").Where("uuid = ? AND deleted is not true", uuid).Exists()
	return exs, err
}

// UpdateByPK update model by PK
func UpdateByPK(mdl interface{}) error {
	_, err := db.Model(mdl).Where("uuid=?uuid").Returning("*").UpdateNotNull()
	if err != nil {
		return err
	}
	return nil
}

// CreateByStruct - created sended struct
func CreateByStruct(mdl interface{}) error {
	_, err := db.Model(mdl).Returning("*").Insert()
	if err != nil {
		return err
	}
	return nil
}

// GetByUUID returning object by interface and UUID
func GetByUUID(uuid string, mdl interface{}) error {
	err := db.Model(mdl).
		Where("uuid = ?", uuid).
		Select()
	if err != nil {
		return err
	}
	return nil
}

type SmallUser struct {
	Name   string `json:"name"`
	Role   string `json:"role"`
	Ava    string `json:"ava"`
	Status string `json:"status"`
}

type ProjectStatuses struct {
	ID        int    `json:"id" sql:",pk"`
	UserID    int    `json:"user_id"`
	ProjectID int    `json:"project_id"`
	Comment   string `json:"comment"`
	UserName  string `json:"user_name"`
	Status    string `json:"status"`
}

func GetAllPSByUserUUID(userUUID string) []ProjectStatuses {
	var result []ProjectStatuses
	err := db.Model(&result).
		Returning("*").
		Select()
	if err != nil {
		fmt.Printf("=======error GetAllPSByUserUUID,%s", err)
	}
	return result
}

type Projects struct {
	ID       int    `json:"id" sql:",pk"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	FullDesc string `json:"full_desc"`
	Xp       int    `json:"xp"`
}

type ProjectTree struct {
	UUID     string        `json:"uuid"`
	Name     string        `json:"name"`
	ImageURL string        `json:"image_url"`
	Mate     []ProjectTree `json:"mate"`
	Children []ProjectTree `json:"children"`
}

type Group struct {
	tableName struct{}   `sql:"groups"`
	Users     []UsersCRM `json:"users"`
	ProjectID int        `json:"project_id"`
	ID        int        `json:"id" sql:",pk"`
}

func GetProjectByUUID(id int) Projects {
	var proj Projects
	err := db.Model(&proj).Where("id = ?", id).
		Select()
	if err != nil {
		fmt.Printf("=======error GetProjectByUUID,%s", err)
	}
	return proj
}

func GetStatusesForTeacher() []ProjectStatuses {
	var statuses []ProjectStatuses
	err := db.Model(&statuses).
		Where("status = ?", NeedStatus).
		Order("id desc").
		Select()
	if err != nil {
		fmt.Printf("=======error GetStatusesForTeacher,%s", err)
	}
	return statuses
}
func GetStatusesByProjectID(projectid int) []ProjectStatuses {
	var statuses []ProjectStatuses
	err := db.Model(&statuses).
		Where("project_id = ?", projectid).
		Order("id desc").
		Select()
	if err != nil {
		fmt.Printf("=======error GetStatusesByProjectID,%s", err)
	}
	return statuses
}
func GetAllGroups() []Group {
	var res []Group
	err := db.Model(&res).Select()
	if err != nil {
		fmt.Print("==GetAllGroxups=======", err)
	}
	return res
}

func AddUserToGroup(usersid []int, projectid int) Group {
	var res Group
	var users UsersCRM
	if len(usersid) == 0 {
		return res
	}

	err := db.Model(&users).Where("id in (?)", pg.In(usersid)).
		Select()
	if err != nil {
		fmt.Print("==error select users by id array=======", err)
		return res
	}
	check, _ := db.Model(&res).Where("project_id = ?", projectid).Exists()
	if check {
		err = db.Model(&res).Where("project_id = ?", projectid).Select()
		if err != nil {
			fmt.Print("==AddUserToGroup=======", err)
		}
		res.Users = append(res.Users, users)
		_, err = db.Model(&res).Where("project_id = ?", projectid).
			Set("users = ?", res.Users).
			Returning("*").
			Update()
		if err != nil {
			fmt.Print("==AddUserToGroupUpdate=======", err)
		}
	} else {
		res.ProjectID = projectid
		res.Users = append(res.Users, users)
		db.Model(&res).Insert()
	}
	return res
}
func GetGroup(projectid int) Group {
	var res Group
	err := db.Model(&res).Where("project_id = ?", projectid).Select()
	if err != nil {
		fmt.Print("==GetGroup=======", err)
	}
	res.ProjectID = projectid
	return res
}
func (status ProjectStatuses) Save() error {
	var user UsersCRM
	err := db.Model(&user).Where("id = ?", status.UserID).
		Select()
	if err != nil {
		fmt.Printf("=======error getting user,%s", err)
	}

	status.UserName = user.Name
	_, err = db.Model(&status).Insert()
	return err
}
