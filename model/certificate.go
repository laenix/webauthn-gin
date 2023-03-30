package model

import "gorm.io/gorm"

type Certificate struct {
	gorm.Model
	UserID     uint
	Credention []byte
}
