package controllers

import (
	"log"
	"net/http"
	"strconv"

	"backendsetup/m/db/sql/dbgen"
	"backendsetup/m/middleware"
	"backendsetup/m/services"

	"github.com/gin-gonic/gin"
)

type CountersController struct {
	countersService *services.CountersService
	userService *services.UserService
}

func InitCountersController(countersService *services.CountersService, userService *services.UserService) CountersController {
	return CountersController{
		countersService,
		userService,
	}
}

func (c *CountersController) GetCountersForOtherUser(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	fromUserID := claims.(middleware.Claims).Sub
	username,exists := ctx.Params.Get("username")
	if !exists{
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
		return

	}
	user, err := c.userService.GetUserByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	counters, err := c.countersService.GetCountersForUserFrom(user.UserID, fromUserID,)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't find counters"})
		log.Printf("failed retrieving counters: %v\r\n", err)
		return
	}
	ctx.JSON(http.StatusOK, counters)
}

func (c *CountersController) GetCounters(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	userID := claims.(middleware.Claims).Sub
	counters, err := c.countersService.GetCountersForUser(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "couldn't find counters"})
		log.Printf("failed retrieving counters: %v\r\n", err)
		return
	}
	ctx.JSON(http.StatusOK, counters)
}

func (c *CountersController) CreateCounter(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}

	type createCounterBody struct {
		Name string `json:"name"`
		Icon string `json:"icon"`
		Users []dbgen.FindUserRow `json:"users"`
	}
	var body createCounterBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect body supplied"})
		log.Printf("Bad Request: %v\r\n", err)
		return
	}
	counter, err := c.countersService.CreateCounter(body.Name, body.Icon, claims.(middleware.Claims).Sub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not create counter"})
		log.Printf("could not create counter: %v\r\n", err)
		return
	}
	for _, user := range(body.Users) {
		err = c.countersService.ShareCounter(user.UserID, int(counter.ID), "participant", claims.(middleware.Claims).Sub)
		if err != nil {
			log.Printf("could not share to user %s: %v", user.UserID, err)
		}
	}
	ctx.JSON(http.StatusOK, counter)
}

func (c *CountersController) ShareCounter(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	type shareCounterBody struct {
		CounterID int    `json:"counterID"`
		Recipient string `json:"recipientID"`
	}
	var body shareCounterBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect body supplied"})
		return
	}

	err := c.countersService.ShareCounter(body.Recipient, body.CounterID, "participant", claims.(middleware.Claims).Sub)
	if err != nil {
		if err.Error() == "no permission" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized to share counter"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not share counter"})
		log.Printf("could not share counter to owner: %v\r\n", err)
		return
	}
	ctx.JSON(http.StatusOK, "done")
}

func (c *CountersController) AddEvent(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}

	counterID, err := strconv.Atoi(ctx.Param("counterID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "incorrect params passed"})
		return
	}

	err = c.countersService.AddEvent(counterID, claims.(middleware.Claims).Sub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not add event"})
		return
	}
	ctx.JSON(http.StatusOK, "done")
}

func (c *CountersController) GetCounter(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	counterID, err := strconv.Atoi(ctx.Param("counterID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "counterID not passed"})
		return
	}
	counter, err := c.countersService.GetCounter(int(counterID), claims.(middleware.Claims).Sub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not find counter"})
		return
	}
	ctx.JSON(http.StatusOK, counter)
}

func (c *CountersController) GetEvents(ctx *gin.Context) {
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "claims not found"})
		return
	}
	counterID, err := strconv.Atoi(ctx.Param("counterID"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "counterID not passed"})
		return
	}
	events, err := c.countersService.GetEvents(counterID, claims.(middleware.Claims).Sub)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "could not find events"})
		return
	}
	ctx.JSON(http.StatusOK, events)
}
