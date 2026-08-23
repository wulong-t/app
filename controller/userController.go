package controller

import (
	"Server/database"
	"Server/models"
	"context"
	"errors"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetUserByID
// @Summary Get user by ID
// @Description Get user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Security BearerAuth
// @Success 200 {object} models.UserModel
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /user/getuser/{id} [get]
func GetUserByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	userID := c.Param("id")
	user, err := gorm.G[models.UserModel](database.DB).Where("id = ?", userID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateUser
// @Summary Update user
// @Description Update user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.UpdateUser true "User"
// @Security BearerAuth
// @Success 200 {object} models.UpdateUser
// @Failure 400 {object} map[string]any
// @Failure 401 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /user/updateuser/{id} [patch]
func UpdateUser(c *gin.Context) {
	updateID := c.Param("id")
	userID := c.GetString("userID")
	if updateID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not authorized to update this user"})
		return
	}
	var user models.UpdateUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	updateUser := map[string]any{}
	if user.Name != "" {
		updateUser["name"] = user.Name
	}
	if user.Bio != "" {
		updateUser["bio"] = user.Bio
	}
	if user.ImageURL != "" {
		updateUser["image_url"] = user.ImageURL
	}
	result := database.DB.WithContext(ctx).Model(&models.UserModel{}).Where("id = ?", updateID).Updates(updateUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"details": result.Error.Error(),
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// FollowUser
// @Summary Follow user
// @Description follow user, if user is already followed, unfollow
// @Tags Users
// @Param id path string true "Target ID"
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /user/follow/{id} [patch]
func FollowUser(c *gin.Context) {
	currentID := c.GetString("userID")
	targetID := c.Param("id")
	targetU, err := strconv.ParseUint(targetID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"details": err.Error(),
		})
		return
	}
	currentU, err := strconv.ParseUint(currentID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"details": err.Error(),
		})
		return
	}
	currentUID := uint(currentU)
	targetUID := uint(targetU)
	if currentUID == targetUID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot follow yourself"})
		return
	}

	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	_, err = gorm.G[models.UserModel](database.DB).Where("id = ?", targetUID).First(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user",
			"details": err.Error(),
		})
		return
	}
	_, err = gorm.G[models.FollowModel](database.DB).Where("follower_id = ? AND followed_id = ?", currentUID, targetUID).First(ctx)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to check follow",
				"details": err.Error(),
			})
			return
		}
		follow := models.FollowModel{FollowerID: currentUID, FollowedID: targetUID}
		err = gorm.G[models.FollowModel](database.DB).Create(ctx, &follow)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to follow user",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Followed user",
		})
	} else {
		_, err = gorm.G[models.FollowModel](database.DB).Where("follower_id = ? && followed_id = ?", currentUID, targetUID).Delete(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to unfollow user",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Unfollowed user",
		})
	}
}

// GetSugUser
// @Summary Get suggested users
// @Description Get suggested users by following
// @Tags Users
// @Produce json
// @Param id query string true "User ID"
// @Security BearerAuth
// @Success 200 {object} []models.UserModel
// @Failure 400 {object} map[string]any
// @Failure 404 {object} map[string]any
// @Failure 500 {object} map[string]any
// @Router /user/suguser [get]
func GetSugUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	mainID := c.Query("id")
	mainU, err := strconv.ParseUint(mainID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid user ID",
			"details": err.Error(),
		})
		return
	}
	mainUID := uint(mainU)
	followingUsers, err := gorm.G[models.FollowModel](database.DB).Where("follower_id = ?", mainUID).Find(ctx)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User no follow any users"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get users",
			"details": err.Error(),
		})
		return
	}
	followingIDs := []uint{}
	for _, record := range followingUsers {
		followingIDs = append(followingIDs, record.FollowedID)
	}
	ffs, err := gorm.G[models.FollowModel](database.DB).Where("follower_id IN ?", followingIDs).Find(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get users",
			"details": err.Error(),
		})
		return
	}
	sugUsers := make(map[uint]struct{})
	for _, record := range ffs {
		if record.FollowedID == mainUID {
			continue
		}
		if slices.Contains(followingIDs, record.FollowedID) {
			continue
		}
		sugUsers[record.FollowedID] = struct{}{}
	}
	if len(sugUsers) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "No suggestions",
		})
		return
	}
	sugIDs := make([]uint, 0, len(sugUsers))
	for id := range sugUsers {
		sugIDs = append(sugIDs, id)
	}
	users, err := gorm.G[models.UserModel](database.DB).Where("id IN ?", sugIDs).Find(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get users",
			"details": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, users)
}
