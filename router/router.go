package router

import (
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"glimpseguru-api/analytics"
	"glimpseguru-api/authent"
	"net/http"
	"time"
)

func getUserAnalytics(c *gin.Context) {
	var errUser error
	var user authent.User
	var isUserStruct bool
	if storedUser, exists := c.Get("user"); exists {
		if user, isUserStruct = storedUser.(authent.User); !isUserStruct {
			errUser = errors.New("cannot bind user from context")
		}
	} else {
		errUser = errors.New("no user found in context")
	}
	if errUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("error setting user: %e", errUser)})
		return
	}
	var request AnalyticsRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	analyticsService := analytics.NewService(request.StartTime, request.EndTime)
	userAnalytics, errsUserAnalytics := analyticsService.GetUsersAnalytics(user)
	if len(errsUserAnalytics) > 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.Join(errsUserAnalytics...)})
		return
	}
	c.JSON(http.StatusOK, userAnalytics)
}

func New() *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "OPTIONS"},
		AllowHeaders:    []string{"x-api-key", "Content-Type"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}))
	r.Use(identityValidationMiddleware())
	r.GET("/userAnalytics", getUserAnalytics)

	return r
}

func identityValidationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "API Key is required in headers"})
			return
		}

		user, errAuthent := authent.GetUser(apiKey)
		if errAuthent != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
