package database

import (
	"fmt"
	"os"
	

	"github.com/Go-FootballTickets/deyki/v2/util"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


type Ticket struct {
	gorm.Model
	Host        	string		`json:"host"`
	Guest       	string		`json:"guest"`
	StadiumName 	string		`json:"stadiumName"`
	Time			string		`json:"time"`
	Price			string		`json:"price"`
	InStock			bool		`json:"inStock"`
}


type Admin struct {
	gorm.Model
	Username	string	`json:"username"`
	Password	string	`json:"password"`
}


type DBConfig struct {
	Host 		string
	User 		string
	Password 	string
	Name 		string
	Port 		string
}


func LoadEnvVariables() *util.ErrorMessage {
	errorMessage := godotenv.Load(".env")
	if errorMessage != nil {
		return util.ErrorMessage{}.ErrorLoadingEnvFile()
	}

	return nil
}


func ConnectDB() (*gorm.DB, *util.ErrorMessage) {
	dbConfig := &DBConfig{
		Host: 		os.Getenv("DB_HOST"),
		User: 		os.Getenv("DB_USER"),
		Password: 	os.Getenv("DB_PASSWORD"),
		Name: 		os.Getenv("DB_NAME"),
		Port: 		os.Getenv("DB_PORT"),
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.Port)
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, util.ErrorMessage{}.FailedToOpenDB()
	}

	db.AutoMigrate(&Ticket{}, &Admin{})

	return db, nil
}



