// this is a file which is used to manage users get,set,delete operatons of users. here users doesnt always refers to the humans who are interacting with the interface here it refers to cache data or object which are stored , retrieved or deleted.
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ragavcr7/caching-library/cache"
)

// UserHandler encapsulates user data operations.
type UserHandler struct {
	userCache cache.Cache
}

// NewUserHandler initializes a new UserHandler with the provided cache.
func NewUserHandler(userCache cache.Cache) *UserHandler {
	return &UserHandler{
		userCache: userCache,
	}
}

// SetupRoutes sets up routes for user operations.
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

		// Attempt to retrieve user from cache
		key := fmt.Sprintf("user:%d", id)
		cachedValue, found := uh.userCache.Get(key)
		if found != nil {
			var user cache.User
			err := json.Unmarshal([]byte(cachedValue.(string)), &user)
			if err != nil {
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
		} else {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("User with ID %d not found", id),
			})
		}
	}
}

// createUserHandler adds a new user to the cache.
func (uh *UserHandler) createUserHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newUser cache.User
		if err := c.BindJSON(&newUser); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid JSON payload",
			})
			return
		}

		// Store user in cache
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

		// Delete user from cache
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
		// Fetch all keys from cache (assuming cache supports this operation)
		keys, err := uh.userCache.FetchAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Failed to fetch all users from cache: %v", err),
			})
			return
		}

		users := make([]cache.User, 0, len(keys))
		for _, key := range keys {
			cachedValue, found := uh.userCache.Get(key)
			if found != nil {
				var user cache.User
				err := json.Unmarshal([]byte(cachedValue.(string)), &user)
				if err != nil {
					continue // Skip this user if unmarshalling fails
				}
				users = append(users, user)
			}
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
		// Fetch all keys from cache (assuming cache supports this operation)
		keys, err := uh.userCache.FetchAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Failed to fetch all users from cache: %v", err),
			})
			return
		}

		// Delete all keys
		for _, key := range keys {
			if err := uh.userCache.Delete(key); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": fmt.Sprintf("Failed to delete user with key %s: %v", key, err),
				})
				return
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "All users deleted successfully",
		})
	}
}
