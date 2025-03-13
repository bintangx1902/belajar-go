package Models

import (
	"time"
)

type Users struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"column:username;type:varchar(100);unique;not null" json:"username" bson:"username"`
	Password  string    `gorm:"column:password;type:varchar(100);not null" json:"password" bson:"password"`
	LastLogin time.Time ``
	Roles     string    `gorm:"column:roles;type:varchar(20);not null" json:"roles" bson:"roles"`
	//Notes     []Notes   `gorm:"foreignKey:id;AssociationForeignKey:user_id" json:"notes,omitempty"`
}
