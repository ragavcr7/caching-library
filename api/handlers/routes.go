/*The routes.go file is where you define the routes (endpoints) for your
  API and link them to the appropriate handler functions. setting up routes for caching operations (GET, POST, DELETE, etc...) and user management.*/

package handlers

import (
	"github.com/gin-gonic/gin"
	cache "github.com/ragavcr7/caching-library/api"
)

func SetupRouter(memcachedCache cache.Cache, redisCache cache.Cache, inMemoryCache cache.Cache) *gin.Engine {
	router := gin.Default()

	// Initialize handlers
	userHandler := NewUserHandler(memcachedCache)

	// User routes
	router.POST("/user", userHandler.createUserHandler())
	router.GET("/user/:id", userHandler.getUserHandler())
	router.DELETE("/user/:id", userHandler.deleteUserHandler())
	router.DELETE("/users", userHandler.deleteAllUsersHandler())

	return router
}
