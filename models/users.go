package models

import (
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

// UsersCRM - order user structure
type UsersCRM struct {
	tableName struct{} `sql:"crm_users"`
	ormStruct
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Login     string `json:"login"`
	Exp       int    `json:"exp"`
	GroupUUID string `json:"group_uuid"`
	Power     int    `json:"power"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

// Create - save structure to database
func (us *UsersCRM) Create() (UsersCRM, error) {
	nUser, err := us.validCreation()
	if err != nil {
		return nUser, fmt.Errorf("Data validation error. %v ", err)
	}

	uu, err := uuid.NewV4()
	if err != nil {
		return nUser, err
	}

	nUser.UUID = uu.String()
	_, err = db.Model(&nUser).Returning("*").Insert()
	if err != nil {
		return nUser, fmt.Errorf("DB error. %v ", err)
	}
	return nUser, nil
}

func (us UsersCRM) validCreation() (UsersCRM, error) {
	var newUser UsersCRM
	if us.Name == "" {
		return newUser, fmt.Errorf("Field name is required")
	}
	if us.Login == "" {
		return newUser, fmt.Errorf("Field login is required")
	}
	if us.Password == "" {
		return newUser, fmt.Errorf("Field password is required")
	}
	if us.Password == "" {
		return newUser, fmt.Errorf("Field role is required")
	}
	if us.Role != "student" && us.Role != "teacher" {
		return newUser, fmt.Errorf("Field role is invalid")
	}
	return us, nil
}

// UsersList - return all Users
func UsersList() ([]UsersCRM, error) {
	var users []UsersCRM
	err := db.Model(&users).
		Select()
	if err != nil {
		return users, err
	}
	return users, nil
}

// Update godoc
func (us *UsersCRM) Update(uuid string) (UsersCRM, error) {

	uUser, err := us.validUpdate(uuid)
	if err != nil {
		return uUser, fmt.Errorf("Data validation error. %v", err)
	}

	err = UpdateByPK(&uUser)
	if err != nil {
		return uUser, fmt.Errorf("DB error. %v", err)
	}
	return uUser, nil
}

//
func (us *UsersCRM) validUpdate(uuid string) (UsersCRM, error) {

	flag := false
	var user UsersCRM

	if us.Name != "" {
		user.Name = us.Name
		flag = true
	}

	if us.Role != "" {
		user.Role = us.Role
		flag = true
	}
	if us.Login != "" {
		user.Login = us.Login
		flag = true
	}
	if us.Password != "" {
		user.Password = us.Password
		flag = true
	}

	if !flag {
		return user, fmt.Errorf("Required fielus are empty")
	}
	// user.Users = us.Users
	user.UpdatedAt = time.Now()
	user.UUID = uuid

	err := CheckExistsUUID(&user, uuid)
	if err != nil {
		return user, err
	}
	user.Deleted = false
	return user, nil
}

// // CheckExists return nil if us.ID not exists or deleted
// // TODO: вынести в общие методы и сделать его универсальным
// func (us *UsersCRM) CheckExists(id int) error {
// 	exs, _ := db.Model(us).Column("id").Where("id = ? AND deleted is not true", id).Exists()
// 	if !exs {
// 		return fmt.Errorf("Not found ID=%v", id)
// 	}
// 	return nil
// }

// // BeforeUpdate godoc
// // TODO: вынести в общие методы и сделать его универсальным
// func (us *UsersCRM) BeforeUpdate(db orm.DB) error {
// 	us.UpdatedAt = time.Nus()
// 	return nil
// }

// SetDeleted godoc
func (us *UsersCRM) SetDeleted() error {
	// TODO: вынести в общие методы и сделать его универсальным
	var user UsersCRM
	err := CheckExistsUUID(&user, us.UUID)
	if err != nil {
		return err
	}
	us.Deleted = true
	_, err = db.Model(us).Where("uuid = ?uuid").UpdateNotNull()
	if err != nil {
		return err
	}
	return nil
}

// UsersByRegion godoc
func UsersByRegion(uuid string) ([]UsersCRM, error) {
	var users []UsersCRM
	err := db.Model(&users).
		Where("region_uuid = ? AND deleted is not true", uuid).
		Select()
	if err != nil {
		return users, err
	}
	return users, nil
}
