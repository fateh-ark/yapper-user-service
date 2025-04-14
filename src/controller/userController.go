package controller

import (
	"errors"
	"fateh-ark/yapper-user-service/request"
	"fateh-ark/yapper-user-service/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController interface {
	CreateUser(c *gin.Context)
	GetUserByID(c *gin.Context)
	GetUserByUsername(c *gin.Context)
	GetUserByEmail(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
}

type userControllerImpl struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &userControllerImpl{userService: userService}
}

func (uc *userControllerImpl) CreateUser(c *gin.Context) {
	var requestUser request.CreateUserReq
	if err := c.ShouldBindBodyWithJSON(&requestUser); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": ErrInvalidRequest, "details": err.Error()})
		return
	}

	newUser, err := uc.userService.CreateUser(c.Request.Context(), &requestUser)

	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": newUser.ID, "created_at": newUser.CreatedAt})
}

func (uc *userControllerImpl) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	user, err := uc.userService.GetUserByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *userControllerImpl) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := uc.userService.GetUserByUsername(c.Request.Context(), username)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *userControllerImpl) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")

	user, err := uc.userService.GetUserByEmail(c.Request.Context(), email)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *userControllerImpl) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	var requestUser request.UpdateUserReq
	if err := c.ShouldBindBodyWithJSON(&requestUser); err != nil {
		c.JSON(
			http.StatusBadRequest,
			gin.H{"error": ErrInvalidRequest, "details": err.Error()})
		return
	}

	updatedUser, err := uc.userService.UpdateUser(c.Request.Context(), id, &requestUser)

	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, updatedUser)
}

func (uc *userControllerImpl) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidUserId})
		return
	}

	if err := uc.userService.DeleteUser(c.Request.Context(), id); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
