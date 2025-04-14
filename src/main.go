package main

import (
	"context"
	"fateh-ark/yapper-user-service/controller"
	"fateh-ark/yapper-user-service/repositories"
	"fateh-ark/yapper-user-service/service"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var dbPool *pgxpool.Pool // Use pgxpool.Pool

func main() {
	var err error

	// Load the env file if any.
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using environment variables directly")
	}

	// Setup the PGPool Connection
	dbUser := os.Getenv("POSTGRESQL_USERNAME")
	dbPassword := os.Getenv("POSTGRESQL_PASSWORD")
	pgPoolPort := os.Getenv("PGPOOL_PORT_NUMBER")
	dbName := os.Getenv("POSTGRESQL_DATABASE")

	connStr := fmt.Sprintf("postgres://%s:%s@pgpool:%s/%s?sslmode=disable", url.QueryEscape(dbUser), url.QueryEscape(dbPassword), pgPoolPort, url.QueryEscape(dbName))

	dbPool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	// Setup the Repositories
	userRepo := repositories.NewUserRepository(dbPool)
	followerRepo := repositories.NewFollowerRepository(dbPool)
	userProfileRepo := repositories.NewUserProfileRepository(dbPool)
	userPreferenceRepo := repositories.NewUserPreferenceRepository(dbPool)

	// Setup the User service
	userService := service.NewUserService(
		userRepo,
		followerRepo,
		userProfileRepo,
		userPreferenceRepo,
	)

	// Setup the Controllers instances
	userController := controller.NewUserController(userService)
	followerController := controller.NewFollowerController(userService)
	userProfileController := controller.NewUserProfileController(userService)
	userPreferenceController := controller.NewUserPreferenceController(userService)

	// Setup Gin Routes
	router := gin.Default()

	userApi := router.Group("/user/v1")
	{
		// User
		userApi.POST("", userController.CreateUser)
		userApi.GET("/username/:username", userController.GetUserByUsername)
		userApi.GET("/email/:email", userController.GetUserByEmail)

		userIdGroup := userApi.Group("/:id")
		{
			// Remaining Users
			userIdGroup.GET("", userController.GetUserByID)
			userIdGroup.PUT("", userController.UpdateUser)
			userIdGroup.DELETE("", userController.DeleteUser)

			// Follower
			userIdGroup.PUT("/follow", followerController.FollowUser)
			userIdGroup.PUT("/unfollow", followerController.UnfollowUser)
			userIdGroup.GET("/followers", followerController.GetFollowers)
			userIdGroup.GET("/following", followerController.GetFollowing)
			userIdGroup.GET("/isFollowing", followerController.IsFollowing)
			userIdGroup.GET("/followStats", followerController.GetFollowStats)

			// User Profile
			userIdGroup.PUT("/profile", userProfileController.UpsertUserProfile)
			userIdGroup.GET("/profile", userProfileController.GetUserProfileByUserID)

			// User Preference
			userIdGroup.PUT("/preference", userPreferenceController.UpsertUserPreference)
			userIdGroup.GET("/preference", userPreferenceController.GetUserPreferenceByUserID)
		}
	}

	apiPort := os.Getenv("USER_SERVICE_PORT_NUMBER")
	router.Run(":" + apiPort)
}
