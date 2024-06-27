// this is a file which is used to manage users get,set,delete operatons of users. here users doesnt always refers to the humans who are interacting with the interface here it refers to cache data or object which are stored , retrieved or deleted.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	cacheinterface "github.com/ragavcr7/caching-library/cache_interface"
)

// UserHandler encapsulates user data operations.
type UserHandler struct {
	userCache cacheinterface.Cache
}

// NewUserHandler initializes a new UserHandler with the provided cache.
func NewUserHandler(userCache cacheinterface.Cache) *UserHandler {
	return &UserHandler{
		userCache: userCache,
	}
}

// SetupRoutes sets up routes for user operations. this by default will fetch from memcache since  NewUserHandler has been only set for memcache server
func (uh *UserHandler) SetupRoutes(router *gin.Engine) {
	router.GET("/user/:id", uh.getUserHandler())
	router.POST("/user", uh.createUserHandler())
	router.DELETE("/user/:id", uh.deleteUserHandler())
	router.GET("/user", uh.getAllUsersHandler())
	router.DELETE("/user", uh.deleteAllUsersHandler())
}

// getUserHandler retrieves a user from the cache based on ID.
func (uh *UserHandler) getUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid user ID",
			})
			return
		}

		key := fmt.Sprintf("user:%d", id)
		cachedValue, err := uh.userCache.Get(key)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("User with ID %d not found", id),
			})
			return
		}

		var user cacheinterface.User
		if err := json.Unmarshal([]byte(cachedValue.(string)), &user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Failed to unmarshal user data: %v", err),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"user":   user,
		})
	}
}

// createUserHandler adds a new user to the cache.
func (uh *UserHandler) createUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser cacheinterface.User
		if err := c.BindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid JSON payload",
			})
			return
		}

		key := fmt.Sprintf("user:%d", newUser.ID)
		userJSON, err := json.Marshal(newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Failed to marshal user data: %v", err),
			})
			return
		}
		if err := uh.userCache.Set(key, string(userJSON), 0); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Failed to store user data in cache: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("User with ID %d created successfully", newUser.ID),
		})
	}
}

// deleteUserHandler removes a user from the cache based on the ID.
func (uh *UserHandler) deleteUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		idParam := c.Param("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid user ID",
			})
			return
		}

		key := fmt.Sprintf("user:%d", id)
		if err := uh.userCache.Delete(key); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Failed to delete user with ID %d: %v", id, err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": fmt.Sprintf("User with ID %d deleted successfully", id),
		})
	}
}

// getAllUsersHandler retrieves all users from the cache.
func (uh *UserHandler) getAllUsersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		keys, err := uh.userCache.GetAllKeys()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Failed to fetch all users from cache: %v", err),
			})
			return
		}

		users := make([]cacheinterface.User, 0, len(keys))
		for _, key := range keys {
			cachedValue, err := uh.userCache.Get(key)
			if err != nil {
				continue // Skip this user if retrieval fails
			}
			var user cacheinterface.User
			err = json.Unmarshal([]byte(cachedValue.(string)), &user)
			if err != nil {
				continue // Skip this user if unmarshalling fails
			}
			users = append(users, user)
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"users":  users,
		})
	}
}

// deleteAllUsersHandler removes all users from the cache.
func (uh *UserHandler) deleteAllUsersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := uh.userCache.DeleteAllKeys()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Failed to delete all users from cache: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "All users deleted successfully",
		})
	}
}
