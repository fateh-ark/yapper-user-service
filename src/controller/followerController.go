package controller

import (
	"fateh-ark/yapper-user-service/request"
	"fateh-ark/yapper-user-service/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type FollowerController interface {
	FollowUser(c *gin.Context)
	UnfollowUser(c *gin.Context)
	GetFollowers(c *gin.Context)
	GetFollowing(c *gin.Context)
	IsFollowing(c *gin.Context)
	GetFollowStats(c *gin.Context)
}

type followerControllerImpl struct {
	userService service.UserService
}

func NewFollowerController(userService service.UserService) FollowerController {
	return &followerControllerImpl{userService: userService}
}

func (fc *followerControllerImpl) FollowUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	var request request.FollowReq
	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": ErrInvalidRequest, "details": err.Error()})
		return
	}

	if err := fc.userService.FollowUser(c.Request.Context(), id, &request); err != nil {
		errMessage := err.Error()
		if strings.Contains(errMessage, service.ErrUserNotFound.Error()) {
			c.JSON(http.StatusNotFound, gin.H{"error": errMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
		return
	}
	c.Status(http.StatusNoContent)
}

func (fc *followerControllerImpl) UnfollowUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	var request request.FollowReq
	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": ErrInvalidRequest, "details": err.Error()})
		return
	}

	if err := fc.userService.UnfollowUser(c.Request.Context(), id, &request); err != nil {
		errMessage := err.Error()
		if strings.Contains(errMessage, service.ErrUserNotFound.Error()) {
			c.JSON(http.StatusNotFound, gin.H{"error": errMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
		return
	}
	c.Status(http.StatusNoContent)
}

func (fc *followerControllerImpl) GetFollowers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	followers, err := fc.userService.GetFollowers(c.Request.Context(), id)
	if err != nil {
		errMessage := err.Error()
		if strings.Contains(errMessage, service.ErrUserNotFound.Error()) {
			c.JSON(http.StatusNotFound, gin.H{"error": errMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
		return
	}
	c.JSON(http.StatusOK, followers)
}

func (fc *followerControllerImpl) GetFollowing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	following, err := fc.userService.GetFollowing(c.Request.Context(), id)
	if err != nil {
		errMessage := err.Error()
		if strings.Contains(errMessage, service.ErrUserNotFound.Error()) {
			c.JSON(http.StatusNotFound, gin.H{"error": errMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
		return
	}
	c.JSON(http.StatusOK, following)
}

func (fc *followerControllerImpl) IsFollowing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	var request request.FollowReq
	if err := c.ShouldBindBodyWithJSON(&request); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": ErrInvalidRequest, "details": err.Error()})
		return
	}

	isFollowing, err := fc.userService.Isfollowing(c.Request.Context(), id, &request)
	if err != nil {
		errMessage := err.Error()
		if strings.Contains(errMessage, service.ErrUserNotFound.Error()) {
			c.JSON(http.StatusNotFound, gin.H{"error": errMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"is_following": isFollowing})
}

func (fc *followerControllerImpl) GetFollowStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	stat, err := fc.userService.GetFollowStats(c.Request.Context(), id)
	if err != nil {
		errMessage := err.Error()
		if strings.Contains(errMessage, service.ErrUserNotFound.Error()) {
			c.JSON(http.StatusNotFound, gin.H{"error": errMessage})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
		return
	}
	c.JSON(http.StatusOK, stat)
}
