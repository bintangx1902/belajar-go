package Configs

import (
	"TestApp/Models"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "123456"
	dbname   = "db_note"
)

func ConnectDB() error {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbname, port)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix: "public.",
		//},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	//Ensure proper migration
	//err = DB.Migrator().DropTable(&Models.InvoiceItems{})
	//if err != nil {
	//	return fmt.Errorf("failed to drop table: %v", err)
	//}

	err = DB.AutoMigrate(
		&Models.Users{},
		&Models.Notes{},
		&Models.Items{},
		&Models.Invoice{},
		&Models.InvoiceItems{},
	)
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}
	//DB.LogMode(true)

	fmt.Println("Database connected and migrated successfully")
	return nil
}
