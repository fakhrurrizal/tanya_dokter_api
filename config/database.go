package config

import (
	"fmt"
	"log"
	"tanya_dokter_app/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Database() *gorm.DB {

	host := LoadConfig().DatabaseHost
	user := LoadConfig().DatabaseUsername
	password := LoadConfig().DatabasePassword
	name := LoadConfig().DatabaseName
	port := LoadConfig().DatabasePort

	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, password, name, port)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	if LoadConfig().EnableDatabaseAutomigration {
		err = DB.AutoMigrate(
			&models.GlobalUser{},
			&models.GlobalRole{},
			&models.GlobalUser{},
			&models.NotificationForEmail{},
			&models.GlobalCategorySpecialist{},
			&models.ResetPasswordRequest{},
			&models.GlobalFile{},
			&models.GlobalDataDrugs{},
			&models.GlobalMessages{},
		)
		if err != nil {
			log.Fatalf("Auto migration failed: %v", err)
		} else {
			fmt.Println("Auto migration success ...")

		}
	}

	fmt.Println("Connected to Database:", name)

	return DB

}
