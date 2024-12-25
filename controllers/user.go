package controllers

import (
	"event-management/database"
	"event-management/helpers"
	"event-management/structs"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var err error

func AddNewUser(c *gin.Context) {
	var newUser structs.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := newUser.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.MinCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	newUser.Password = string(hashedPassword)

	if err := database.DB.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var user structs.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	if user.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email is required"})
		return
	} else if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Password is required"})
		return
	}

	var userDb structs.User
	fmt.Println(user.Email)
	tx := database.DB.Where("email = ?", user.Email).First(&userDb)
	if tx.Error != nil {
		fmt.Println("Email", user.Email, "Ngga ke cek")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Email / Password"})
		return
	}

	result := bcrypt.CompareHashAndPassword([]byte(userDb.Password), []byte(user.Password))
	if result != nil {
		fmt.Println("Ini pas cek passwordnya bener engga dengan yang di hash")
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid Email / Password"})
		return
	}

	token, tokenError := helpers.SignPayload(userDb)
	if tokenError != nil {
		fmt.Println(tokenError)
		c.JSON(http.StatusInternalServerError, "Internal Server Error")
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}
