package main

import (
	"context"
	"fateh-ark/yapper-user-service/controller"
	"fateh-ark/yapper-user-service/logger"
	"fateh-ark/yapper-user-service/repositories"
	"fateh-ark/yapper-user-service/service"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
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

	connStr := fmt.Sprintf(
		"postgres://%s:%s@pgpool:%s/%s?sslmode=disable",
		url.QueryEscape(dbUser),
		url.QueryEscape(dbPassword),
		pgPoolPort, url.QueryEscape(dbName),
	)

	// log.Println("connecting pg to:", connStr)

	dbPool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("failed to connect to database pool: %v", err)
	}
	defer dbPool.Close()

	// Setup RabbitMQ Connection
	rmqUser := os.Getenv("RABBITMQ_USERNAME")
	rmqPassword := os.Getenv("RABBITMQ_PASSWORD")
	rmqPort := os.Getenv("RABBITMQ_PORT_NUMBER")

	rmqConnStr := fmt.Sprintf("amqp://%s:%s@rabbitmq:%s/",
		url.QueryEscape(rmqUser),
		url.QueryEscape(rmqPassword),
		rmqPort,
	)

	log.Println("connecting rmq to:", rmqConnStr)

	rmqConn, err := amqp.Dial(rmqConnStr)
	if err != nil {
		log.Fatal("failed to connect to rabbitmq:", err)
	}
	defer rmqConn.Close()

	// Setup Logger
	loggerInstance, err := logger.NewLogger(rmqConn, "service_log")
	if err != nil {
		log.Fatal("failed to setup logger:", err)
	}
	defer loggerInstance.CloseChannel()

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

	router.Use(controller.LoggerHandler(loggerInstance))

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
