package controller

import (
	"errors"
	"fateh-ark/yapper-user-service/request"
	"fateh-ark/yapper-user-service/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserPreferenceController interface {
	UpsertUserPreference(c *gin.Context)
	GetUserPreferenceByUserID(c *gin.Context)
}

type userPreferenceControllerImpl struct {
	userService service.UserService
}

func NewUserPreferenceController(userService service.UserService) UserPreferenceController {
	return &userPreferenceControllerImpl{userService: userService}
}

func (upc userPreferenceControllerImpl) UpsertUserPreference(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	var request request.UserPreferenceReq
	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": ErrInvalidRequest, "details": err.Error()})
		return
	}

	preference, err := upc.userService.UpsertUserPreference(c.Request.Context(), id, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, preference)
}

func (upc userPreferenceControllerImpl) GetUserPreferenceByUserID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	preference, err := upc.userService.GetUserPreferenceByUserID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserPrefNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, preference)
}
