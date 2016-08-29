package models

import (
	"github.com/jinzhu/gorm"
)

type Namespace struct {
	gorm.Model
	Name   string
	Images []Image
}