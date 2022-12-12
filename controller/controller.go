package controller

import (
	"net/http"
	"strconv"

	"github.com/Go-FootballTickets/deyki/v2/database"
	"github.com/Go-FootballTickets/deyki/v2/middleware"
	"github.com/Go-FootballTickets/deyki/v2/service"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var loginRequest database.Admin

	if c.BindJSON(&loginRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	signInResponse, errorMessage := service.Login(&loginRequest)
	if errorMessage != nil {
		c.IndentedJSON(errorMessage.HttpStatus, gin.H{"error": errorMessage})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", signInResponse.Message, 3600 * 24 * 30, "", "", false, true)
	c.IndentedJSON(http.StatusOK, gin.H{"Authenticated": true})
}


func NewTicket(c *gin.Context) {
	var ticketRequest service.TicketRequest

	if c.BindJSON(&ticketRequest) != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read body"})
		return
	}

	responseMessage, errorMessage := service.NewTicket(&ticketRequest)
	if errorMessage != nil {
		c.IndentedJSON(errorMessage.HttpStatus, gin.H{"error": errorMessage})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"Response message": responseMessage.Message})
}


func GetTickets(c *gin.Context) {
	response, errorMessage := service.GetTickets()
	if errorMessage != nil {
		c.IndentedJSON(errorMessage.HttpStatus, gin.H{"error": errorMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}


func UpdateTicketAvailability(c *gin.Context) {
	pathVariable := c.Param("ticketID")

	ticketID, _ := strconv.Atoi(pathVariable)

	responseMessage, errorMessage := service.UpdateTicketAvailability(ticketID)
	if errorMessage != nil {
		c.IndentedJSON(errorMessage.HttpStatus, gin.H{"error": errorMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, responseMessage)
}


func GetAvailableTickets(c *gin.Context) {
	response, errorMessage := service.GetAvailableTickets()
	if errorMessage != nil {
		c.IndentedJSON(errorMessage.HttpStatus, gin.H{"error": errorMessage})
		return
	}

	c.IndentedJSON(http.StatusOK, response)
}


func GinRouter() {
	router := gin.Default()
	router.POST("/login", middleware.RequireAuth, Login)
	router.POST("/newTicket", middleware.RequireAuth, NewTicket)
	router.GET("/tickets", middleware.RequireAuth, GetTickets)
	router.PUT("/updateTicket/:ticketID", middleware.RequireAuth, UpdateTicketAvailability)
	router.GET("/availableTickets", GetAvailableTickets)
	router.Run()
}