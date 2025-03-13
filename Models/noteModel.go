package Models

type Notes struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	NoteTitle string `gorm:"column:note_title;type:varchar(50);not null" json:"note_title" validate:"required,max=50"`
	NoteMsg   string `gorm:"column:note;type:text;not null" json:"note" bson:"note"`
	UserID    uint   `gorm:"column:user_id;not null" json:"user_id" bson:"user_id"`
	Users     Users  `gorm:"foreignKey:UserID" json:"user"`
}
