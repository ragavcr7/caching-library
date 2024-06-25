/*The routes.go file is where you define the routes (endpoints) for your
  API and link them to the appropriate handler functions. setting up routes for caching operations (GET, POST, DELETE, etc...) and user management.*/

package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/ragavcr7/caching-library/cache_interface"
)

// SetupRouter sets up the routes for the API and links them to the appropriate handler functions.
func SetupRouter(memcachedCache cache_interface.Cache, redisCache cache_interface.Cache, inMemoryCache cache_interface.Cache) *gin.Engine {
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
