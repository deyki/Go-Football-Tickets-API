package service

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/Go-FootballTickets/deyki/v2/database"
	"github.com/Go-FootballTickets/deyki/v2/util"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type ResponseMessage struct {
	Message string
}


type TicketRequest struct {
	Host        	string		`json:"host"`
	Guest       	string		`json:"guest"`
	StadiumName 	string		`json:"stadiumName"`
	Time		string		`json:"time"`
	Price		string		`json:"price"`	
}


func CreateAdmin() (*ResponseMessage, *util.ErrorMessage) {
	db, err := database.ConnectDB()
	if err != nil {
		return nil, util.ErrorMessage{}.FailedToOpenDB()
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter admin username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Enter admin password: ")
	scanner.Scan()
	password := scanner.Text()

	hash, errorMessage := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errorMessage != nil {
		return nil, util.ErrorMessage{}.FailedToCreateHashFromPassword()
	}

	admin := &database.Admin{Username: username, Password: string(hash)}

	db.Create(admin)

	return &ResponseMessage{"Admin created!"}, nil
}


func Login(loginRequest *database.Admin) (*ResponseMessage, *util.ErrorMessage) {
	db, err := database.ConnectDB()
	if err != nil {
		return nil, util.ErrorMessage{}.FailedToOpenDB()
	}

	var admin database.Admin

	errorMessage := db.First(&admin, "username = ?", loginRequest.Username).Error
	if errorMessage != nil {
		return nil, util.ErrorMessage{}.AdminNotFound()
	}

	compareHashAndPass := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(loginRequest.Password))
	if compareHashAndPass != nil {
		return nil, util.ErrorMessage{}.AdminNotFound()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": admin.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, errorMessage := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if errorMessage != nil {
		return nil, util.ErrorMessage{}.FailedToCreateJWToken()
	}

	return &ResponseMessage{tokenString}, nil
}


func NewTicket(ticketRequest *TicketRequest) (*ResponseMessage, *util.ErrorMessage) {
	db, err := database.ConnectDB()
	if err != nil {
		return nil, util.ErrorMessage{}.FailedToOpenDB()
	}

	ticket := &database.Ticket{
		Host: 			ticketRequest.Host,
		Guest: 			ticketRequest.Guest,
		StadiumName: 		ticketRequest.StadiumName,
		Time: 			ticketRequest.Time,
		Price: 			ticketRequest.Price,
		InStock: 		true,
	}

	db.Create(ticket)

	return &ResponseMessage{"Ticket created!"}, nil
}


func GetTickets() (*[]database.Ticket, *util.ErrorMessage) {
	db, err := database.ConnectDB()
	if err != nil {
		return nil, util.ErrorMessage{}.FailedToOpenDB()
	}

	var tickets []database.Ticket

	db.Find(&tickets)

	return &tickets, nil
}


func UpdateTicketAvailability(ticketID int) (*ResponseMessage, *util.ErrorMessage) {
	db, err := database.ConnectDB()
	if err != nil {
		return nil, util.ErrorMessage{}.FailedToOpenDB()
	}

	var ticket database.Ticket

	errorMessage := db.First(&ticket, ticketID).Error
	if errorMessage != nil {
		return nil, util.ErrorMessage{}.TicketNotFound()
	}

	switch ticket.InStock {
	case true:
		ticket.InStock = false
	case false:
		ticket.InStock = true
	}

	db.Save(&ticket)

	return &ResponseMessage{"Ticket updated!"}, nil
}


func GetAvailableTickets() (*[]database.Ticket, *util.ErrorMessage) {
	db, err := database.ConnectDB()
	if err != nil {
		return nil, util.ErrorMessage{}.FailedToOpenDB()
	}

	var tickets []database.Ticket

	db.Where("InStock = ?", true).Find(&tickets)

	return &tickets, nil
}
