package Models

type User struct {
	ID       uint64 `json:"id" gorm:"primary_key;autoIncrement;type:serial"`
	Name     string `json:"name" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`
}
