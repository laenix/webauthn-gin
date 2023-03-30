package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name         string `gorm:"unique"`
	DisplayName  string `gorm:"type:varchar(100)"`
	Certificates []byte `gorm:""`
}
