package controller

import (
	"Server/database"
	"Server/models"
	"context"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
// @Success 201 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
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
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})
	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "JWT error",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"Result": user,
		"Token":  token,
	})
}

// Login
// @Summary Login a user
// @Description Login a user with email, password
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body models.LoginUser true "user Login details"
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /user/login [post]
func Login(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	var userLogin models.LoginUser
	if err := c.ShouldBindJSON(&userLogin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	user, err := gorm.G[models.UserModel](database.DB).Where("email = ?", userLogin.Email).First(ctx)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Email does not exist",
			"details": "Please use a different email",
		})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Incorrect password",
			"details": "Please use a different password",
		})
		return
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   strconv.Itoa(int(user.ID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
	})

	token, err := claims.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "JWT error",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Result": user,
		"Token":  token,
	})
}
