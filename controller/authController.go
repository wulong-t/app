package controller

import (
	"Server/database"
	"Server/models"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Register
// @Summary Register a new user
// @Description Register a new user with email, password, first name, and last name
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.CreateUser true "user register details"
// @Success 201 {object} models.CreateUser
// @Failure 400 {object} map[string]any
// @Router /user/register [post]
func Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var userRegister models.CreateUser
	if err := c.ShouldBindJSON(&userRegister); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	count, err := gorm.G[models.UserModel](database.DB).Where("email = ?", userRegister.Email).Count(ctx, "*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Email already exists",
			"details": "Please use a different email",
		})
		return
	}
	pwdHash, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Password hash error",
			"details": err.Error(),
		})
		return
	}
	user := models.UserModel{
		Email:    userRegister.Email,
		Name:     userRegister.FirstName + " " + userRegister.LastName,
		Password: string(pwdHash),
	}
	if err := gorm.G[models.UserModel](database.DB).Create(ctx, &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"userID": user.ID,
	})
}
