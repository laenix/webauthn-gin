package model

import "gorm.io/gorm"

type Challenge struct {
	gorm.Model
	Username  string
	Challenge string
}
