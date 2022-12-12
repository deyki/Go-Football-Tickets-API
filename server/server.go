package server

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Go-FootballTickets/deyki/v2/controller"
	"github.com/Go-FootballTickets/deyki/v2/database"
	"github.com/Go-FootballTickets/deyki/v2/service"
)

func AppRun() {
	database.LoadEnvVariables()
	database.ConnectDB()

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter --> R for RunServer or A for Create Admin: ")
	scanner.Scan()
	action := scanner.Text()

	switch action {
	case "R":
		controller.GinRouter()
	case "A":
		service.CreateAdmin()
	default:
		fmt.Println("You must enter R or A..")
	}
}