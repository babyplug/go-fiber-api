package model

import (
	"gorm.io/gorm"
)

type UserStatus string

const (
	UserStatusNormal  UserStatus = "NORMAL"
	UserStatusLocked  UserStatus = "LOCKED"
	UserStatusBlocked UserStatus = "BLOCKED"
)

type User struct {
	gorm.Model

	Username  string     `json:"username"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Status    UserStatus `json:"status"`
}

func (u User) ToDTO() UserDTO {
	return UserDTO(u)
}

func (u User) FromDTO(dto UserDTO) any {
	return User(dto)
}

type UserDTO struct {
	gorm.Model

	Username  string     `json:"username"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
	Status    UserStatus `json:"status"`
}
