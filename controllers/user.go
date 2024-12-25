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
	tx := database.DB.Where("email = ?", user.Email).First(&userDb)
	if tx.Error != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": tokenError.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": token})
}

func GetAllUser(c *gin.Context) {
	var users []structs.User
	database.DB.Table("users").Select("id, name, email, role, created_at, updated_at").Find(&users)
	c.JSON(http.StatusOK, gin.H{"result": users})
}

func GetUserById(c *gin.Context) {
	idParam := c.Param("id")
	var userDB structs.User
	if err := database.DB.First(&userDB, idParam).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User is not found", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": userDB})
}

func UpdateUserById(c *gin.Context) {
	idParam := c.Param("id")
	var userDB structs.User
	database.DB.First(&userDB, idParam)
	if userDB.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User is not found"})
		return
	}

	var userUpdate structs.User
	if err := c.ShouldBind(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if userUpdate.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userUpdate.Password), bcrypt.MinCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"result": err})
			return
		}
		userUpdate.Password = string(hashedPassword)
	} else {
		userUpdate.Password = userDB.Password
	}

	if err := database.DB.Model(&userDB).Updates(userUpdate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User has been updated"})
}

func DeleteUserById(c *gin.Context) {
	idParam := c.Param("id")
	var userDB structs.User

	if err := database.DB.First(&userDB, idParam).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User is not found"})
		return
	}

	if err := database.DB.Unscoped().Delete(&userDB).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User has been deleted"})
}
