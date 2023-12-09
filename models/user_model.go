package models

import (
	"task-5-pbi-btpns-muhammad-zachrie-kurniawan/helpers"
	"time"
)

type User struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"unique" json:"username"`
	Email     string    `gorm:"unique" json:"email"`
	Password  string    `json:"password"`
	Photo     []Photo   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (user *User) HashPassword(password string) error {
	hashed, err := helpers.HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	result, err := helpers.CheckPasswordHash(providedPassword, user.Password)
	if !result {
		return err
	}
	return nil
}
