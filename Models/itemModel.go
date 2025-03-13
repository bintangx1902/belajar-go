package Models

type Items struct {
	ID          uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	Item        string  `json:"item" gorm:"column:item;type:varchar(50);not null"`
	Description string  `json:"description" gorm:"column:description;type:text;"`
	Quantity    int     `json:"quantity" gorm:"column:quantity;not null"`
	Price       float64 `json:"price" gorm:"column:price;not null"`
	Discount    float64 `json:"discount" gorm:"column:discount;"`
}

type Invoice struct {
	ID     uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID uint    `json:"user_id" gorm:"column:user_id;not null"`
	Users  Users   `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items  []Items `json:"items" gorm:"many2many:invoice_items;"`
}

type InvoiceItems struct {
	InvoiceID uint `gorm:"primaryKey"`
	ItemsID   uint `gorm:"primaryKey"`
}

type PrintInvoiceItems struct {
	No              string
	Item            string
	Description     string
	Quantity        string
	Price           string
	DiscountedPrice string
	Total           string
}
