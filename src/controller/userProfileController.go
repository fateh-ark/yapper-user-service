package controller

import (
	"errors"
	"fateh-ark/yapper-user-service/request"
	"fateh-ark/yapper-user-service/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserProfileController interface {
	UpsertUserProfile(c *gin.Context)
	GetUserProfileByUserID(c *gin.Context)
}

type userProfileControllerImpl struct {
	userService service.UserService
}

func NewUserProfileController(userService service.UserService) UserProfileController {
	return &userProfileControllerImpl{userService: userService}
}

func (upc userProfileControllerImpl) UpsertUserProfile(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	var request request.UserProfileReq
	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": ErrInvalidRequest, "details": err.Error()})
		return
	}

	Profile, err := upc.userService.UpsertUserProfile(c.Request.Context(), id, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, Profile)
}

func (upc userProfileControllerImpl) GetUserProfileByUserID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	Profile, err := upc.userService.GetUserProfileByUserID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserPrefNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, Profile)
}
